package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"groundTurn/src/config"
	"groundTurn/src/constants"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type apiClient struct {
	t   *testing.T
	r   *gin.Engine
	tok string
}

func (a *apiClient) do(method, path string, body interface{}) (int, map[string]interface{}) {
	a.t.Helper()
	var reader *bytes.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if a.tok != "" {
		req.Header.Set("Authorization", "Bearer "+a.tok)
	}
	w := httptest.NewRecorder()
	a.r.ServeHTTP(w, req)
	var out map[string]interface{}
	if w.Body.Len() > 0 {
		if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
			a.t.Fatalf("decode %s %s: %v body=%s", method, path, err, w.Body.String())
		}
	}
	return w.Code, out
}

func (a *apiClient) login(username string) {
	code, out := a.do(http.MethodPost, "/api/auth/login", map[string]string{"username": username})
	if code != 200 {
		a.t.Fatalf("login %s: %d %v", username, code, out)
	}
	a.tok, _ = out["token"].(string)
}

func setupAPITest(t *testing.T) (*apiClient, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_busy_timeout=5000"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.Migrate(db); err != nil {
		t.Fatal(err)
	}
	db.Exec("DELETE FROM resource_booking")
	db.Exec("DELETE FROM ground_task")
	db.Exec("DELETE FROM delay_event")
	db.Exec("DELETE FROM ground_resource")
	db.Exec("DELETE FROM flight_turnaround")
	db.Exec("DELETE FROM audit_log")
	db.Exec("DELETE FROM app_user")
	if err := config.Seed(db); err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{JWTSecret: "test-secret", RateLimitQPS: 1000}
	r := NewRouter(db, cfg)
	return &apiClient{t: t, r: r}, db
}

func asMap(v interface{}) map[string]interface{} {
	m, _ := v.(map[string]interface{})
	return m
}

func asList(v interface{}) []interface{} {
	l, _ := v.([]interface{})
	return l
}

// TestAPIClosedLoop drives the real HTTP stack: rejected batch, fixes,
// saved plan, sign-off, delay reschedule with frozen task + pending
// bookings, adjustment, and RBAC enforcement.
func TestAPIClosedLoop(t *testing.T) {
	c, _ := setupAPITest(t)

	// Unauthenticated request is rejected.
	anonymous := &apiClient{t: t, r: c.r}
	if code, _ := anonymous.do(http.MethodGet, "/api/turnarounds", nil); code != 401 {
		t.Fatalf("expected 401 without token, got %d", code)
	}

	// A team user cannot register turnarounds (RBAC).
	team := &apiClient{t: t, r: c.r}
	team.login("team-bag")
	if code, out := team.do(http.MethodPost, "/api/turnarounds", map[string]string{
		"flight_no": "ZZ999", "aircraft_reg": "B-1", "stand_no": "1",
		"arrival_time": "2026-09-22T08:00:00Z", "departure_time": "2026-09-22T09:00:00Z",
	}); code != 403 {
		t.Fatalf("team register expected 403, got %d %v", code, out)
	}

	// Dispatcher logs in.
	c.login("dispatcher")

	// Find CA101 id + arrival.
	code, out := c.do(http.MethodGet, "/api/turnarounds?page_size=100", nil)
	if code != 200 {
		t.Fatalf("list turnarounds: %d", code)
	}
	var caID float64
	var arrival string
	for _, item := range asList(out["items"]) {
		row := asMap(item)
		if row["flight_no"] == "CA101" {
			caID, _ = row["id"].(float64)
			arrival, _ = row["arrival_time"].(string)
		}
	}
	if caID == 0 {
		t.Fatal("CA101 not seeded")
	}
	base, err := time.Parse(time.RFC3339, arrival)
	if err != nil {
		t.Fatal(err)
	}
	at := func(min int) string { return base.Add(time.Duration(min) * time.Minute).Format(time.RFC3339) }

	// 1. First plan generation: whole batch rejected, nothing saved.
	code, out = c.do(http.MethodPost, "/api/turnarounds/"+itoa(int(caID))+"/plan", map[string][]string{})
	if code != 409 || out["code"] != constants.CodeConflict {
		t.Fatalf("expected 409 PLAN_CONFLICT, got %d %v", code, out["code"])
	}
	details := asList(out["details"])
	if len(details) == 0 {
		t.Fatal("rejected batch must carry itemized conflicts")
	}
	codes := map[string]int{}
	for _, d := range details {
		codes[asMap(d)["code"].(string)]++
	}
	if codes[constants.ConflictOffline] == 0 || codes[constants.ConflictMaintenance] == 0 || codes[constants.ConflictOverlap] == 0 {
		t.Fatalf("expected offline/maintenance/overlap reasons, got %v", codes)
	}
	// Nothing persisted.
	code, detail0 := c.do(http.MethodGet, "/api/turnarounds/"+itoa(int(caID)), nil)
	if code != 200 {
		t.Fatalf("detail: %d", code)
	}
	if len(asList(detail0["tasks"])) != 0 {
		t.Fatal("rejected batch must not persist tasks")
	}

	// 2. Resource manager clears blockers: BL-02 online; FT-02 online with
	//    maintenance window 10:40-11:15; release MU202 fuel booking.
	rm := &apiClient{t: t, r: c.r}
	rm.login("resource")
	code, resOut := rm.do(http.MethodGet, "/api/ground-resources?page_size=200", nil)
	if code != 200 {
		t.Fatal("list resources")
	}
	idByCode := map[string]int{}
	for _, item := range asList(resOut["items"]) {
		row := asMap(item)
		idByCode[row["resource_code"].(string)] = int(row["id"].(float64))
	}
	if code, out = rm.do(http.MethodPatch, "/api/ground-resources/"+itoa(idByCode["BL-02"])+"/status",
		map[string]string{"status": "AVAILABLE"}); code != 200 {
		t.Fatalf("BL-02 online: %d %v", code, out)
	}
	if code, out = rm.do(http.MethodPatch, "/api/ground-resources/"+itoa(idByCode["FT-02"])+"/status",
		map[string]string{"status": "AVAILABLE", "maintenance_start": at(40), "maintenance_end": at(75)}); code != 200 {
		t.Fatalf("FT-02 status: %d %v", code, out)
	}
	code, bookOut := rm.do(http.MethodGet, "/api/resource-bookings?page_size=200", nil)
	if code != 200 {
		t.Fatal("list bookings")
	}
	for _, item := range asList(bookOut["items"]) {
		row := asMap(item)
		if row["resource_code"] == "FT-02" && row["booking_status"] != constants.BookingReleased {
			bid := itoa(int(row["id"].(float64)))
			if code, out = rm.do(http.MethodPost, "/api/resource-bookings/"+bid+"/release", nil); code != 200 {
				t.Fatalf("release FT-02 booking: %d %v", code, out)
			}
		}
	}

	// 3. Regenerate plan: whole batch saved.
	if code, out = c.do(http.MethodPost, "/api/turnarounds/"+itoa(int(caID))+"/plan", map[string][]string{}); code != 201 {
		t.Fatalf("plan expected 201, got %d %v", code, out)
	}
	result := asMap(out)
	if result["saved"] != true || len(asList(result["tasks"])) != 6 {
		t.Fatalf("expected 6 saved tasks, got %v", result["saved"])
	}

	// 4. Team signs the baggage task; refresh detail shows ACCEPTED.
	code, detail1 := c.do(http.MethodGet, "/api/turnarounds/"+itoa(int(caID)), nil)
	if code != 200 {
		t.Fatal("detail after plan")
	}
	var baggageID int
	taskByType := map[string]map[string]interface{}{}
	for _, item := range asList(detail1["tasks"]) {
		row := asMap(item)
		taskByType[row["task_type"].(string)] = row
		if row["task_type"] == constants.TaskTypeBaggage {
			baggageID = int(row["id"].(float64))
		}
	}
	if code, out = team.do(http.MethodPost, "/api/ground-tasks/"+itoa(baggageID)+"/accept", nil); code != 200 {
		t.Fatalf("accept baggage: %d %v", code, out)
	}
	code, detailAccepted := c.do(http.MethodGet, "/api/turnarounds/"+itoa(int(caID)), nil)
	if taskByType2(detailAccepted, constants.TaskTypeBaggage)["status"] != constants.TaskAccepted {
		t.Fatal("refreshed detail must show baggage ACCEPTED")
	}

	// 5. Register +30 minutes: unsigned tasks move, signed baggage freezes,
	//    CK-02 and FT-02 bookings go PENDING.
	if code, out = c.do(http.MethodPost, "/api/turnarounds/"+itoa(int(caID))+"/delays",
		map[string]interface{}{"delay_type": "LATE_ARRIVAL", "minutes": 30,
			"root_cause": "前站流量控制", "responsibility_team": "AIRLINE"}); code != 201 {
		t.Fatalf("delay expected 201, got %d %v", code, out)
	}
	if int(out["rescheduled_tasks"].(float64)) != 5 {
		t.Fatalf("expected 5 rescheduled unsigned tasks, got %v", out["rescheduled_tasks"])
	}
	if int(out["frozen_tasks"].(float64)) != 1 {
		t.Fatalf("expected 1 frozen signed task, got %v", out["frozen_tasks"])
	}
	if int(out["pending_bookings"].(float64)) < 2 {
		t.Fatalf("expected >=2 pending bookings, got %v", out["pending_bookings"])
	}

	// Refreshed detail: same result visible.
	code, detail2 := c.do(http.MethodGet, "/api/turnarounds/"+itoa(int(caID)), nil)
	if code != 200 {
		t.Fatal("detail after delay")
	}
	baggageAfter := taskByType2(detail2, constants.TaskTypeBaggage)
	expectedBaggageStart := base.Add(2 * time.Minute).Format(time.RFC3339) // unsigned? no — signed, stays 10:02
	if baggageAfter["planned_start"] != expectedBaggageStart {
		t.Fatalf("signed baggage must freeze at %s, got %v", expectedBaggageStart, baggageAfter["planned_start"])
	}
	pendingCodes := map[string]string{}
	for _, item := range asList(detail2["bookings"]) {
		row := asMap(item)
		if row["booking_status"] == constants.BookingPending {
			pendingCodes[row["resource_code"].(string)] = row["conflict_code"].(string)
		}
	}
	if pendingCodes["CK-02"] != constants.ConflictOverlap {
		t.Fatalf("CK-02 must be PENDING overlap, got %v", pendingCodes)
	}
	if pendingCodes["FT-02"] != constants.ConflictMaintenance {
		t.Fatalf("FT-02 must be PENDING maintenance, got %v", pendingCodes)
	}
	// Unsigned cleaning moved +30m.
	cleaningAfter := taskByType2(detail2, constants.TaskTypeCleaning)
	if cleaningAfter["planned_start"] != base.Add(35*time.Minute).Format(time.RFC3339) {
		t.Fatalf("cleaning should shift to 10:35, got %v", cleaningAfter["planned_start"])
	}

	// Dashboard exposes the same pending bookings.
	if code, out = c.do(http.MethodGet, "/api/dashboard/pending-bookings", nil); code != 200 ||
		int(out["total"].(float64)) < 2 {
		t.Fatalf("dashboard pending mismatch: %d %v", code, out)
	}

	// 6. Resource manager adjusts CK-02 pending booking to a clean window.
	var ck02Booking int
	for _, item := range asList(detail2["bookings"]) {
		row := asMap(item)
		if row["resource_code"] == "CK-02" && row["booking_status"] == constants.BookingPending {
			ck02Booking = int(row["id"].(float64))
		}
	}
	if code, out = rm.do(http.MethodPost, "/api/resource-bookings/"+itoa(ck02Booking)+"/adjust",
		map[string]string{"start_time": at(70), "end_time": at(85)}); code != 200 {
		t.Fatalf("adjust CK-02: %d %v", code, out)
	}
	// Adjust FT-02 into maintenance stays PENDING with 409.
	var ft02Booking int
	for _, item := range asList(detail2["bookings"]) {
		row := asMap(item)
		if row["resource_code"] == "FT-02" && row["booking_status"] == constants.BookingPending {
			ft02Booking = int(row["id"].(float64))
		}
	}
	if code, out = rm.do(http.MethodPost, "/api/resource-bookings/"+itoa(ft02Booking)+"/adjust",
		map[string]string{"start_time": at(45), "end_time": at(65)}); code != 409 {
		t.Fatalf("adjust FT-02 inside maintenance must be 409, got %d %v", code, out)
	}
	if len(asList(out["details"])) == 0 {
		t.Fatal("409 adjust must list per-item conflicts")
	}

	// Audit log was written for write operations.
	if code, out = c.do(http.MethodGet, "/api/audit-logs?page_size=50", nil); code != 200 {
		t.Fatalf("audit logs: %d", code)
	}
	found := false
	for _, item := range asList(out["items"]) {
		row := asMap(item)
		if msg, _ := row["action"].(string); contains(msg, "登记延误") {
			found = true
		}
	}
	if !found {
		t.Fatal("audit log missing delay registration entry")
	}
}

func taskByType2(detail map[string]interface{}, taskType string) map[string]interface{} {
	for _, item := range asList(detail["tasks"]) {
		row := asMap(item)
		if row["task_type"] == taskType {
			return row
		}
	}
	return nil
}

func contains(s, sub string) bool {
	return bytes.Contains([]byte(s), []byte(sub))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	out := ""
	for n > 0 {
		out = string(rune('0'+n%10)) + out
		n /= 10
	}
	return out
}

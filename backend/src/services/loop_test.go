package services

import (
	"path/filepath"
	"testing"
	"time"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/database"
	"groundTurn/src/models"
	"groundTurn/src/types"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "loop.db")
	cfg := &config.Config{SQLitePath: path, JWTSecret: "test-secret"}
	db, err := database.Open(cfg)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db
}

func seedTestResources(t *testing.T, db *gorm.DB) {
	rows := []models.GroundResource{
		{ResourceCode: "BELT-01", ResourceType: "BAGGAGE", AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-BAG"},
		{ResourceCode: "BELT-02", ResourceType: "BAGGAGE", AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-BAG"},
		{ResourceCode: "CART-01", ResourceType: "CLEANING", AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-CLEAN"},
		{ResourceCode: "CART-02", ResourceType: "CLEANING", AvailabilityStatus: "MAINTENANCE", OwnerTeam: "TEAM-CLEAN"},
		{ResourceCode: "CATER-01", ResourceType: "CATERING", AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-CATER"},
		{ResourceCode: "CATER-02", ResourceType: "CATERING", AvailabilityStatus: "OFFLINE", OwnerTeam: "TEAM-CATER"},
		{ResourceCode: "WATER-01", ResourceType: "WATER_SERVICE", AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-WATER"},
		{ResourceCode: "FUEL-01", ResourceType: "REFUEL", AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-FUEL"},
		{ResourceCode: "PUSH-01", ResourceType: "PUSHBACK", AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-RAMP"},
	}
	for i := range rows {
		if err := db.Create(&rows[i]).Error; err != nil {
			t.Fatalf("seed resource: %v", err)
		}
	}
}

var dispatcher = Actor{Username: "dispatcher", Role: "DISPATCHER", Name: "调度"}

func TestGeneratePlanClosedLoop(t *testing.T) {
	db := setupTestDB(t)
	seedTestResources(t, db)
	svc := NewServices(db, &config.Config{JWTSecret: "test-secret"})

	base := time.Date(2026, 9, 21, 8, 0, 0, 0, time.Local)

	// 航班 A：到站时间 08:00，75 分钟过站。
	flightA, err := svc.Turnarounds.Create(dispatcher, types.CreateTurnaroundRequest{
		FlightNo: "CA1831", AircraftReg: "B-5821", StandNo: "203",
		ArrivalTime: base, DepartureTime: base.Add(75 * time.Minute),
	})
	if err != nil {
		t.Fatalf("create flight A: %v", err)
	}

	result, err := svc.Turnarounds.GeneratePlan(dispatcher, flightA.ID)
	if err != nil {
		t.Fatalf("generate plan A: %v", err)
	}
	if !result.Saved || result.TaskCount != 6 {
		t.Fatalf("expected 6 saved tasks, got saved=%v count=%d items=%+v", result.Saved, result.TaskCount, result.Items)
	}

	var bookings []models.ResourceBooking
	if err := db.Where("turnaround_id = ?", flightA.ID).Find(&bookings).Error; err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 6 {
		t.Fatalf("expected 6 bookings, got %d", len(bookings))
	}
	for _, b := range bookings {
		if b.BookingStatus != "CONFIRMED" {
			t.Fatalf("booking %d not confirmed: %s", b.ID, b.BookingStatus)
		}
	}

	// 航班 C：只有 40 分钟过站，推送任务截止 09:10 > 离港 08:40，整批拒绝。
	flightC, err := svc.Turnarounds.Create(dispatcher, types.CreateTurnaroundRequest{
		FlightNo: "CZ3108", AircraftReg: "B-9941", StandNo: "305",
		ArrivalTime: base.Add(2 * time.Hour), DepartureTime: base.Add(2*time.Hour + 40*time.Minute),
	})
	if err != nil {
		t.Fatalf("create flight C: %v", err)
	}
	rejected, err := svc.Turnarounds.GeneratePlan(dispatcher, flightC.ID)
	if err != nil {
		t.Fatalf("generate plan C: %v", err)
	}
	if rejected.Saved {
		t.Fatal("flight C plan should be rejected")
	}
	if len(rejected.Items) == 0 {
		t.Fatal("expected per-item conflict reasons")
	}
	foundDeadline := false
	for _, item := range rejected.Items {
		if item.Code == constants.ItemDeadlineConflict {
			foundDeadline = true
		}
	}
	if !foundDeadline {
		t.Fatalf("expected deadline conflict item, got %+v", rejected.Items)
	}
	var cTasks int64
	db.Model(&models.GroundTask{}).Where("turnaround_id = ?", flightC.ID).Count(&cTasks)
	var cBookings int64
	db.Model(&models.ResourceBooking{}).Where("turnaround_id = ?", flightC.ID).Count(&cBookings)
	if cTasks != 0 || cBookings != 0 {
		t.Fatalf("rejected batch must not persist anything, tasks=%d bookings=%d", cTasks, cBookings)
	}
}

func TestDelayReplanKeepsSignedAndBumpsBooking(t *testing.T) {
	db := setupTestDB(t)
	seedTestResources(t, db)
	svc := NewServices(db, &config.Config{JWTSecret: "test-secret"})

	base := time.Date(2026, 9, 21, 8, 0, 0, 0, time.Local)

	// 预先占用 CATER-01 在 08:50~09:20：延误 +40 分钟后配餐窗口正好撞上。
	blockerFlight := &models.FlightTurnaround{
		FlightNo: "HU7801", AircraftReg: "B-3217", StandNo: "201",
		ArrivalTime: base.Add(-2 * time.Hour), DepartureTime: base.Add(-1 * time.Hour),
		TurnaroundStatus: "IN_SERVICE", PlanGenerated: true,
	}
	if err := db.Create(blockerFlight).Error; err != nil {
		t.Fatal(err)
	}
	blockerTask := &models.GroundTask{
		TurnaroundID: blockerFlight.ID, TaskType: "CATERING", TeamID: "TEAM-CATER",
		PlannedStart: base.Add(50 * time.Minute), PlannedEnd: base.Add(80 * time.Minute),
		Deadline: base.Add(85 * time.Minute), Status: "SIGNED",
	}
	if err := db.Create(blockerTask).Error; err != nil {
		t.Fatal(err)
	}
	blockerBooking := &models.ResourceBooking{
		ResourceID: 5, TurnaroundID: blockerFlight.ID, TaskID: &blockerTask.ID,
		StartTime: base.Add(50 * time.Minute), EndTime: base.Add(80 * time.Minute),
		BookingStatus: "CONFIRMED",
	}
	if err := db.Create(blockerBooking).Error; err != nil {
		t.Fatal(err)
	}

	flight, err := svc.Turnarounds.Create(dispatcher, types.CreateTurnaroundRequest{
		FlightNo: "MU5102", AircraftReg: "B-6690", StandNo: "118",
		ArrivalTime: base, DepartureTime: base.Add(75 * time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Turnarounds.GeneratePlan(dispatcher, flight.ID); err != nil {
		t.Fatal(err)
	}

	tasks, _ := svc.Tasks.ListByTurnaround(flight.ID)
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}
	findTask := func(kind string) *models.GroundTask {
		for i := range tasks {
			if tasks[i].TaskType == kind {
				return &tasks[i]
			}
		}
		t.Fatalf("task %s not found", kind)
		return nil
	}
	baggage := findTask("BAGGAGE")
	catering := findTask("CATERING")

	// 签收行李任务；其计划时点必须在延误后保持不变。
	if _, err := svc.Tasks.Sign(dispatcher, baggage.ID, types.SignTaskRequest{Operator: "周海"}); err != nil {
		t.Fatalf("sign baggage: %v", err)
	}
	originalBaggageStart := baggage.PlannedStart

	replan, err := svc.Delays.Register(dispatcher, flight.ID, types.RegisterDelayRequest{
		DelayType: "ATC", Minutes: 40, RootCause: "出港流控", ResponsibilityTeam: "TEAM-OPS",
	})
	if err != nil {
		t.Fatalf("register delay: %v", err)
	}
	if len(replan.KeptTasks) != 1 {
		t.Fatalf("expected 1 kept (signed) task, got %d", len(replan.KeptTasks))
	}
	if !replan.KeptTasks[0].PlannedStart.Equal(originalBaggageStart) {
		t.Fatalf("signed task moved: %v -> %v", originalBaggageStart, replan.KeptTasks[0].PlannedStart)
	}
	if len(replan.MovedTasks) != 5 {
		t.Fatalf("expected 5 moved tasks, got %d", len(replan.MovedTasks))
	}
	if len(replan.Bumped) == 0 {
		t.Fatal("expected catering booking bumped to PENDING")
	}

	// 刷新后从数据库读取：已签收任务时点不变；配餐预约 PENDING 且带原因。
	updatedTasks, _ := svc.Tasks.ListByTurnaround(flight.ID)
	var updatedBaggage, updatedCatering models.GroundTask
	for _, task := range updatedTasks {
		if task.ID == baggage.ID {
			updatedBaggage = task
		}
		if task.ID == catering.ID {
			updatedCatering = task
		}
	}
	if !updatedBaggage.PlannedStart.Equal(originalBaggageStart) {
		t.Fatalf("persisted signed task changed: %v", updatedBaggage.PlannedStart)
	}
	if !updatedCatering.PlannedStart.Equal(catering.PlannedStart.Add(40 * time.Minute)) {
		t.Fatalf("unsigned task not shifted: %v", updatedCatering.PlannedStart)
	}

	var bumpedBooking models.ResourceBooking
	db.Where("task_id = ? AND booking_status = ?", catering.ID, "PENDING").First(&bumpedBooking)
	if bumpedBooking.ID == 0 {
		t.Fatal("bumped booking not persisted as PENDING")
	}
	if bumpedBooking.ConflictReason == "" {
		t.Fatal("bumped booking must carry conflict_reason")
	}

	// 调度在页面上调整时段到空闲窗口 → 重新确认。
	adjusted, err := svc.Bookings.Adjust(dispatcher, bumpedBooking.ID, types.AdjustBookingRequest{
		ResourceID: 5,
		StartTime:  base.Add(100 * time.Minute),
		EndTime:    base.Add(120 * time.Minute),
	})
	if err != nil {
		t.Fatalf("adjust booking: %v", err)
	}
	if adjusted.BookingStatus != "CONFIRMED" || adjusted.ConflictReason != "" {
		t.Fatalf("adjusted booking should be confirmed, got %s/%s", adjusted.BookingStatus, adjusted.ConflictReason)
	}
}

func TestGeneratePlanRejectsMaintenanceAndOfflineResources(t *testing.T) {
	db := setupTestDB(t)
	// 只保留一台维护中的清洁设备、一台离线配餐设备，无可用替代。
	db.Create(&models.GroundResource{ResourceCode: "CART-09", ResourceType: "CLEANING", AvailabilityStatus: "MAINTENANCE", OwnerTeam: "TEAM-CLEAN"})
	db.Create(&models.GroundResource{ResourceCode: "CATER-09", ResourceType: "CATERING", AvailabilityStatus: "OFFLINE", OwnerTeam: "TEAM-CATER"})
	for _, def := range []struct{ code, typ string }{
		{"BELT-09", "BAGGAGE"}, {"WATER-09", "WATER_SERVICE"}, {"FUEL-09", "REFUEL"}, {"PUSH-09", "PUSHBACK"},
	} {
		db.Create(&models.GroundResource{ResourceCode: def.code, ResourceType: def.typ, AvailabilityStatus: "AVAILABLE", OwnerTeam: "TEAM-X"})
	}
	svc := NewServices(db, &config.Config{JWTSecret: "test-secret"})
	base := time.Date(2026, 9, 21, 8, 0, 0, 0, time.Local)
	flight, err := svc.Turnarounds.Create(dispatcher, types.CreateTurnaroundRequest{
		FlightNo: "CA1839", AircraftReg: "B-1111", StandNo: "210",
		ArrivalTime: base, DepartureTime: base.Add(75 * time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := svc.Turnarounds.GeneratePlan(dispatcher, flight.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.Saved {
		t.Fatal("expected whole batch rejected")
	}
	codes := map[string]bool{}
	for _, item := range result.Items {
		codes[item.Code] = true
	}
	if !codes[constants.ItemResourceMaintenance] || !codes[constants.ItemResourceOffline] {
		t.Fatalf("expected maintenance/offline reasons, got %+v", codes)
	}
}

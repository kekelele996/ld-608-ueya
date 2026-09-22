package services

import (
	"testing"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/utils"
)

// TestClosedLoop runs the full operational loop:
//  1. plan generation is rejected ALL-OR-NOTHING with itemized conflicts;
//  2. after freeing blockers, the same plan saves atomically;
//  3. delay registration reschedules unsigned tasks, freezes signed ones,
//     and squeezes conflicting bookings into PENDING;
//  4. a pending booking can be adjusted to a clean window and confirmed,
//     while adjusting into a maintenance window keeps it PENDING;
//  5. task sign-off survives a second delay (time point frozen).
func TestClosedLoop(t *testing.T) {
	f := setupFixture(t)

	// 1. Rejected batch: belt loader overlap+offline, fuel maintenance+overlap.
	result, err := f.plan.GeneratePlan(f.caID, nil, "dispatcher", constants.RoleDispatcher)
	if err == nil {
		t.Fatalf("expected conflict error, got saved result: %+v", result)
	}
	apiErr, ok := err.(*utils.APIError)
	if !ok || apiErr.HTTPCode != 409 {
		t.Fatalf("expected 409 APIError, got %#v", err)
	}
	if result == nil || result.Saved {
		t.Fatalf("rejected batch must report saved=false")
	}
	details, ok := apiErr.Details.([]typesConflictItem)
	if !ok {
		t.Fatalf("expected conflict details []ConflictItem, got %T", apiErr.Details)
	}
	reasonCount := map[string]int{}
	for _, item := range details {
		reasonCount[item.Code]++
	}
	// Belt loader: BL-01 overlaps MU202 signed booking, BL-02 is offline.
	// Fuel truck: FT-01 is under maintenance, FT-02 offline, FT-03 occupied.
	if reasonCount[constants.ConflictOverlap] < 2 {
		t.Fatalf("expected >=2 overlap reasons (BL-01, FT-03), got %v", reasonCount)
	}
	if reasonCount[constants.ConflictOffline] < 2 {
		t.Fatalf("expected 2 offline reasons (BL-02, FT-02), got %v", reasonCount)
	}
	if reasonCount[constants.ConflictMaintenance] < 1 {
		t.Fatalf("FT-01 maintenance collides with fuel task; expected reason, got %v", reasonCount)
	}

	// Nothing persisted for CA101: tasks and bookings must roll back.
	tasks, _ := f.plan.Task.ListByTurnaround(f.caID)
	if len(tasks) != 0 {
		t.Fatalf("rejected batch must rollback tasks, found %d", len(tasks))
	}
	bookings, _ := f.plan.Booking.ListByTurnaround(f.caID)
	if len(bookings) != 0 {
		t.Fatalf("rejected batch must rollback bookings, found %d", len(bookings))
	}

	// 2. Free blockers:
	//    - bring BL-02 online (second belt loader, free window);
	//    - clear FT-03 (the parked blocker is removed below in the live flow
	//      by taking it out of the scenario) — instead we keep FT-03 blocked
	//      all day and switch FT-02 to MAINTENANCE across the shifted fuel
	//      window only. FT-02 remains usable for the initial plan but rejects
	//      the +30m shifted fuel window, forcing the booking to PENDING.
	bl02ID := f.resources["BL-02"]
	var bl02 models.GroundResource
	if err := f.db.First(&bl02, bl02ID).Error; err != nil {
		t.Fatalf("load BL-02: %v", err)
	}
	bl02.AvailabilityStatus = constants.ResourceAvailable
	if err := f.db.Save(&bl02).Error; err != nil {
		t.Fatalf("online BL-02: %v", err)
	}
	// FT-02: release MU202's old booking and put FT-02 under maintenance
	// starting at 10:40 (after CA101's initial 10:12-10:37 fuel window but
	// covering its shifted 10:42-11:07 window).
	muBookings, _ := f.plan.Booking.ListByTurnaround(f.muID)
	for _, b := range muBookings {
		if b.ResourceID == f.resources["FT-02"] {
			b.BookingStatus = constants.BookingReleased
			if err := f.plan.Booking.Save(&b); err != nil {
				t.Fatalf("release FT-02 booking: %v", err)
			}
		}
	}
	ft02ID := f.resources["FT-02"]
	var ft02 models.GroundResource
	if err := f.db.First(&ft02, ft02ID).Error; err != nil {
		t.Fatalf("load FT-02: %v", err)
	}
	ft02.AvailabilityStatus = constants.ResourceAvailable
	mStart := f.arrival.Add(40 * time.Minute)
	mEnd := f.arrival.Add(75 * time.Minute)
	ft02.MaintenanceDueAt = &mStart
	ft02.MaintenanceEnd = &mEnd
	if err := f.db.Save(&ft02).Error; err != nil {
		t.Fatalf("maintenance FT-02: %v", err)
	}
	// FT-03 stays blocked all day by the seeded CZ303 confirmed booking.

	result, err = f.plan.GeneratePlan(f.caID, nil, "dispatcher", constants.RoleDispatcher)
	if err != nil {
		t.Fatalf("second plan should save, got %v", err)
	}
	if !result.Saved || len(result.Tasks) != 6 || len(result.Bookings) != 6 {
		t.Fatalf("expected 6 tasks/bookings saved, got saved=%v tasks=%d bookings=%d",
			result.Saved, len(result.Tasks), len(result.Bookings))
	}
	for _, b := range result.Bookings {
		if b.BookingStatus != constants.BookingConfirmed {
			t.Fatalf("fresh booking should be CONFIRMED, got %s (%s)", b.BookingStatus, b.ConflictCode)
		}
	}

	// 3. Register +30 minutes:
	//    CLEANING shifts 10:05-10:35 -> 10:35-11:05 and collides with CZ303's
	//    SIGNED CK-02 window 10:35-11:05 -> CK-02 booking PENDING (overlap).
	//    REFUEL shifts 10:12-10:37 -> 10:42-11:07; the assigned truck FT-02
	//    enters maintenance at 10:40 -> its booking becomes PENDING.
	res, err := f.delay.RegisterDelay(f.caID, constants.DelayLateArrival, 30,
		"前站晚到", "AIRLINE", "dispatcher", constants.RoleDispatcher)
	if err != nil {
		t.Fatalf("register delay: %v", err)
	}
	if res.RescheduledTasks != 6 {
		t.Fatalf("expected 6 rescheduled tasks, got %d", res.RescheduledTasks)
	}
	if res.FrozenTasks != 0 {
		t.Fatalf("CA101 had no signed tasks, frozen=%d", res.FrozenTasks)
	}
	if res.PendingBookings < 2 {
		t.Fatalf("expected >=2 pending bookings, got %d", res.PendingBookings)
	}
	caTasks, _ := f.plan.Task.ListByTurnaround(f.caID)
	baggage := findTask(caTasks, constants.TaskTypeBaggage)
	if baggage.PlannedStart != f.arrival.Add(32*time.Minute) {
		t.Fatalf("baggage expected 10:32 after +30 shift, got %s", baggage.PlannedStart)
	}
	caBookings, _ := f.plan.Booking.ListByTurnaround(f.caID)
	pendingByCode := map[string]string{}
	for _, b := range caBookings {
		if b.BookingStatus == constants.BookingPending {
			pendingByCode[f.codeByResourceID(b.ResourceID)] = b.ConflictReason
		}
	}
	if pendingByCode["CK-02"] != constants.ConflictOverlap {
		t.Fatalf("CK-02 booking should be PENDING overlap, got %v", pendingByCode)
	}
	if pendingByCode["FT-02"] != constants.ConflictMaintenance {
		t.Fatalf("FT-02 booking should be PENDING maintenance, got %v", pendingByCode)
	}

	// MU202 signed baggage keeps its original window (10:02 arrival +2m).
	muTasks, _ := f.plan.Task.ListByTurnaround(f.muID)
	muBaggage := findTask(muTasks, constants.TaskTypeBaggage)
	if !muBaggage.PlannedStart.Equal(f.arrival.Add(-3 * time.Minute)) {
		t.Fatalf("signed MU202 baggage must freeze at 09:57, got %s", muBaggage.PlannedStart)
	}

	// Adjust pending CK-02 booking to 11:10-11:25: within the delayed
	// departure 11:30 and clear of CZ303's signed 10:35-11:05 window.
	var ck02Pending uint
	for _, b := range caBookings {
		if b.ResourceID == f.resources["CK-02"] && b.BookingStatus == constants.BookingPending {
			ck02Pending = b.ID
		}
	}
	newStart := f.arrival.Add(70 * time.Minute) // 11:10
	newEnd := f.arrival.Add(85 * time.Minute)   // 11:25
	booking, conflicts, err := f.booking.Adjust(ck02Pending, newStart, newEnd,
		"resource", constants.RoleResourceManager)
	if err != nil {
		t.Fatalf("adjust CK-02 should confirm, got %v conflicts=%v", err, conflicts)
	}
	if booking.BookingStatus != constants.BookingConfirmed {
		t.Fatalf("adjusted booking should confirm, got %s", booking.BookingStatus)
	}

	// Adjusting pending FT-02 into its maintenance window must fail and the
	// booking must remain PENDING with itemized reasons.
	var ftPending uint
	caBookings2, _ := f.plan.Booking.ListByTurnaround(f.caID)
	for _, b := range caBookings2 {
		if b.ResourceID == f.resources["FT-02"] && b.BookingStatus == constants.BookingPending {
			ftPending = b.ID
		}
	}
	_, badConflicts, err := f.booking.Adjust(ftPending,
		f.arrival.Add(40*time.Minute), f.arrival.Add(60*time.Minute),
		"resource", constants.RoleResourceManager)
	if err == nil {
		t.Fatalf("adjusting FT-02 inside maintenance must fail")
	}
	if len(badConflicts) == 0 {
		t.Fatalf("failed adjust must return itemized conflict reasons")
	}
	stillPending, _ := f.plan.Booking.Get(ftPending)
	if stillPending.BookingStatus != constants.BookingPending {
		t.Fatalf("failed adjust must keep booking PENDING, got %s", stillPending.BookingStatus)
	}

	// Move it to 11:15-11:25 after maintenance ends; FT-02 itself is clear.
	_, _, err = f.booking.Adjust(ftPending,
		f.arrival.Add(75*time.Minute), f.arrival.Add(90*time.Minute),
		"resource", constants.RoleResourceManager)
	if err != nil {
		t.Fatalf("FT-02 adjust to 11:15-11:30 should confirm, got %v", err)
	}

	// 5. Sign CA101 baggage then delay again +10m: baggage stays at 10:32.
	if _, err := f.task.Accept(baggage.ID, "team-bag", constants.RoleTeam, "TEAM-BAG"); err != nil {
		t.Fatalf("accept baggage: %v", err)
	}
	res2, err := f.delay.RegisterDelay(f.caID, constants.DelayGround, 10,
		"二次延误", "GROUND", "dispatcher", constants.RoleDispatcher)
	if err != nil {
		t.Fatalf("second delay: %v", err)
	}
	caTasks2, _ := f.plan.Task.ListByTurnaround(f.caID)
	baggage2 := findTask(caTasks2, constants.TaskTypeBaggage)
	if !baggage2.PlannedStart.Equal(f.arrival.Add(32 * time.Minute)) {
		t.Fatalf("signed CA101 baggage must freeze at 10:32, got %s", baggage2.PlannedStart)
	}
	if res2.FrozenTasks != 1 {
		t.Fatalf("expected 1 frozen signed task, got %d", res2.FrozenTasks)
	}
}

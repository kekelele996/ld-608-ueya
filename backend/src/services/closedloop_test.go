package services

import (
	"testing"
	"time"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_busy_timeout=5000"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := config.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	// Shared in-memory DB persists across handles in same process; clean it.
	db.Exec("DELETE FROM resource_booking")
	db.Exec("DELETE FROM ground_task")
	db.Exec("DELETE FROM delay_event")
	db.Exec("DELETE FROM ground_resource")
	db.Exec("DELETE FROM flight_turnaround")
	db.Exec("DELETE FROM audit_log")
	db.Exec("DELETE FROM app_user")
	return db
}

// seedFixture builds resources and two turnarounds mirroring config.Seed but
// relative to the test run time.
type fixture struct {
	db        *gorm.DB
	plan      *PlanService
	delay     *DelayService
	task      *TaskService
	booking   *BookingService
	query     *QueryService
	arrival   time.Time
	caID      uint
	muID      uint
	czID      uint
	resources map[string]uint
}

func setupFixture(t *testing.T) *fixture {
	db := testDB(t)
	base := time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC)
	arrival := base

	resources := []models.GroundResource{
		{ResourceCode: "CK-01", ResourceType: "CLEANING_KIT", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CLEAN"},
		{ResourceCode: "CK-02", ResourceType: "CLEANING_KIT", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CLEAN"},
		{ResourceCode: "CT-01", ResourceType: "CATERING_TRUCK", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CATER"},
		{ResourceCode: "CT-02", ResourceType: "CATERING_TRUCK", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CATER"},
		{ResourceCode: "BL-01", ResourceType: "BELT_LOADER", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-BAG"},
		{ResourceCode: "BL-02", ResourceType: "BELT_LOADER", AvailabilityStatus: constants.ResourceOffline, OwnerTeam: "TEAM-BAG"},
		{ResourceCode: "FT-01", ResourceType: "FUEL_TRUCK", AvailabilityStatus: constants.ResourceMaintenance, OwnerTeam: "TEAM-FUEL",
			MaintenanceDueAt: pTime(base.Add(5 * time.Minute)), MaintenanceEnd: pTime(base.Add(70 * time.Minute))},
		{ResourceCode: "FT-02", ResourceType: "FUEL_TRUCK", AvailabilityStatus: constants.ResourceOffline, OwnerTeam: "TEAM-FUEL"},
		{ResourceCode: "FT-03", ResourceType: "FUEL_TRUCK", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-FUEL"},
		{ResourceCode: "WT-01", ResourceType: "WATER_TRUCK", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-WATER"},
		{ResourceCode: "WT-02", ResourceType: "WATER_TRUCK", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-WATER"},
		{ResourceCode: "PT-01", ResourceType: "PUSHBACK_TRACTOR", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-RAMP"},
		{ResourceCode: "PT-02", ResourceType: "PUSHBACK_TRACTOR", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-RAMP"},
	}
	if err := db.Create(&resources).Error; err != nil {
		t.Fatalf("resources: %v", err)
	}
	resID := map[string]uint{}
	for _, r := range resources {
		resID[r.ResourceCode] = r.ID
	}

	ca := models.FlightTurnaround{FlightNo: "CA101", AircraftReg: "B-1", StandNo: "203",
		ArrivalTime: base, DepartureTime: base.Add(60 * time.Minute), TurnaroundStatus: constants.StatusOnStand}
	db.Create(&ca)

	muBase := base.Add(-5 * time.Minute)
	mu := models.FlightTurnaround{FlightNo: "MU202", AircraftReg: "B-2", StandNo: "204",
		ArrivalTime: muBase, DepartureTime: muBase.Add(55 * time.Minute), TurnaroundStatus: constants.StatusInService}
	db.Create(&mu)
	specs := []struct {
		typ    string
		team   string
		code   string
		off    int
		dur    int
		status string
	}{
		{constants.TaskTypeBaggage, "TEAM-BAG", "BL-01", 2, 25, constants.TaskAccepted},
		{constants.TaskTypeCleaning, "TEAM-CLEAN", "CK-01", 5, 30, constants.TaskPlanned},
		{constants.TaskTypeWaterService, "TEAM-WATER", "WT-01", 8, 20, constants.TaskPlanned},
		{constants.TaskTypeCatering, "TEAM-CATER", "CT-01", 10, 25, constants.TaskPlanned},
		{constants.TaskTypeRefuel, "TEAM-FUEL", "FT-02", 12, 25, constants.TaskPlanned},
		{constants.TaskTypePushback, "TEAM-RAMP", "PT-01", 42, 12, constants.TaskPlanned},
	}
	for _, sp := range specs {
		task := models.GroundTask{TurnaroundID: mu.ID, TaskType: sp.typ, TeamID: sp.team,
			PlannedStart: muBase.Add(time.Duration(sp.off) * time.Minute),
			Deadline:     muBase.Add(time.Duration(sp.off+sp.dur) * time.Minute),
			Status:       sp.status}
		if sp.status == constants.TaskAccepted {
			task.SignedAt = pTime(task.PlannedStart)
		}
		db.Create(&task)
		db.Create(&models.ResourceBooking{ResourceID: resID[sp.code], TurnaroundID: mu.ID, TaskID: task.ID,
			StartTime: task.PlannedStart, EndTime: task.Deadline, BookingStatus: constants.BookingConfirmed})
	}

	czBase := base.Add(35 * time.Minute)
	cz := models.FlightTurnaround{FlightNo: "CZ303", AircraftReg: "B-3", StandNo: "118",
		ArrivalTime: czBase, DepartureTime: czBase.Add(60 * time.Minute), TurnaroundStatus: constants.StatusInService}
	db.Create(&cz)
	czTask := models.GroundTask{TurnaroundID: cz.ID, TaskType: constants.TaskTypeCleaning, TeamID: "TEAM-CLEAN",
		PlannedStart: czBase, Deadline: czBase.Add(30 * time.Minute), Status: constants.TaskAccepted,
		SignedAt: pTime(czBase)}
	db.Create(&czTask)
	db.Create(&models.ResourceBooking{ResourceID: resID["CK-02"], TurnaroundID: cz.ID, TaskID: czTask.ID,
		StartTime: czBase, EndTime: czBase.Add(30 * time.Minute), BookingStatus: constants.BookingConfirmed})

	// FT-03 (third fuel truck) is blocked across the whole CA101 window by a
	// confirmed CZ303 booking, so rejected-batch validation reports overlap
	// for it as well.
	db.Create(&models.ResourceBooking{ResourceID: resID["FT-03"], TurnaroundID: cz.ID,
		StartTime: base.Add(10 * time.Minute), EndTime: base.Add(80 * time.Minute),
		BookingStatus: constants.BookingConfirmed})

	return &fixture{
		db:        db,
		plan:      NewPlanService(db),
		delay:     NewDelayService(db, NewPlanService(db)),
		task:      NewTaskService(db),
		booking:   NewBookingService(db),
		query:     NewQueryService(db),
		arrival:   arrival,
		caID:      ca.ID,
		muID:      mu.ID,
		czID:      cz.ID,
		resources: resID,
	}
}

func pTime(t time.Time) *time.Time { return &t }

package config

import (
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"

	"gorm.io/gorm"
)

// Seed inserts local-only demo data covering the full closed loop:
//   - flight CA101: no plan yet; first generation is REJECTED batch with
//     itemized conflicts (belt loader + fuel truck);
//   - flight MU202: full plan already generated, overlapping CA101 windows,
//     with one signed task that must freeze on delay;
//	 - flight CZ303: one signed task to demonstrate cross-flight freeze.
//
// Every timestamp is fixed so conflict demos are reproducible.
func Seed(db *gorm.DB) error {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		return nil
	}

	loc := time.Local
	base := time.Date(2026, 9, 22, 10, 0, 0, 0, loc) // CA101 arrival
	at := func(min int) time.Time { return base.Add(time.Duration(min) * time.Minute) }

	users := []models.User{
		{Username: "dispatcher", DisplayName: "调度-王敏", Role: constants.RoleDispatcher},
		{Username: "team-bag", DisplayName: "行李班长-李强", Role: constants.RoleTeam, TeamID: "TEAM-BAG"},
		{Username: "team-fuel", DisplayName: "加油班长-赵雷", Role: constants.RoleTeam, TeamID: "TEAM-FUEL"},
		{Username: "team-ramp", DisplayName: "机坪班长-孙伟", Role: constants.RoleTeam, TeamID: "TEAM-RAMP"},
		{Username: "resource", DisplayName: "资源管理员-周倩", Role: constants.RoleResourceManager},
		{Username: "supervisor", DisplayName: "运行督导-陈刚", Role: constants.RoleSupervisor},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	ptrTime := func(t time.Time) *time.Time { return &t }

	resources := []models.GroundResource{
		// CLEANING_KIT x2 available
		{ResourceCode: "CK-01", ResourceType: "CLEANING_KIT", Location: "T2 保洁间", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CLEAN"},
		{ResourceCode: "CK-02", ResourceType: "CLEANING_KIT", Location: "T2 保洁间", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CLEAN"},
		// CATERING_TRUCK x2
		{ResourceCode: "CT-01", ResourceType: "CATERING_TRUCK", Location: "航食通道", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CATER"},
		{ResourceCode: "CT-02", ResourceType: "CATERING_TRUCK", Location: "航食通道", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-CATER"},
		// BELT_LOADER x2: BL-01 occupied at CA101 BAGGAGE window, BL-02 offline
		{ResourceCode: "BL-01", ResourceType: "BELT_LOADER", Location: "行李分拣区", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-BAG"},
		{ResourceCode: "BL-02", ResourceType: "BELT_LOADER", Location: "行李分拣区", AvailabilityStatus: constants.ResourceOffline, OwnerTeam: "TEAM-BAG"},
		// FUEL_TRUCK x3:
		// FT-01 maintenance covers both initial (10:12) and +30m (10:42)
		// fuel windows; FT-02 is offline (with a pre-registered maintenance
		// window 10:40-11:15 that bites after it comes back online and the
		// flight is delayed); FT-03 is blocked all morning by CZ303.
		{ResourceCode: "FT-01", ResourceType: "FUEL_TRUCK", Location: "油站", AvailabilityStatus: constants.ResourceMaintenance,
			OwnerTeam: "TEAM-FUEL", MaintenanceDueAt: ptrTime(at(5)), MaintenanceEnd: ptrTime(at(70))},
		{ResourceCode: "FT-02", ResourceType: "FUEL_TRUCK", Location: "油站", AvailabilityStatus: constants.ResourceOffline,
			OwnerTeam: "TEAM-FUEL", MaintenanceDueAt: ptrTime(at(40)), MaintenanceEnd: ptrTime(at(75))},
		{ResourceCode: "FT-03", ResourceType: "FUEL_TRUCK", Location: "油站", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-FUEL"},
		// WATER_TRUCK
		{ResourceCode: "WT-01", ResourceType: "WATER_TRUCK", Location: "清水站", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-WATER"},
		{ResourceCode: "WT-02", ResourceType: "WATER_TRUCK", Location: "清水站", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-WATER"},
		// PUSHBACK_TRACTOR
		{ResourceCode: "PT-01", ResourceType: "PUSHBACK_TRACTOR", Location: "机坪 A 区", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-RAMP"},
		{ResourceCode: "PT-02", ResourceType: "PUSHBACK_TRACTOR", Location: "机坪 B 区", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "TEAM-RAMP"},
	}
	if err := db.Create(&resources).Error; err != nil {
		return err
	}
	resID := map[string]uint{}
	for _, r := range resources {
		resID[r.ResourceCode] = r.ID
	}

	// CA101: arrived ON_STAND, no tasks yet. Generate plan -> rejected batch.
	ca101 := models.FlightTurnaround{
		FlightNo: "CA101", AircraftReg: "B-6271", StandNo: "203",
		ArrivalTime: at(0), DepartureTime: at(60),
		TurnaroundStatus: constants.StatusOnStand,
	}
	if err := db.Create(&ca101).Error; err != nil {
		return err
	}

	// MU202: another turnaround whose plan occupies BL-01 / FT-02 / WT-01 /
	// PT-01 during CA101 windows. Arrival 09:55, departure 10:50. Its baggage
	// task is already signed (frozen), so delay on CA101 cannot move it.
	muBase := base.Add(-5 * time.Minute)
	muAt := func(min int) time.Time { return muBase.Add(time.Duration(min) * time.Minute) }
	mu202 := models.FlightTurnaround{
		FlightNo: "MU202", AircraftReg: "B-1088", StandNo: "204",
		ArrivalTime: muAt(0), DepartureTime: muAt(55),
		TurnaroundStatus: constants.StatusInService,
	}
	if err := db.Create(&mu202).Error; err != nil {
		return err
	}

	// MU202 tasks mirror offsets: BAGGAGE 2-27, CLEAN 5-35, WATER 8-28,
	// CATER 10-35, FUEL 12-37, PUSHBACK 42-54.
	mkTask := func(turnID uint, taskType, team string, off, dur int, status string) models.GroundTask {
		t := models.GroundTask{
			TurnaroundID: turnID, TaskType: taskType, TeamID: team,
			PlannedStart: muAt(off), Deadline: muAt(off + dur),
			Status: status,
		}
		if status == constants.TaskAccepted || status == constants.TaskCompleted {
			signed := muAt(off)
			t.SignedAt = &signed
		}
		return t
	}
	muTasks := []models.GroundTask{
		mkTask(mu202.ID, constants.TaskTypeBaggage, "TEAM-BAG", 2, 25, constants.TaskAccepted),
		mkTask(mu202.ID, constants.TaskTypeCleaning, "TEAM-CLEAN", 5, 30, constants.TaskPlanned),
		mkTask(mu202.ID, constants.TaskTypeWaterService, "TEAM-WATER", 8, 20, constants.TaskPlanned),
		mkTask(mu202.ID, constants.TaskTypeCatering, "TEAM-CATER", 10, 25, constants.TaskPlanned),
		mkTask(mu202.ID, constants.TaskTypeRefuel, "TEAM-FUEL", 12, 25, constants.TaskPlanned),
		mkTask(mu202.ID, constants.TaskTypePushback, "TEAM-RAMP", 42, 12, constants.TaskPlanned),
	}
	if err := db.Create(&muTasks).Error; err != nil {
		return err
	}
	taskID := map[string]uint{}
	for _, t := range muTasks {
		taskID[t.TaskType] = t.ID
	}

	mkBooking := func(resourceCode string, taskID uint, off, dur int, status string) models.ResourceBooking {
		return models.ResourceBooking{
			ResourceID: resID[resourceCode], TurnaroundID: mu202.ID, TaskID: taskID,
			StartTime: muAt(off), EndTime: muAt(off + dur), BookingStatus: status,
		}
	}
	muBookings := []models.ResourceBooking{
		mkBooking("BL-01", taskID[constants.TaskTypeBaggage], 2, 25, constants.BookingConfirmed),
		mkBooking("CK-01", taskID[constants.TaskTypeCleaning], 5, 30, constants.BookingConfirmed),
		mkBooking("WT-01", taskID[constants.TaskTypeWaterService], 8, 20, constants.BookingConfirmed),
		mkBooking("CT-01", taskID[constants.TaskTypeCatering], 10, 25, constants.BookingConfirmed),
		mkBooking("FT-02", taskID[constants.TaskTypeRefuel], 12, 25, constants.BookingConfirmed),
		mkBooking("PT-01", taskID[constants.TaskTypePushback], 42, 12, constants.BookingConfirmed),
	}
	if err := db.Create(&muBookings).Error; err != nil {
		return err
	}

	// CZ303 at a remote stand: one signed cleaning task overlapping the
	// CA101 delayed CLEANING window (10:35-11:05) so reschedule squeezes
	// CA101 CK-02 out into PENDING at +30.
	czBase := base.Add(35 * time.Minute)
	czAt := func(min int) time.Time { return czBase.Add(time.Duration(min) * time.Minute) }
	cz303 := models.FlightTurnaround{
		FlightNo: "CZ303", AircraftReg: "B-5520", StandNo: "118",
		ArrivalTime: czAt(0), DepartureTime: czAt(60),
		TurnaroundStatus: constants.StatusInService,
	}
	if err := db.Create(&cz303).Error; err != nil {
		return err
	}
	czSigned := czAt(0)
	czTasks := []models.GroundTask{
		{TurnaroundID: cz303.ID, TaskType: constants.TaskTypeCleaning, TeamID: "TEAM-CLEAN",
			PlannedStart: czAt(0), Deadline: czAt(30), Status: constants.TaskAccepted, SignedAt: &czSigned},
	}
	if err := db.Create(&czTasks).Error; err != nil {
		return err
	}
	czBookings := []models.ResourceBooking{
		{ResourceID: resID["CK-02"], TurnaroundID: cz303.ID, TaskID: czTasks[0].ID,
			StartTime: czAt(0), EndTime: czAt(30), BookingStatus: constants.BookingConfirmed},
		// FT-03 is held all morning for CZ303, so CA101 fuel can never fall
		// back to it (itemized overlap in rejected batches).
		{ResourceID: resID["FT-03"], TurnaroundID: cz303.ID, TaskID: czTasks[0].ID,
			StartTime: at(10), EndTime: at(80), BookingStatus: constants.BookingConfirmed},
	}
	if err := db.Create(&czBookings).Error; err != nil {
		return err
	}

	// One prior delay for MU202 so the delays page has history.
	muDelay := models.DelayEvent{
		TurnaroundID: mu202.ID, DelayType: constants.DelayLateArrival, Minutes: 5,
		RootCause: "前站流量控制晚到 5 分钟", ResponsibilityTeam: "AIRLINE",
	}
	if err := db.Create(&muDelay).Error; err != nil {
		return err
	}

	logs := []models.AuditLog{
		{Actor: "dispatcher", Role: constants.RoleDispatcher,
			Action: "航班过站 CA101 已登记，机位 203，状态 已靠桥", TargetType: "FlightTurnaround", TargetID: itoa(ca101.ID)},
		{Actor: "team-bag", Role: constants.RoleTeam,
			Action: "任务 7（行李装卸）由班组 TEAM-BAG 签收，计划时点保持 09-22 09:57",
			TargetType: "GroundTask", TargetID: itoa(taskID[constants.TaskTypeBaggage])},
	}
	return db.Create(&logs).Error
}

func itoa(id uint) string {
	if id == 0 {
		return "0"
	}
	out := ""
	for id > 0 {
		out = string(rune('0'+id%10)) + out
		id /= 10
	}
	return out
}

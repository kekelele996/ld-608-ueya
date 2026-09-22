package database

import (
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"

	"gorm.io/gorm"
)

// Seed populates local demo data: resources for every task type plus a few
// turnarounds (one ready to generate a plan, one deliberately conflicting).
func Seed(db *gorm.DB) error {
	var turnCount int64
	if err := db.Model(&models.FlightTurnaround{}).Count(&turnCount).Error; err != nil {
		return err
	}
	if turnCount > 0 {
		return nil
	}

	base := time.Date(2026, 9, 21, 8, 0, 0, 0, time.Local)

	// 资源台账：每个任务类型至少两台设备，便于演示挤出/改期。
	resourceDefs := []struct {
		Code, Type, Location, Team string
		Status                     constants.ResourceStatus
		Maintenance                *time.Time
	}{
		{"BELT-01", string(constants.TaskBaggage), "T1-远机位", "TEAM-BAG", constants.ResourceAvailable, nil},
		{"BELT-02", string(constants.TaskBaggage), "T2-近机位", "TEAM-BAG", constants.ResourceAvailable, nil},
		{"CART-01", string(constants.TaskCleaning), "T1-远机位", "TEAM-CLEAN", constants.ResourceAvailable, nil},
		{"CART-02", string(constants.TaskCleaning), "T2-近机位", "TEAM-CLEAN", constants.ResourceMaintenance, &base},
		{"CATER-01", string(constants.TaskCatering), "配餐间A", "TEAM-CATER", constants.ResourceAvailable, nil},
		{"CATER-02", string(constants.TaskCatering), "配餐间B", "TEAM-CATER", constants.ResourceOffline, nil},
		{"WATER-01", string(constants.TaskWaterService), "清水站", "TEAM-WATER", constants.ResourceAvailable, nil},
		{"WATER-02", string(constants.TaskWaterService), "污水站", "TEAM-WATER", constants.ResourceAvailable, nil},
		{"FUEL-01", string(constants.TaskRefuel), "加油栓区", "TEAM-FUEL", constants.ResourceAvailable, nil},
		{"FUEL-02", string(constants.TaskRefuel), "加油栓区", "TEAM-FUEL", constants.ResourceAvailable, nil},
		{"PUSH-01", string(constants.TaskPushback), "牵引车坪", "TEAM-RAMP", constants.ResourceAvailable, nil},
		{"PUSH-02", string(constants.TaskPushback), "牵引车坪", "TEAM-RAMP", constants.ResourceAvailable, nil},
	}
	for _, def := range resourceDefs {
		res := &models.GroundResource{
			ResourceCode:       def.Code,
			ResourceType:       def.Type,
			Location:           def.Location,
			AvailabilityStatus: string(def.Status),
			MaintenanceDueAt:   def.Maintenance,
			OwnerTeam:          def.Team,
		}
		if err := db.Create(res).Error; err != nil {
			return err
		}
	}

	// 航班 A：正常航班，到站后可一键生成全部任务与预约。
	flightA := &models.FlightTurnaround{
		FlightNo:         "CA1831",
		AircraftReg:      "B-5821",
		StandNo:          "203",
		ArrivalTime:      base,
		DepartureTime:    base.Add(75 * time.Minute),
		TurnaroundStatus: string(constants.StatusArriving),
	}
	// 航班 B：已到站，演示生成与后续延误重排。
	flightB := &models.FlightTurnaround{
		FlightNo:         "MU5102",
		AircraftReg:      "B-6690",
		StandNo:          "118",
		ArrivalTime:      base.Add(30 * time.Minute),
		DepartureTime:    base.Add(30*time.Minute + 75*time.Minute),
		TurnaroundStatus: string(constants.StatusArriving),
	}
	// 航班 C：离港时间过紧，推送任务必然超截止，演示整批拒绝。
	flightC := &models.FlightTurnaround{
		FlightNo:         "CZ3108",
		AircraftReg:      "B-9941",
		StandNo:          "305",
		ArrivalTime:      base.Add(2 * time.Hour),
		DepartureTime:    base.Add(2*time.Hour + 40*time.Minute),
		TurnaroundStatus: string(constants.StatusArriving),
	}
	for _, f := range []*models.FlightTurnaround{flightA, flightB, flightC} {
		if err := db.Create(f).Error; err != nil {
			return err
		}
	}

	// 预置一条占用 BELT-01 的历史航班预约，用于演示同资源时段重叠。
	other := &models.FlightTurnaround{
		FlightNo:         "HU7801",
		AircraftReg:      "B-3217",
		StandNo:          "201",
		ArrivalTime:      base.Add(-2 * time.Hour),
		DepartureTime:    base.Add(-2*time.Hour + 70*time.Minute),
		TurnaroundStatus: string(constants.StatusInService),
		PlanGenerated:    true,
	}
	if err := db.Create(other).Error; err != nil {
		return err
	}
	task := &models.GroundTask{
		TurnaroundID: other.ID,
		TaskType:     string(constants.TaskBaggage),
		TeamID:       "TEAM-BAG",
		PlannedStart: base.Add(10 * time.Minute),
		PlannedEnd:   base.Add(40 * time.Minute),
		Deadline:     base.Add(45 * time.Minute),
		Status:       string(constants.TaskSigned),
		SignedBy:     "行李班组长·赵岩",
	}
	if err := db.Create(task).Error; err != nil {
		return err
	}
	booking := &models.ResourceBooking{
		ResourceID:    1, // BELT-01
		TurnaroundID:  other.ID,
		TaskID:        &task.ID,
		StartTime:     base.Add(10 * time.Minute),
		EndTime:       base.Add(40 * time.Minute),
		BookingStatus: string(constants.BookingConfirmed),
	}
	if err := db.Create(booking).Error; err != nil {
		return err
	}

	// 另一条历史航班占用 CATER-01 的 08:50~09:20：CA1831 生成计划时配餐窗口
	// 08:10~08:40 尚空闲（选用 CATER-01），登记 40 分钟延误后窗口平移到
	// 08:50~09:20 正好撞车，预约被挤出进入待处理，供页面演示改期确认。
	later := &models.FlightTurnaround{
		FlightNo:         "9C8903",
		AircraftReg:      "B-7032",
		StandNo:          "202",
		ArrivalTime:      base.Add(-90 * time.Minute),
		DepartureTime:    base.Add(-20 * time.Minute),
		TurnaroundStatus: string(constants.StatusInService),
		PlanGenerated:    true,
	}
	if err := db.Create(later).Error; err != nil {
		return err
	}
	laterTask := &models.GroundTask{
		TurnaroundID: later.ID,
		TaskType:     string(constants.TaskCatering),
		TeamID:       "TEAM-CATER",
		PlannedStart: base.Add(50 * time.Minute),
		PlannedEnd:   base.Add(80 * time.Minute),
		Deadline:     base.Add(85 * time.Minute),
		Status:       string(constants.TaskSigned),
		SignedBy:     "配餐班组长·许静",
	}
	if err := db.Create(laterTask).Error; err != nil {
		return err
	}
	laterBooking := &models.ResourceBooking{
		ResourceID:    5, // CATER-01
		TurnaroundID:  later.ID,
		TaskID:        &laterTask.ID,
		StartTime:     base.Add(50 * time.Minute),
		EndTime:       base.Add(80 * time.Minute),
		BookingStatus: string(constants.BookingConfirmed),
	}
	if err := db.Create(laterBooking).Error; err != nil {
		return err
	}

	fmt.Println("seed data inserted")
	return nil
}

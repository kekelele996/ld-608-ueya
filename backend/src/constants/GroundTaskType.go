package constants

// GroundTaskType enumerates every turnaround support task type.
// Appears in: types/GroundTaskType, constructors, logTemplates, errorMessages,
// filters and display components on both frontend and backend.
const (
	TaskTypeCleaning     = "CLEANING"
	TaskTypeCatering     = "CATERING"
	TaskTypeBaggage      = "BAGGAGE"
	TaskTypeRefuel       = "REFUEL"
	TaskTypeWaterService = "WATER_SERVICE"
	TaskTypePushback     = "PUSHBACK"
)

// GroundTaskType is the ordered shared enum list.
var GroundTaskType = []string{
	TaskTypeCleaning,
	TaskTypeCatering,
	TaskTypeBaggage,
	TaskTypeRefuel,
	TaskTypeWaterService,
	TaskTypePushback,
}

// GroundTaskTypeText drives status text formatters and filters.
var GroundTaskTypeText = map[string]string{
	TaskTypeCleaning:     "客舱清洁",
	TaskTypeCatering:     "航食配送",
	TaskTypeBaggage:      "行李装卸",
	TaskTypeRefuel:       "航油加注",
	TaskTypeWaterService: "清水污水",
	TaskTypePushback:     "飞机推出",
}

// TaskDefaultTeam maps each task type to its responsible team code.
var TaskDefaultTeam = map[string]string{
	TaskTypeCleaning:     "TEAM-CLEAN",
	TaskTypeCatering:     "TEAM-CATER",
	TaskTypeBaggage:      "TEAM-BAG",
	TaskTypeRefuel:       "TEAM-FUEL",
	TaskTypeWaterService: "TEAM-WATER",
	TaskTypePushback:     "TEAM-RAMP",
}

// TaskOffsetMinutes describes (offset after arrival in minutes, duration in
// minutes) per task type. The plan builder stamps tasks from arrival time.
var TaskOffsetMinutes = map[string][2]int{
	TaskTypeBaggage:      {2, 25},
	TaskTypeCleaning:     {5, 30},
	TaskTypeWaterService: {8, 20},
	TaskTypeCatering:     {10, 25},
	TaskTypeRefuel:       {12, 25},
	TaskTypePushback:     {42, 12},
}

// TaskResourceType maps a task type to the resource type it consumes.
var TaskResourceType = map[string]string{
	TaskTypeCleaning:     "CLEANING_KIT",
	TaskTypeCatering:     "CATERING_TRUCK",
	TaskTypeBaggage:      "BELT_LOADER",
	TaskTypeRefuel:       "FUEL_TRUCK",
	TaskTypeWaterService: "WATER_TRUCK",
	TaskTypePushback:     "PUSHBACK_TRACTOR",
}

func IsValidTaskType(value string) bool {
	for _, item := range GroundTaskType {
		if item == value {
			return true
		}
	}
	return false
}

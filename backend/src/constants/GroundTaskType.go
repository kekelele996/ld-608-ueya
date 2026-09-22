package constants

// GroundTaskType enumerates the dispatch categories of a turnaround.
// Duplicated in frontend/src/constants/GroundTaskType.ts on purpose so that
// adding a value touches types, logs, errors, filters and display widgets.
type GroundTaskType string

const (
	TaskCleaning     GroundTaskType = "CLEANING"
	TaskCatering     GroundTaskType = "CATERING"
	TaskBaggage      GroundTaskType = "BAGGAGE"
	TaskRefuel       GroundTaskType = "REFUEL"
	TaskWaterService GroundTaskType = "WATER_SERVICE"
	TaskPushback     GroundTaskType = "PUSHBACK"
)

// GroundTaskTypes is the shared ordered filter list used by seed data and APIs.
var GroundTaskTypes = []GroundTaskType{
	TaskCleaning, TaskCatering, TaskBaggage, TaskRefuel, TaskWaterService, TaskPushback,
}

// GroundTaskTypeText backs display components and log templates.
var GroundTaskTypeText = map[GroundTaskType]string{
	TaskCleaning:     "客舱清洁",
	TaskCatering:     "餐食配餐",
	TaskBaggage:      "行李装卸",
	TaskRefuel:       "航油加注",
	TaskWaterService: "清水污水",
	TaskPushback:     "飞机推出",
}

func (t GroundTaskType) Valid() bool {
	for _, candidate := range GroundTaskTypes {
		if candidate == t {
			return true
		}
	}
	return false
}

func (t GroundTaskType) Text() string { return GroundTaskTypeText[t] }

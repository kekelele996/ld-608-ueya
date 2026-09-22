package constants

// TeamInfo binds a team id to its display label and the task types it serves.
// Task generation chooses the owning team per GroundTaskType through this table.
type TeamInfo struct {
	ID        string
	Name      string
	TaskKinds []GroundTaskType
}

var Teams = []TeamInfo{
	{"TEAM-CLEAN", "客舱清洁班", []GroundTaskType{TaskCleaning}},
	{"TEAM-CATER", "航机配餐班", []GroundTaskType{TaskCatering}},
	{"TEAM-BAG", "行李装卸班", []GroundTaskType{TaskBaggage}},
	{"TEAM-FUEL", "航油保障班", []GroundTaskType{TaskRefuel}},
	{"TEAM-WATER", "清水污水班", []GroundTaskType{TaskWaterService}},
	{"TEAM-RAMP", "机坪牵引班", []GroundTaskType{TaskPushback}},
}

// TeamForTask returns the owning team for a task type.
func TeamForTask(kind GroundTaskType) TeamInfo {
	for _, team := range Teams {
		for _, served := range team.TaskKinds {
			if served == kind {
				return team
			}
		}
	}
	return Teams[0]
}

func TeamText(teamID string) string {
	for _, team := range Teams {
		if team.ID == teamID {
			return team.Name
		}
	}
	return teamID
}

// ResourceTypeForTask maps a task type to the resource pool it consumes.
func ResourceTypeForTask(kind GroundTaskType) string {
	return string(kind)
}

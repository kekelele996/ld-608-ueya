package constants

// Machine-readable conflict reason codes returned per rejected item so the
// UI ConflictBadge / conflict panel can group and explain each one.
const (
	ConflictOverlap     = "RESOURCE_TIME_OVERLAP"
	ConflictMaintenance = "RESOURCE_MAINTENANCE"
	ConflictOffline     = "RESOURCE_OFFLINE"
	ConflictNoResource  = "NO_MATCHING_RESOURCE"
	ConflictDeadline    = "TASK_DEADLINE_AFTER_DEPARTURE"
	ConflictBeforeArrival = "TASK_BEFORE_ARRIVAL"
	ConflictPlanExists  = "PLAN_ALREADY_EXISTS"
)

var ConflictReasonText = map[string]string{
	ConflictOverlap:       "同一资源时段重叠",
	ConflictMaintenance:   "资源处于维护期",
	ConflictOffline:       "资源已离线",
	ConflictNoResource:    "无匹配可用资源",
	ConflictDeadline:      "任务截止晚于离港",
	ConflictBeforeArrival: "任务早于到站时间",
	ConflictPlanExists:    "保障计划已存在",
}

// Conflict reason codes mirror backend constants/conflictReasons.go.
// ConflictBadge and the conflict panel group and explain each one.
export const CONFLICT_REASON_TEXT: Record<string, string> = {
  RESOURCE_TIME_OVERLAP: "同一资源时段重叠",
  RESOURCE_MAINTENANCE: "资源处于维护期",
  RESOURCE_OFFLINE: "资源已离线",
  NO_MATCHING_RESOURCE: "无匹配可用资源",
  TASK_DEADLINE_AFTER_DEPARTURE: "任务截止晚于离港",
  TASK_BEFORE_ARRIVAL: "任务早于到站时间",
  PLAN_ALREADY_EXISTS: "保障计划已存在",
};

export const CONFLICT_REASON_COLOR: Record<string, string> = {
  RESOURCE_TIME_OVERLAP: "error",
  RESOURCE_MAINTENANCE: "warning",
  RESOURCE_OFFLINE: "default",
  NO_MATCHING_RESOURCE: "error",
  TASK_DEADLINE_AFTER_DEPARTURE: "error",
  TASK_BEFORE_ARRIVAL: "warning",
  PLAN_ALREADY_EXISTS: "default",
};

export const GROUND_TASK_STATUS_COLOR: Record<string, string> = {
  PLANNED: "default",
  ACCEPTED: "processing",
  COMPLETED: "success",
  BLOCKED: "error",
};

export const BOOKING_STATUS_COLOR: Record<string, string> = {
  CONFIRMED: "success",
  PENDING: "warning",
  RELEASED: "default",
};

package constants

// Error codes are returned in the JSON error envelope and mirrored in
// frontend/src/constants/errorCodes.ts.
const (
	CodeAuthRequired      = "AUTH_REQUIRED"
	CodeAuthInvalid       = "AUTH_INVALID"
	CodeRBACDenied        = "RBAC_DENIED"
	CodeValidationFailed  = "VALIDATION_FAILED"
	CodeRateLimited       = "RATE_LIMITED"
	CodeNotFound          = "NOT_FOUND"
	CodePlanConflict      = "PLAN_CONFLICT"
	CodeBookingConflict   = "BOOKING_CONFLICT"
	CodeInvalidTransition = "INVALID_TRANSITION"
	CodeInternal          = "INTERNAL"
)

// Item-level conflict codes explain why a single generated task/booking was
// rejected and are surfaced through ConflictBadge / 冲突面板.
const (
	ItemResourceOverlap     = "RESOURCE_TIME_OVERLAP"
	ItemResourceMaintenance = "RESOURCE_MAINTENANCE"
	ItemResourceOffline     = "RESOURCE_OFFLINE"
	ItemResourceWindow      = "RESOURCE_WINDOW_UNAVAILABLE"
	ItemDeadlineConflict    = "TASK_DEADLINE_CONFLICT"
	ItemNoResource          = "NO_MATCHING_RESOURCE"
	ItemBookingBumped       = "BOOKING_BUMPED"
)

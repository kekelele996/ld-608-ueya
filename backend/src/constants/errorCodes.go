package constants

// Error codes are referenced by services, controllers, frontend errorMessages
// and the global error handler. Keep codes stable across layers.
const (
	CodeAuthRequired       = "AUTH_REQUIRED"
	CodeRBACDenied         = "RBAC_DENIED"
	CodeValidationFailed   = "VALIDATION_FAILED"
	CodeRateLimited        = "RATE_LIMITED"
	CodeNotFound           = "NOT_FOUND"
	CodeConflict           = "PLAN_CONFLICT"
	CodeBookingConflict    = "BOOKING_CONFLICT"
	CodeResourceUnusable   = "RESOURCE_UNUSABLE"
	CodeDeadlineConflict   = "DEADLINE_CONFLICT"
	CodeAlreadyExists      = "ALREADY_EXISTS"
	CodeInvalidTransition  = "INVALID_TRANSITION"
	CodeInternal           = "INTERNAL_ERROR"
)

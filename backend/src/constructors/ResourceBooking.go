package constructors

import (
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
)

// NewResourceBooking constructs a CONFIRMED appointment for a generated task.
func NewResourceBooking(resourceID, turnaroundID uint, taskID uint, start, end time.Time) *models.ResourceBooking {
	return &models.ResourceBooking{
		ResourceID:    resourceID,
		TurnaroundID:  turnaroundID,
		TaskID:        &taskID,
		StartTime:     start,
		EndTime:       end,
		BookingStatus: string(constants.BookingConfirmed),
	}
}

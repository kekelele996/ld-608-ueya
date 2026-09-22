package services

import (
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/utils"
)

// overlaps reports half-open interval overlap: [start,end) x [otherStart,otherEnd).
func overlaps(start, end, otherStart, otherEnd time.Time) bool {
	return start.Before(otherEnd) && otherStart.Before(end)
}

// resourceCovers verifies status, maintenance, offline and availability window.
// Returns a conflict code + message when the resource cannot host the window.
func resourceCovers(resource *models.GroundResource, start, end time.Time) (string, string) {
	switch constants.ResourceStatus(resource.AvailabilityStatus) {
	case constants.ResourceOffline:
		return constants.ItemResourceOffline,
			fmt.Sprintf(constants.ConflictMessages[constants.ItemResourceOffline], resource.ResourceCode)
	case constants.ResourceMaintenance:
		return constants.ItemResourceMaintenance,
			fmt.Sprintf(constants.ConflictMessages[constants.ItemResourceMaintenance],
				resource.ResourceCode, utils.FormatTimePtr(resource.MaintenanceDueAt))
	}

	if resource.AvailableFrom != nil || resource.AvailableTo != nil {
		from, to := time.Time{}, time.Time{}
		if resource.AvailableFrom != nil {
			from = *resource.AvailableFrom
		}
		if resource.AvailableTo != nil {
			to = *resource.AvailableTo
		}
		if (!from.IsZero() && start.Before(from)) || (!to.IsZero() && end.After(to)) {
			return constants.ItemResourceWindow,
				fmt.Sprintf(constants.ConflictMessages[constants.ItemResourceWindow],
					resource.ResourceCode, utils.FormatTime(from), utils.FormatTime(to),
					utils.FormatTime(start), utils.FormatTime(end))
		}
	}
	return "", ""
}

// bookingOverlap finds a confirmed booking colliding with [start,end).
// excludeID lets adjustment ignore itself; ignoreTurnaround ignores the same
// turnaround during in-batch candidate scoring.
func bookingOverlap(bookings []models.ResourceBooking, resourceID uint, start, end time.Time, excludeID, ignoreTurnaround uint) *models.ResourceBooking {
	for i := range bookings {
		b := &bookings[i]
		if b.ResourceID != resourceID || b.ID == excludeID {
			continue
		}
		if b.BookingStatus != string(constants.BookingConfirmed) {
			continue
		}
		if ignoreTurnaround > 0 && b.TurnaroundID == ignoreTurnaround {
			continue
		}
		if overlaps(start, end, b.StartTime, b.EndTime) {
			return b
		}
	}
	return nil
}

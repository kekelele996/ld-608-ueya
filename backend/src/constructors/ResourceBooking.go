package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewBookingViews builds booking DTOs with joined resource display fields.
func NewBookingViews(rows []models.ResourceBooking, resourceMap map[uint]models.GroundResource) []types.BookingView {
	out := make([]types.BookingView, 0, len(rows))
	for _, row := range rows {
		resource, ok := resourceMap[row.ResourceID]
		var ptr *models.GroundResource
		if ok {
			ptr = &resource
		}
		out = append(out, NewBookingView(row, ptr))
	}
	return out
}

// NewPendingBooking moves a confirmed booking into the pending queue and
// stamps the machine-readable reason shown by ConflictBadge.
func NewPendingBooking(booking *models.ResourceBooking, reasonCode string) {
	booking.BookingStatus = constants.BookingPending
	booking.ConflictReason = reasonCode
}

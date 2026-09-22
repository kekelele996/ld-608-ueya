package services

import (
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func taskViews(rows []models.GroundTask) []types.TaskView {
	return constructors.NewTaskViews(rows)
}

func bookingViews(rows []models.ResourceBooking, resourceMap map[uint]models.GroundResource) []types.BookingView {
	return constructors.NewBookingViews(rows, resourceMap)
}

package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// TurnaroundView is the detail response carrying nested tasks/bookings.
type TurnaroundView struct {
	ID                  uint                 `json:"id"`
	FlightNo            string               `json:"flight_no"`
	AircraftReg         string               `json:"aircraft_reg"`
	StandNo             string               `json:"stand_no"`
	ArrivalTime         string               `json:"arrival_time"`
	DepartureTime       string               `json:"departure_time"`
	TurnaroundStatus    string               `json:"turnaround_status"`
	StatusText          string               `json:"status_text"`
	DelayReason         string               `json:"delay_reason"`
	TotalDelayMinutes   int                  `json:"total_delay_minutes"`
	TaskCount           int                  `json:"task_count"`
	CompletedCount      int                  `json:"completed_count"`
	AcceptedCount       int                  `json:"accepted_count"`
	PendingBookings     int                  `json:"pending_bookings"`
	Tasks               []types.TaskView     `json:"tasks"`
	Bookings            []types.BookingView  `json:"bookings"`
}

// NewTurnaroundView assembles a turnaround detail with tasks and bookings.
func NewTurnaroundView(t models.FlightTurnaround, tasks []models.GroundTask, bookings []models.ResourceBooking, resourceMap map[uint]models.GroundResource) TurnaroundView {
	view := TurnaroundView{
		ID:                t.ID,
		FlightNo:          t.FlightNo,
		AircraftReg:       t.AircraftReg,
		StandNo:           t.StandNo,
		ArrivalTime:       utils.FormatTime(t.ArrivalTime),
		DepartureTime:     utils.FormatTime(t.DepartureTime),
		TurnaroundStatus:  t.TurnaroundStatus,
		StatusText:        constants.TurnaroundStatusText[t.TurnaroundStatus],
		DelayReason:       t.DelayReason,
		TotalDelayMinutes: t.TotalDelayMinutes,
		Tasks:             NewTaskViews(tasks),
		Bookings:          make([]types.BookingView, 0, len(bookings)),
	}
	for _, task := range tasks {
		view.TaskCount++
		if task.Status == constants.TaskCompleted {
			view.CompletedCount++
		}
		if task.Status == constants.TaskAccepted {
			view.AcceptedCount++
		}
	}
	for _, booking := range bookings {
		resource, ok := resourceMap[booking.ResourceID]
		var ptr *models.GroundResource
		if ok {
			ptr = &resource
		}
		bv := NewBookingView(booking, ptr)
		view.Bookings = append(view.Bookings, bv)
		if booking.BookingStatus == constants.BookingPending {
			view.PendingBookings++
		}
	}
	return view
}

// NewTurnaroundSummary builds list-row payloads without nested entities.
func NewTurnaroundSummary(t models.FlightTurnaround) map[string]interface{} {
	return map[string]interface{}{
		"id":                  t.ID,
		"flight_no":           t.FlightNo,
		"aircraft_reg":        t.AircraftReg,
		"stand_no":            t.StandNo,
		"arrival_time":        utils.FormatTime(t.ArrivalTime),
		"departure_time":      utils.FormatTime(t.DepartureTime),
		"turnaround_status":   t.TurnaroundStatus,
		"status_text":         constants.TurnaroundStatusText[t.TurnaroundStatus],
		"delay_reason":        t.DelayReason,
		"total_delay_minutes": t.TotalDelayMinutes,
	}
}

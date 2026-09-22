package types

import (
	"time"

	"groundTurn/src/models"
)

// TurnaroundDetail aggregates the closed-loop view served to the detail page.
type TurnaroundDetail struct {
	models.FlightTurnaround
	Tasks    []models.GroundTask      `json:"tasks"`
	Bookings []models.ResourceBooking `json:"bookings"`
	Delays   []models.DelayEvent      `json:"delays"`
}

// ReplanResult is returned after delay registration, describing what moved.
type ReplanResult struct {
	DelayEventID  uint                     `json:"delay_event_id"`
	ShiftMinutes  int                      `json:"shift_minutes"`
	KeptTasks     []KeptTaskView           `json:"kept_tasks"`
	MovedTasks    []models.GroundTask      `json:"moved_tasks"`
	MovedBookings []models.ResourceBooking `json:"moved_bookings"`
	Bumped        []ConflictItem           `json:"bumped"`
}

// KeptTaskView records an acknowledged task whose time points were frozen.
type KeptTaskView struct {
	TaskID       uint      `json:"task_id"`
	TaskType     string    `json:"task_type"`
	Status       string    `json:"status"`
	PlannedStart time.Time `json:"planned_start"`
	SignedBy     string    `json:"signed_by"`
	Reason       string    `json:"reason"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

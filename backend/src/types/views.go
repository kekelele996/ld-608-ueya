// Package types carries request/response envelope declarations shared by
// controllers and middleware. Entity DTO constructors live in /constructors.
package types

// ConflictItem describes a single per-item failure in a rejected batch.
type ConflictItem struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	TaskType    string `json:"task_type,omitempty"`
	ResourceID  uint   `json:"resource_id,omitempty"`
	ResourceCode string `json:"resource_code,omitempty"`
	BookingID   uint   `json:"booking_id,omitempty"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
}

// PlanResult is returned by plan generation / delay registration.
// Saved=false means the whole batch was rejected and nothing persisted.
type PlanResult struct {
	Saved        bool           `json:"saved"`
	TurnaroundID uint          `json:"turnaround_id"`
	Tasks        []TaskView    `json:"tasks"`
	Bookings     []BookingView `json:"bookings"`
	Conflicts    []ConflictItem `json:"conflicts"`
}

// RescheduleResult mirrors PlanResult with reschedule counters.
type RescheduleResult struct {
	PlanResult
	RescheduledTasks  int `json:"rescheduled_tasks"`
	FrozenTasks       int `json:"frozen_tasks"`
	PendingBookings   int `json:"pending_bookings"`
	AddedDelayMinutes int `json:"added_delay_minutes"`
}

type TaskView struct {
	ID           uint   `json:"id"`
	TurnaroundID uint   `json:"turnaround_id"`
	TaskType     string `json:"task_type"`
	TaskTypeText string `json:"task_type_text"`
	TeamID       string `json:"team_id"`
	PlannedStart string `json:"planned_start"`
	Deadline     string `json:"deadline"`
	ActualFinish string `json:"actual_finish"`
	Status       string `json:"status"`
	StatusText   string `json:"status_text"`
	BlockerNote  string `json:"blocker_note"`
	SignedAt     string `json:"signed_at"`
}

type BookingView struct {
	ID             uint   `json:"id"`
	ResourceID     uint   `json:"resource_id"`
	ResourceCode   string `json:"resource_code"`
	ResourceType   string `json:"resource_type"`
	TurnaroundID   uint   `json:"turnaround_id"`
	TaskID         uint   `json:"task_id"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	BookingStatus  string `json:"booking_status"`
	StatusText     string `json:"status_text"`
	ConflictCode   string `json:"conflict_code"`
	ConflictReason string `json:"conflict_reason"`
}

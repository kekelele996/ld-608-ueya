package types

import "fmt"

// AppError is wrapped independently by services and controllers (never only
// swallowed by a single global handler), carrying an error code and HTTP hint.
type AppError struct {
	Code       string
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ConflictItem explains a single rejected task / booking during batch
// generation or a pending booking after delay replanning.
type ConflictItem struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	TaskType      string `json:"task_type,omitempty"`
	TaskID        *uint  `json:"task_id,omitempty"`
	BookingID     *uint  `json:"booking_id,omitempty"`
	ResourceID    *uint  `json:"resource_id,omitempty"`
	ResourceCode  string `json:"resource_code,omitempty"`
	TeamID        string `json:"team_id,omitempty"`
	PlannedStart  string `json:"planned_start,omitempty"`
	PlannedEnd    string `json:"planned_end,omitempty"`
	Deadline      string `json:"deadline,omitempty"`
}

// ConflictResult is the all-or-nothing generation response payload.
type ConflictResult struct {
	Saved      bool           `json:"saved"`
	Items      []ConflictItem `json:"items"`
	TaskCount  int            `json:"task_count"`
}

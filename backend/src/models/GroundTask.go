package models

import "time"

// GroundTask belongs to one turnaround and consumes resources via bookings.
type GroundTask struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	TurnaroundID uint       `gorm:"not null;index" json:"turnaround_id"`
	TaskType     string     `gorm:"size:32;not null;index" json:"task_type"`
	TeamID       string     `gorm:"size:32;not null" json:"team_id"`
	PlannedStart time.Time  `gorm:"not null;index" json:"planned_start"`
	Deadline     time.Time  `gorm:"not null" json:"deadline"`
	ActualFinish *time.Time `json:"actual_finish"`
	Status       string     `gorm:"size:16;not null;index" json:"status"`
	BlockerNote  string     `gorm:"size:255" json:"blocker_note"`
	// SignedAt freezes the task against later delay-driven reschedules.
	SignedAt  *time.Time `json:"signed_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (GroundTask) TableName() string { return "ground_task" }

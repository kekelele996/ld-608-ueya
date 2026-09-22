package models

import "time"

// GroundTask 地勤任务：按任务类型派发给班组的单项保障作业。
type GroundTask struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	TurnaroundID uint       `gorm:"not null;index" json:"turnaround_id"`
	TaskType     string     `gorm:"size:32;not null;index" json:"task_type"`
	TeamID       string     `gorm:"size:32;not null;index" json:"team_id"`
	PlannedStart time.Time  `gorm:"not null" json:"planned_start"`
	PlannedEnd   time.Time  `gorm:"not null" json:"planned_end"`
	Deadline     time.Time  `gorm:"not null" json:"deadline"`
	ActualFinish *time.Time `json:"actual_finish"`
	Status       string     `gorm:"size:32;not null;default:PLANNED" json:"status"`
	BlockerNote  string     `gorm:"size:255" json:"blocker_note"`
	SignedBy     string     `gorm:"size:64" json:"signed_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (GroundTask) TableName() string { return "ground_task" }

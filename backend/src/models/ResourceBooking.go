package models

import "time"

// ResourceBooking links flight, task and resource across one time window.
type ResourceBooking struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ResourceID     uint      `gorm:"not null;index" json:"resource_id"`
	TurnaroundID   uint      `gorm:"not null;index" json:"turnaround_id"`
	TaskID         uint      `gorm:"not null;index" json:"task_id"`
	StartTime      time.Time `gorm:"not null;index" json:"start_time"`
	EndTime        time.Time `gorm:"not null" json:"end_time"`
	BookingStatus  string    `gorm:"size:16;not null;index" json:"booking_status"`
	ConflictReason string    `gorm:"size:64" json:"conflict_reason"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ResourceBooking) TableName() string { return "resource_booking" }

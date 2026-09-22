package models

import "time"

// ResourceBooking 资源预约：连接航班、任务和资源的时段占用记录。
type ResourceBooking struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ResourceID     uint      `gorm:"not null;index" json:"resource_id"`
	TurnaroundID   uint      `gorm:"not null;index" json:"turnaround_id"`
	TaskID         *uint     `gorm:"index" json:"task_id"`
	StartTime      time.Time `gorm:"not null;index" json:"start_time"`
	EndTime        time.Time `gorm:"not null" json:"end_time"`
	BookingStatus  string    `gorm:"size:32;not null;default:CONFIRMED" json:"booking_status"`
	ConflictReason string    `gorm:"size:255" json:"conflict_reason"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ResourceBooking) TableName() string { return "resource_booking" }

// TaskIDValue renders the optional owning task id for conflict messages.
func (b ResourceBooking) TaskIDValue() uint {
	if b.TaskID == nil {
		return 0
	}
	return *b.TaskID
}

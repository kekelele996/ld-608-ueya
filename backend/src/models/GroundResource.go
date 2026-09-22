package models

import "time"

// GroundResource 保障资源：可被预约的车辆/设备/工位。
type GroundResource struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	ResourceCode       string     `gorm:"size:32;not null;uniqueIndex" json:"resource_code"`
	ResourceType       string     `gorm:"size:32;not null;index" json:"resource_type"`
	Location           string     `gorm:"size:64;not null" json:"location"`
	AvailabilityStatus string     `gorm:"size:32;not null;default:AVAILABLE" json:"availability_status"`
	MaintenanceDueAt   *time.Time `json:"maintenance_due_at"`
	OwnerTeam          string     `gorm:"size:32;not null" json:"owner_team"`
	// 班组/资源管理员维护的可用时段窗口（为空表示全天可用）。
	AvailableFrom *time.Time `json:"available_from"`
	AvailableTo   *time.Time `json:"available_to"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (GroundResource) TableName() string { return "ground_resource" }

package models

import "time"

// GroundResource is the station equipment ledger.
type GroundResource struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	ResourceCode       string     `gorm:"size:32;not null;uniqueIndex" json:"resource_code"`
	ResourceType       string     `gorm:"size:32;not null;index" json:"resource_type"`
	Location           string     `gorm:"size:64" json:"location"`
	AvailabilityStatus string     `gorm:"size:16;not null;index" json:"availability_status"`
	MaintenanceDueAt   *time.Time `json:"maintenance_due_at"`
	// MaintenanceEnd is exclusive: unavailable during
	// [maintenance_due_at, maintenance_end).
	MaintenanceEnd *time.Time `json:"maintenance_end"`
	OwnerTeam      string     `gorm:"size:32;not null" json:"owner_team"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (GroundResource) TableName() string { return "ground_resource" }

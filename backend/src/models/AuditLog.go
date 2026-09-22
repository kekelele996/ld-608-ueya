package models

import "time"

// AuditLog backs the operation-log cross-cutting concern.
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Actor      string    `gorm:"size:64;not null;index" json:"actor"`
	Role       string    `gorm:"size:32;not null" json:"role"`
	Action     string    `gorm:"size:255;not null" json:"action"`
	TargetType string    `gorm:"size:32;not null;index" json:"target_type"`
	TargetID   string    `gorm:"size:64;index" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_log" }

// User backs JWT + RBAC. Seeded locally, never sourced from a third party.
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	DisplayName string    `gorm:"size:64;not null" json:"display_name"`
	Role        string    `gorm:"size:32;not null" json:"role"`
	TeamID      string    `gorm:"size:32" json:"team_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (User) TableName() string { return "app_user" }

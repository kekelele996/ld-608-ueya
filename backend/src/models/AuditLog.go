package models

import "time"

// AuditLog 操作日志：由 auditLogMiddleware 与 service 层共同写入。
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Actor      string    `gorm:"size:64;not null;index" json:"actor"`
	Role       string    `gorm:"size:32;not null" json:"role"`
	Action     string    `gorm:"size:64;not null;index" json:"action"`
	TargetType string    `gorm:"size:32;not null" json:"target_type"`
	TargetID   string    `gorm:"size:64" json:"target_id"`
	Detail     string    `gorm:"size:512" json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_log" }

// User 登录账号，承载 RBAC 角色与所属班组。
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:128;not null" json:"-"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Role      string    `gorm:"size:32;not null" json:"role"`
	TeamID    string    `gorm:"size:32" json:"team_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string { return "app_user" }

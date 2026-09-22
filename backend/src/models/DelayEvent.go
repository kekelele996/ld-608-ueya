package models

import "time"

// DelayEvent 延误事件：登记后触发未签收任务与资源预约的整批重排。
type DelayEvent struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	TurnaroundID       uint       `gorm:"not null;index" json:"turnaround_id"`
	DelayType          string     `gorm:"size:32;not null" json:"delay_type"`
	Minutes            int        `gorm:"not null" json:"minutes"`
	RootCause          string     `gorm:"size:255" json:"root_cause"`
	ResponsibilityTeam string     `gorm:"size:32" json:"responsibility_team"`
	ResolvedAt         *time.Time `json:"resolved_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (DelayEvent) TableName() string { return "delay_event" }

package models

import "time"

// DelayEvent records a registered delay and drives plan rescheduling.
type DelayEvent struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	TurnaroundID       uint       `gorm:"not null;index" json:"turnaround_id"`
	DelayType          string     `gorm:"size:32;not null" json:"delay_type"`
	Minutes            int        `gorm:"not null" json:"minutes"`
	RootCause          string     `gorm:"size:255" json:"root_cause"`
	ResponsibilityTeam string     `gorm:"size:32" json:"responsibility_team"`
	ResolvedAt         *time.Time `json:"resolved_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

func (DelayEvent) TableName() string { return "delay_event" }

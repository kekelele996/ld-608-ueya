package models

import "time"

// FlightTurnaround is the root entity of the station turnaround closed loop.
type FlightTurnaround struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	FlightNo          string    `gorm:"size:32;not null;index" json:"flight_no"`
	AircraftReg       string    `gorm:"size:32;not null" json:"aircraft_reg"`
	StandNo           string    `gorm:"size:16;not null" json:"stand_no"`
	ArrivalTime       time.Time `gorm:"not null;index" json:"arrival_time"`
	DepartureTime     time.Time `gorm:"not null" json:"departure_time"`
	TurnaroundStatus  string    `gorm:"size:16;not null;index" json:"turnaround_status"`
	DelayReason       string    `gorm:"size:255" json:"delay_reason"`
	TotalDelayMinutes int       `gorm:"not null;default:0" json:"total_delay_minutes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (FlightTurnaround) TableName() string { return "flight_turnaround" }

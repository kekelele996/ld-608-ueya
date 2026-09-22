package types

import "time"

// Request DTOs bound from JSON bodies (constructors build default forms).

type CreateTurnaroundRequest struct {
	FlightNo      string    `json:"flight_no" binding:"required"`
	AircraftReg   string    `json:"aircraft_reg" binding:"required"`
	StandNo       string    `json:"stand_no" binding:"required"`
	ArrivalTime   time.Time `json:"arrival_time" binding:"required"`
	DepartureTime time.Time `json:"departure_time" binding:"required"`
}

type ArriveRequest struct {
	ArrivalTime *time.Time `json:"arrival_time"` // 为空则取当前时间
}

type RegisterDelayRequest struct {
	DelayType          string `json:"delay_type" binding:"required"`
	Minutes            int    `json:"minutes" binding:"required,min=1"`
	RootCause          string `json:"root_cause"`
	ResponsibilityTeam string `json:"responsibility_team"`
}

type SignTaskRequest struct {
	Operator string `json:"operator"`
}

type FinishTaskRequest struct {
	Note string `json:"note"`
}

type BlockTaskRequest struct {
	Note string `json:"note" binding:"required"`
}

type AdjustBookingRequest struct {
	ResourceID uint      `json:"resource_id" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" binding:"required"`
}

type MaintenanceRequest struct {
	Status          string     `json:"status" binding:"required"` // MAINTENANCE / AVAILABLE / OFFLINE
	MaintenanceDueAt *time.Time `json:"maintenance_due_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

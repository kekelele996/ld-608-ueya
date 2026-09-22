package constructors

// RegisterTurnaroundRequest is the arrival registration form.
type RegisterTurnaroundRequest struct {
	FlightNo      string `json:"flight_no" binding:"required"`
	AircraftReg   string `json:"aircraft_reg" binding:"required"`
	StandNo       string `json:"stand_no" binding:"required"`
	ArrivalTime   string `json:"arrival_time" binding:"required"`
	DepartureTime string `json:"departure_time" binding:"required"`
}

// GeneratePlanRequest triggers whole-batch task + booking generation.
type GeneratePlanRequest struct {
	// TaskTypes optional subset; defaults to the full enum when empty.
	TaskTypes []string `json:"task_types"`
}

// RegisterDelayRequest is the delay registration form.
type RegisterDelayRequest struct {
	DelayType          string `json:"delay_type" binding:"required"`
	Minutes            int    `json:"minutes" binding:"required"`
	RootCause          string `json:"root_cause"`
	ResponsibilityTeam string `json:"responsibility_team"`
}

// AdjustBookingRequest moves a pending booking to a new window.
type AdjustBookingRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

// BlockTaskRequest marks a task blocked with a note.
type BlockTaskRequest struct {
	BlockerNote string `json:"blocker_note" binding:"required"`
}

// ResourceMaintenanceRequest registers a maintenance window / status.
type ResourceStatusRequest struct {
	Status           string `json:"status"`
	MaintenanceStart string `json:"maintenance_start"`
	MaintenanceEnd   string `json:"maintenance_end"`
}

// LoginRequest is the local JWT login form.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
}

// LoginResponse returns token + identity for the auth store.
type LoginResponse struct {
	Token       string `json:"token"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	TeamID      string `json:"team_id"`
}

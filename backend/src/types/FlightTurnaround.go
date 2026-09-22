package types

// FlightTurnaround-specific view aliases keep turnaround changes coupled
// across layers (routes -> controllers -> services -> repositories).

// PlanPreviewRow is one projected task/booking pair before persistence.
type PlanPreviewRow struct {
	TaskType     string
	TeamID       string
	PlannedStart string
	PlannedEnd   string
	Deadline     string
	ResourceCode string
}

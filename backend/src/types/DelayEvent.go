package types

// DelayImpact is one row of the delay attribution statistics panel.
type DelayImpact struct {
	DelayType          string `json:"delay_type"`
	EventCount         int    `json:"event_count"`
	TotalMinutes       int    `json:"total_minutes"`
	ResponsibilityTeam string `json:"responsibility_team"`
}

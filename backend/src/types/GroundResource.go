package types

// ResourceAvailabilityView joins a resource with its upcoming bookings.
type ResourceAvailabilityView struct {
	ResourceID      uint     `json:"resource_id"`
	ResourceCode    string   `json:"resource_code"`
	Status          string   `json:"status"`
	BookedWindows   []string `json:"booked_windows"`
	HasConflict     bool     `json:"has_conflict"`
}

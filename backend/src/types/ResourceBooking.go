package types

// BookingConflictView is the payload behind ConflictBadge on resources page.
type BookingConflictView struct {
	BookingID      uint   `json:"booking_id"`
	ResourceCode   string `json:"resource_code"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	BookingStatus  string `json:"booking_status"`
	ConflictReason string `json:"conflict_reason"`
}

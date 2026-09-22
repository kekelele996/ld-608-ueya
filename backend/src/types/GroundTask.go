package types

// TaskProgress is computed by useTurnaroundProgress's backend counterpart and
// consumed by dashboard/stat responses.
type TaskProgress struct {
	Total     int `json:"total"`
	Planned   int `json:"planned"`
	Signed    int `json:"signed"`
	Finished  int `json:"finished"`
	Blocked   int `json:"blocked"`
}

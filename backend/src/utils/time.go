package utils

import (
	"time"
)

// TimeLayout is the canonical on-the-wire format shared with frontend formatters.
const TimeLayout = "2006-01-02 15:04"

// FormatTime renders local time consistently for log templates and conflict
// reasons; frontend utils/formatters.ts mirrors this presentation.
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Local().Format(TimeLayout)
}

func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return FormatTime(*t)
}

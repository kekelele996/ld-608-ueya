package utils

import "time"

// APILayout is the wire format shared with frontend constructors/formatters.
const APILayout = time.RFC3339

// FormatTime renders time for JSON views.
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(APILayout)
}

// FormatTimePtr renders an optional time.
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatTime(*t)
}

// ParseTime parses incoming RFC3339 timestamps.
func ParseTime(value string) (time.Time, error) {
	return time.Parse(APILayout, value)
}

// Overlap reports whether half-open windows [aStart,aEnd) and [bStart,bEnd)
// touch. Zero-length windows never overlap.
func Overlap(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}

// AddMinutes is the single helper used by plan build and delay reschedule.
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

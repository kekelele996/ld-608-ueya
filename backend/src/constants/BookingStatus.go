package constants

// BookingStatus is the resource appointment state.
type BookingStatus string

const (
	BookingConfirmed BookingStatus = "CONFIRMED" // 已确认占用资源时段
	BookingPending   BookingStatus = "PENDING"   // 被挤出/待调度处理
	BookingReleased  BookingStatus = "RELEASED"  // 已释放
)

var BookingStatuses = []BookingStatus{
	BookingConfirmed, BookingPending, BookingReleased,
}

var BookingStatusText = map[BookingStatus]string{
	BookingConfirmed: "已确认",
	BookingPending:   "待处理",
	BookingReleased:  "已释放",
}

func (s BookingStatus) Valid() bool {
	for _, candidate := range BookingStatuses {
		if candidate == s {
			return true
		}
	}
	return false
}

func (s BookingStatus) Text() string { return BookingStatusText[s] }

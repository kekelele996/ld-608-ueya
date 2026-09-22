package constants

// ResourceStatus values, duplicated in frontend constants/ResourceStatus.
const (
	ResourceAvailable   = "AVAILABLE"
	ResourceBooked      = "BOOKED"
	ResourceMaintenance = "MAINTENANCE"
	ResourceOffline     = "OFFLINE"
)

var ResourceStatus = []string{
	ResourceAvailable,
	ResourceBooked,
	ResourceMaintenance,
	ResourceOffline,
}

var ResourceStatusText = map[string]string{
	ResourceAvailable:   "可用",
	ResourceBooked:      "占用中",
	ResourceMaintenance: "维护中",
	ResourceOffline:     "离线",
}

func IsValidResourceStatus(value string) bool {
	for _, item := range ResourceStatus {
		if item == value {
			return true
		}
	}
	return false
}

// Resource booking status values (ResourceBooking.booking_status).
const (
	BookingConfirmed = "CONFIRMED"
	BookingPending   = "PENDING"
	BookingReleased  = "RELEASED"
)

var BookingStatus = []string{BookingConfirmed, BookingPending, BookingReleased}

var BookingStatusText = map[string]string{
	BookingConfirmed: "已确认",
	BookingPending:   "待处理",
	BookingReleased:  "已释放",
}

// Ground task lifecycle status values.
const (
	TaskPlanned   = "PLANNED"
	TaskAccepted  = "ACCEPTED"
	TaskCompleted = "COMPLETED"
	TaskBlocked   = "BLOCKED"
)

var GroundTaskStatus = []string{TaskPlanned, TaskAccepted, TaskCompleted, TaskBlocked}

var GroundTaskStatusText = map[string]string{
	TaskPlanned:   "待签收",
	TaskAccepted:  "已签收",
	TaskCompleted: "已完成",
	TaskBlocked:   "阻塞",
}

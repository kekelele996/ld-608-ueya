package constants

// ResourceStatus describes whether a ground resource can be booked.
// Duplicated in frontend/src/constants/ResourceStatus.ts.
type ResourceStatus string

const (
	ResourceAvailable   ResourceStatus = "AVAILABLE"
	ResourceBooked      ResourceStatus = "BOOKED"
	ResourceMaintenance ResourceStatus = "MAINTENANCE"
	ResourceOffline     ResourceStatus = "OFFLINE"
)

var ResourceStatuses = []ResourceStatus{
	ResourceAvailable, ResourceBooked, ResourceMaintenance, ResourceOffline,
}

var ResourceStatusText = map[ResourceStatus]string{
	ResourceAvailable:   "可用",
	ResourceBooked:      "已占用",
	ResourceMaintenance: "维护中",
	ResourceOffline:     "离线",
}

func (s ResourceStatus) Valid() bool {
	for _, candidate := range ResourceStatuses {
		if candidate == s {
			return true
		}
	}
	return false
}

func (s ResourceStatus) Bookable() bool { return s == ResourceAvailable }

func (s ResourceStatus) Text() string { return ResourceStatusText[s] }

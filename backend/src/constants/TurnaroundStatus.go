package constants

// TurnaroundStatus values, duplicated in frontend constants/TurnaroundStatus.
const (
	StatusArriving  = "ARRIVING"
	StatusOnStand   = "ON_STAND"
	StatusInService = "IN_SERVICE"
	StatusReady     = "READY"
	StatusDeparted  = "DEPARTED"
	StatusDelayed   = "DELAYED"
)

var TurnaroundStatus = []string{
	StatusArriving,
	StatusOnStand,
	StatusInService,
	StatusReady,
	StatusDeparted,
	StatusDelayed,
}

var TurnaroundStatusText = map[string]string{
	StatusArriving:  "即将到站",
	StatusOnStand:   "已靠桥",
	StatusInService: "保障中",
	StatusReady:     "待放行",
	StatusDeparted:  "已离港",
	StatusDelayed:   "延误",
}

func IsValidTurnaroundStatus(value string) bool {
	for _, item := range TurnaroundStatus {
		if item == value {
			return true
		}
	}
	return false
}

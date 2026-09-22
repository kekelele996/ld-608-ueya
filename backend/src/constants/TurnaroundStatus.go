package constants

// TurnaroundStatus is the lifecycle of a flight turnaround.
// Duplicated in frontend/src/constants/TurnaroundStatus.ts.
type TurnaroundStatus string

const (
	StatusArriving   TurnaroundStatus = "ARRIVING"
	StatusOnStand    TurnaroundStatus = "ON_STAND"
	StatusInService  TurnaroundStatus = "IN_SERVICE"
	StatusReady      TurnaroundStatus = "READY"
	StatusDeparted   TurnaroundStatus = "DEPARTED"
	StatusDelayed    TurnaroundStatus = "DELAYED"
)

var TurnaroundStatuses = []TurnaroundStatus{
	StatusArriving, StatusOnStand, StatusInService, StatusReady, StatusDeparted, StatusDelayed,
}

var TurnaroundStatusText = map[TurnaroundStatus]string{
	StatusArriving:  "即将到站",
	StatusOnStand:   "已上轮挡",
	StatusInService: "保障中",
	StatusReady:     "待放飞",
	StatusDeparted:  "已离港",
	StatusDelayed:   "延误",
}

func (s TurnaroundStatus) Valid() bool {
	for _, candidate := range TurnaroundStatuses {
		if candidate == s {
			return true
		}
	}
	return false
}

func (s TurnaroundStatus) Text() string { return TurnaroundStatusText[s] }

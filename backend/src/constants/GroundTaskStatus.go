package constants

// GroundTaskStatus tracks dispatch / sign-off lifecycle.
type GroundTaskStatus string

const (
	TaskPlanned   GroundTaskStatus = "PLANNED"   // 已派发未签收
	TaskSigned    GroundTaskStatus = "SIGNED"    // 班组已签收
	TaskFinished  GroundTaskStatus = "FINISHED"  // 已完成
	TaskBlocked   GroundTaskStatus = "BLOCKED"   // 阻塞
)

var GroundTaskStatuses = []GroundTaskStatus{
	TaskPlanned, TaskSigned, TaskFinished, TaskBlocked,
}

var GroundTaskStatusText = map[GroundTaskStatus]string{
	TaskPlanned:  "待签收",
	TaskSigned:   "已签收",
	TaskFinished: "已完成",
	TaskBlocked:  "阻塞",
}

// Acknowledged means the crew already signed for the task; delay replanning
// must keep these time points untouched.
func (s GroundTaskStatus) Acknowledged() bool { return s == TaskSigned || s == TaskFinished }
func (s GroundTaskStatus) Finished() bool     { return s == TaskFinished }

func (s GroundTaskStatus) Text() string { return GroundTaskStatusText[s] }

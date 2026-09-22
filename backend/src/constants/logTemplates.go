package constants

// Log action templates. Every write operation records one audit entry;
// field changes require updating both the template and its call site.
const (
	LogTurnaroundCreated = "航班过站 %s 已登记，机位 %s，状态 %s"
	LogTurnaroundStatus  = "航班过站 %s 状态由 %s 变更为 %s"
	LogTurnaroundRelease = "航班过站 %s 已放行离港"

	LogTaskPlanGenerated = "航班 %s 整批生成保障计划：%d 个任务、%d 条资源预约"
	LogTaskPlanRejected  = "航班 %s 保障计划校验失败，整批未保存：%d 项冲突"
	LogTaskAccepted      = "任务 %d（%s）由班组 %s 签收，计划时点保持 %s"
	LogTaskCompleted     = "任务 %d（%s）完成，实际完成时间 %s，预约 %d 已释放"
	LogTaskBlocked       = "任务 %d（%s）标记阻塞：%s"

	LogResourceCreated     = "保障资源 %s（%s）登记入库，归属班组 %s"
	LogResourceStatus      = "保障资源 %s 状态由 %s 变更为 %s"
	LogResourceMaintenance = "保障资源 %s 维护窗口登记：%s ~ %s"
	LogBookingAdjusted     = "预约 %d（资源 %s）调整时段为 %s ~ %s，状态 %s"

	LogBookingCreated  = "资源预约 %d 生成：资源 %s，%s ~ %s"
	LogBookingPending  = "预约 %d 被挤出确认时段，进入待处理：%s"
	LogBookingReleased = "资源预约 %d 已释放"

	LogDelayRegistered = "航班 %s 登记延误 +%d 分钟（%s），重排 %d 个未签收任务，%d 条预约待处理"
	LogDelayResolved   = "延误事件 %d 已关闭归因：%s"
)

// LogTemplates is consumed by audit middleware/tests and mirrors the
// frontend LOG_TEMPLATES registry. Each entity exposes >= 4 templates.
var LogTemplates = map[string][]string{
	"FlightTurnaround": {LogTurnaroundCreated, LogTurnaroundStatus, LogTurnaroundRelease, "航班过站导出"},
	"GroundTask":       {LogTaskPlanGenerated, LogTaskAccepted, LogTaskCompleted, LogTaskBlocked, LogTaskPlanRejected},
	"GroundResource":   {LogResourceCreated, LogResourceStatus, LogResourceMaintenance, LogBookingAdjusted},
	"ResourceBooking":  {LogBookingCreated, LogBookingPending, LogBookingReleased, "资源预约导出"},
	"DelayEvent":       {LogDelayRegistered, LogDelayResolved, "延误事件更新", "延误事件导出"},
}

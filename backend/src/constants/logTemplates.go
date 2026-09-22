package constants

// Audit actions for operation logging. Every write path emits one of these.
const (
	ActionTurnaroundCreate  = "TURNAROUND_CREATE"
	ActionTurnaroundArrive  = "TURNAROUND_ARRIVE"
	ActionTurnaroundRelease = "TURNAROUND_RELEASE"

	ActionTaskDispatch = "TASK_DISPATCH"
	ActionTaskSign     = "TASK_SIGN"
	ActionTaskFinish   = "TASK_FINISH"
	ActionTaskBlock    = "TASK_BLOCK"

	ActionResourceCreate     = "RESOURCE_CREATE"
	ActionResourceUpdate     = "RESOURCE_UPDATE"
	ActionResourceMaintenance = "RESOURCE_SET_MAINTENANCE"

	ActionBookingCreate   = "BOOKING_CREATE"
	ActionBookingAdjust   = "BOOKING_ADJUST"
	ActionBookingBump     = "BOOKING_BUMP"
	ActionBookingRelease  = "BOOKING_RELEASE"
	ActionBookingConfirm  = "BOOKING_CONFIRM"

	ActionDelayRegister = "DELAY_REGISTER"
	ActionDelayReplan   = "DELAY_REPLAN"
	ActionDelayResolve  = "DELAY_RESOLVE"

	ActionPlanGenerate = "PLAN_GENERATE"
)

// LogTemplates are centralized audit message templates. Changing an entity
// field requires touching these templates and every call site in services.
var LogTemplates = map[string]string{
	ActionTurnaroundCreate:  "航班过站登记：航班 %s 机位 %s 计划到站 %s",
	ActionTurnaroundArrive:  "航班 %s 到站确认，状态 %s -> %s",
	ActionTurnaroundRelease: "航班 %s 放行，保障任务完成率 %d/%d",

	ActionTaskDispatch: "任务批量派发：航班 %s 生成 %d 项任务",
	ActionTaskSign:     "任务签收：任务 #%d（%s）由班组 %s 签收，计划开始 %s",
	ActionTaskFinish:   "任务完成：任务 #%d（%s）实际完成 %s",
	ActionTaskBlock:    "任务阻塞：任务 #%d（%s）阻塞原因 %s",

	ActionResourceCreate:      "资源台账新增：资源 %s（%s）归属 %s",
	ActionResourceUpdate:      "资源信息更新：资源 %s 字段 %s 变更",
	ActionResourceMaintenance: "资源维护登记：资源 %s 状态 -> %s，维护到期 %s",

	ActionBookingCreate:  "资源预约生成：资源 %s 预约 %s ~ %s 关联任务 #%d",
	ActionBookingAdjust:  "预约时段调整：预约 #%d 资源 %s 调整为 %s ~ %s，结果 %s",
	ActionBookingBump:    "预约挤出：预约 #%d 资源 %s 原因 %s",
	ActionBookingRelease: "预约释放：预约 #%d 资源 %s",
	ActionBookingConfirm: "预约确认：预约 #%d 资源 %s %s ~ %s",

	ActionDelayRegister: "延误登记：航班 %s 延误 %d 分钟（%s）责任 %s",
	ActionDelayReplan:   "延误重排：航班 %s 偏移 %d 分钟，调整任务 %d 项，挤出预约 %d 项",
	ActionDelayResolve:  "延误关闭：事件 #%d 航班 %s",

	ActionPlanGenerate: "保障计划生成：航班 %s 校验 %d 项任务/预约",
}

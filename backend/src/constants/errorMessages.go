package constants

// Error message templates. Services wrap per-entity errors and controllers
// wrap them again; neither layer is allowed to swallow errors globally.
const (
	MsgAuthRequired      = "缺少登录令牌，请先登录"
	MsgRBACDenied        = "当前角色无权执行该操作"
	MsgValidationFailed  = "表单字段缺失或格式错误：%s"
	MsgRateLimited       = "请求过于频繁，请稍后再试"
	MsgNotFound          = "%s不存在：%v"
	MsgInternal          = "服务内部错误"
	MsgPlanConflict      = "整批保障计划校验失败，任务与预约均未保存（%d 项冲突）"
	MsgBookingConflict   = "资源 %s 在 %s ~ %s 已被预约 %d 占用"
	MsgResourceOffline   = "资源 %s 当前处于离线状态"
	MsgResourceMaint     = "资源 %s 处于维护期（至 %s），与 %s ~ %s 重叠"
	MsgNoResource        = "任务 %s 需要 %s 类型资源，但无可用资源"
	MsgDeadlineConflict  = "任务 %s 的截止时间 %s 晚于离港时间 %s"
	MsgTaskOrder         = "任务 %s 计划开始 %s 早于到站时间 %s"
	MsgAlreadyExists     = "航班 %s 已生成过保障计划，如需重排请先登记延误或释放旧计划"
	MsgInvalidTransition = "任务 %d 当前状态 %s 不允许执行 %s"
)

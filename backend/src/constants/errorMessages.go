package constants

// Error message templates are composed by services / controllers separately;
// nobody swallows every error at a single global layer.
const (
	MsgAuthRequired      = "缺少登录凭证，请先登录"
	MsgAuthInvalid       = "登录凭证无效或已过期"
	MsgRBACDenied        = "当前角色无权执行该操作"
	MsgValidationFailed  = "请求参数校验失败：%s"
	MsgRateLimited       = "操作过于频繁，请稍后再试"
	MsgNotFound          = "%s不存在或已被释放"
	MsgPlanConflict      = "保障计划整批未保存：存在 %d 项冲突"
	MsgBookingConflict   = "资源预约失败：%s"
	MsgInvalidTransition = "当前状态不允许该操作：%s"
	MsgInternal          = "服务内部错误，请稍后重试"
)

// Item conflict message templates, indexed by item-level conflict codes.
var ConflictMessages = map[string]string{
	ItemResourceOverlap:     "资源 %s 在 %s ~ %s 已被其他任务占用（任务 #%d）",
	ItemResourceMaintenance: "资源 %s 处于维护中，维护到期 %s",
	ItemResourceOffline:     "资源 %s 已离线，不可预约",
	ItemResourceWindow:      "资源 %s 的可用时段 %s ~ %s 不覆盖任务窗口 %s ~ %s",
	ItemDeadlineConflict:    "任务 %s 计划结束 %s 晚于截止时间 %s",
	ItemNoResource:          "任务 %s 没有可匹配的 %s 资源",
	ItemBookingBumped:       "延误重排后资源 %s 在 %s ~ %s 被占用，预约挤出待处理",
}

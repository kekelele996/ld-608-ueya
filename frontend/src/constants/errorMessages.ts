// Error message templates shared by api client, stores and pages. Services
// wrap per-entity errors; the UI renders backend message plus conflict list.
export const ERROR_MESSAGES: Record<string, string> = {
  AUTH_REQUIRED: "请先登录后再继续操作",
  RBAC_DENIED: "当前角色没有执行该动作的权限",
  VALIDATION_FAILED: "表单字段缺失或格式错误",
  RATE_LIMITED: "请求过于频繁，请稍后再试",
  NOT_FOUND: "记录不存在",
  PLAN_CONFLICT: "整批保障计划校验失败，任务与预约均未保存",
  BOOKING_CONFLICT: "调整时段仍有冲突，预约保持待处理",
  INVALID_TRANSITION: "当前状态不允许执行该动作",
  INTERNAL_ERROR: "服务内部错误",
};

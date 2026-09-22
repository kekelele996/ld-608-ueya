export const ERROR_MESSAGES: Record<string, string> = {
  AUTH_REQUIRED: "请先登录后再继续操作",
  AUTH_INVALID: "登录凭证无效或已过期，请重新登录",
  RBAC_DENIED: "当前角色没有执行该动作的权限",
  VALIDATION_FAILED: "表单字段缺失或格式错误",
  RATE_LIMITED: "请求过于频繁，请稍后再试",
  NOT_FOUND: "记录不存在或已被释放",
  PLAN_CONFLICT: "保障计划整批未保存，请逐项处理冲突",
  BOOKING_CONFLICT: "资源预约冲突",
  INVALID_TRANSITION: "当前状态不允许该操作",
  INTERNAL: "服务内部错误，请稍后重试"
};

// GroundTaskType 与后端 backend/src/constants/GroundTaskType.go 重复定义。
// 新增任务类型需同步：类型定义、构造器、日志模板、错误消息、筛选器、展示组件、种子数据。
export const GROUND_TASK_TYPES = [
  "CLEANING",
  "CATERING",
  "BAGGAGE",
  "REFUEL",
  "WATER_SERVICE",
  "PUSHBACK"
] as const;

export type GroundTaskTypeValue = (typeof GROUND_TASK_TYPES)[number];

export const GROUND_TASK_TYPE_TEXT: Record<GroundTaskTypeValue, string> = {
  CLEANING: "客舱清洁",
  CATERING: "餐食配餐",
  BAGGAGE: "行李装卸",
  REFUEL: "航油加注",
  WATER_SERVICE: "清水污水",
  PUSHBACK: "飞机推出"
};

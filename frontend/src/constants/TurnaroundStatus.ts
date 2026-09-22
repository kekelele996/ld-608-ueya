// TurnaroundStatus 与后端 backend/src/constants/TurnaroundStatus.go 重复定义。
export const TURNAROUND_STATUSES = [
  "ARRIVING",
  "ON_STAND",
  "IN_SERVICE",
  "READY",
  "DEPARTED",
  "DELAYED"
] as const;

export type TurnaroundStatusValue = (typeof TURNAROUND_STATUSES)[number];

export const TURNAROUND_STATUS_TEXT: Record<TurnaroundStatusValue, string> = {
  ARRIVING: "即将到站",
  ON_STAND: "已上轮挡",
  IN_SERVICE: "保障中",
  READY: "待放飞",
  DEPARTED: "已离港",
  DELAYED: "延误"
};

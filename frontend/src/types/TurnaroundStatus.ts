export type TurnaroundStatus =
  | "ARRIVING"
  | "ON_STAND"
  | "IN_SERVICE"
  | "READY"
  | "DEPARTED"
  | "DELAYED";

export const TURNAROUND_STATUSES: TurnaroundStatus[] = [
  "ARRIVING",
  "ON_STAND",
  "IN_SERVICE",
  "READY",
  "DEPARTED",
  "DELAYED",
];

export const TURNAROUND_STATUS_TEXT: Record<TurnaroundStatus, string> = {
  ARRIVING: "即将到站",
  ON_STAND: "已靠桥",
  IN_SERVICE: "保障中",
  READY: "待放行",
  DEPARTED: "已离港",
  DELAYED: "延误",
};

import type { TurnaroundStatus } from "../types/TurnaroundStatus";

export const TURNAROUND_STATUS_COLOR: Record<TurnaroundStatus, string> = {
  ARRIVING: "default",
  ON_STAND: "blue",
  IN_SERVICE: "processing",
  READY: "warning",
  DEPARTED: "success",
  DELAYED: "error",
};

export const TURNAROUND_STATUS_FILTERS = [
  { value: "", label: "全部状态" },
  { value: "ON_STAND", label: "已靠桥" },
  { value: "IN_SERVICE", label: "保障中" },
  { value: "DELAYED", label: "延误" },
  { value: "READY", label: "待放行" },
  { value: "DEPARTED", label: "已离港" },
];

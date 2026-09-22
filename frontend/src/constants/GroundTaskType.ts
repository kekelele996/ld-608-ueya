// GroundTaskType enum mirrored from types/GroundTaskType.
// Appears in: constructors, logTemplates, errorMessages, filters and
// display components; change a value here and across all of those layers.
import type { GroundTaskType } from "../types/GroundTaskType";

export const GROUND_TASK_TYPE_OPTIONS: { value: GroundTaskType; label: string }[] = [
  { value: "CLEANING", label: "客舱清洁" },
  { value: "CATERING", label: "航食配送" },
  { value: "BAGGAGE", label: "行李装卸" },
  { value: "REFUEL", label: "航油加注" },
  { value: "WATER_SERVICE", label: "清水污水" },
  { value: "PUSHBACK", label: "飞机推出" },
];

export const GROUND_TASK_TYPE_COLOR: Record<GroundTaskType, string> = {
  CLEANING: "cyan",
  CATERING: "gold",
  BAGGAGE: "geekblue",
  REFUEL: "orange",
  WATER_SERVICE: "blue",
  PUSHBACK: "purple",
};

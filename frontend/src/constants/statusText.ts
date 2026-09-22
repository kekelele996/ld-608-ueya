import { GROUND_TASK_TYPE_TEXT } from "./GroundTaskType";
import { TURNAROUND_STATUS_TEXT } from "./TurnaroundStatus";
import { RESOURCE_STATUS_TEXT } from "./ResourceStatus";
import { GROUND_TASK_STATUS_TEXT } from "./GroundTaskStatus";
import { BOOKING_STATUS_TEXT } from "./BookingStatus";
import { DELAY_TYPE_TEXT } from "./DelayType";

// 汇总状态文案，供 formatters 与共享组件依赖。
export const STATUS_TEXT = {
  GroundTaskType: GROUND_TASK_TYPE_TEXT,
  TurnaroundStatus: TURNAROUND_STATUS_TEXT,
  ResourceStatus: RESOURCE_STATUS_TEXT,
  GroundTaskStatus: GROUND_TASK_STATUS_TEXT,
  BookingStatus: BOOKING_STATUS_TEXT,
  DelayType: DELAY_TYPE_TEXT
};

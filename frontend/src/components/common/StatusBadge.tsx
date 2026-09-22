import React from "react";
import { Tag } from "antd";
import { TURNAROUND_STATUS_TEXT } from "../../types/TurnaroundStatus";
import { RESOURCE_STATUS_TEXT, BOOKING_STATUS_TEXT, GROUND_TASK_STATUS_TEXT } from "../../types/ResourceStatus";
import { GROUND_TASK_TYPE_TEXT } from "../../types/GroundTaskType";
import { TURNAROUND_STATUS_COLOR } from "../../constants/TurnaroundStatus";
import { RESOURCE_STATUS_COLOR } from "../../constants/ResourceStatus";
import {
  BOOKING_STATUS_COLOR,
  GROUND_TASK_STATUS_COLOR,
} from "../../constants/statusText";
import { GROUND_TASK_TYPE_COLOR } from "../../constants/GroundTaskType";

export type BadgeKind =
  | "turnaround"
  | "resource"
  | "booking"
  | "task"
  | "taskType";

interface StatusBadgeProps {
  kind: BadgeKind;
  value: string;
}

// StatusBadge is shared across dashboard, turnarounds, tasks and resources.
export const StatusBadge: React.FC<StatusBadgeProps> = ({ kind, value }) => {
  let text = value;
  let color = "default";
  switch (kind) {
    case "turnaround":
      text = TURNAROUND_STATUS_TEXT[value as keyof typeof TURNAROUND_STATUS_TEXT] ?? value;
      color = TURNAROUND_STATUS_COLOR[value as keyof typeof TURNAROUND_STATUS_COLOR] ?? "default";
      break;
    case "resource":
      text = RESOURCE_STATUS_TEXT[value as keyof typeof RESOURCE_STATUS_TEXT] ?? value;
      color = RESOURCE_STATUS_COLOR[value as keyof typeof RESOURCE_STATUS_COLOR] ?? "default";
      break;
    case "booking":
      text = BOOKING_STATUS_TEXT[value as keyof typeof BOOKING_STATUS_TEXT] ?? value;
      color = BOOKING_STATUS_COLOR[value as keyof typeof BOOKING_STATUS_COLOR] ?? "default";
      break;
    case "task":
      text = GROUND_TASK_STATUS_TEXT[value as keyof typeof GROUND_TASK_STATUS_TEXT] ?? value;
      color = GROUND_TASK_STATUS_COLOR[value] ?? "default";
      break;
    case "taskType":
      text = GROUND_TASK_TYPE_TEXT[value as keyof typeof GROUND_TASK_TYPE_TEXT] ?? value;
      color = GROUND_TASK_TYPE_COLOR[value as keyof typeof GROUND_TASK_TYPE_COLOR] ?? "default";
      break;
  }
  return <Tag color={color}>{text}</Tag>;
};

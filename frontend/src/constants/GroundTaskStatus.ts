export const GROUND_TASK_STATUSES = ["PLANNED", "SIGNED", "FINISHED", "BLOCKED"] as const;
export type GroundTaskStatusValue = (typeof GROUND_TASK_STATUSES)[number];

export const GROUND_TASK_STATUS_TEXT: Record<GroundTaskStatusValue, string> = {
  PLANNED: "待签收",
  SIGNED: "已签收",
  FINISHED: "已完成",
  BLOCKED: "阻塞"
};

// 已签收/已完成任务在延误重排时保持原计划时点。
export const ACKNOWLEDGED_TASK_STATUSES: GroundTaskStatusValue[] = ["SIGNED", "FINISHED"];

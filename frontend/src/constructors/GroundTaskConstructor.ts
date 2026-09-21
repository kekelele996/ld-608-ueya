import type { GroundTask } from "../types/GroundTask";

export const createDefaultGroundTask = (overrides: Partial<GroundTask> = {}): GroundTask => ({
  id: 1 as never,
  turnaround_id: 1 as never,
  task_type: "CATERING" as never,
  team_id: 1 as never,
  planned_start: "planned start 1" as never,
  deadline: "deadline 1" as never,
  actual_finish: "actual finish 1" as never,
  status: "ON_STAND" as never,
  blocker_note: "blocker note 1" as never,
  ...overrides
});

export const createGroundTaskForm = createDefaultGroundTask;
export const createGroundTaskResponse = createDefaultGroundTask;

import type { GroundTask } from "../types/GroundTask";

export function createDefaultGroundTask(overrides: Partial<GroundTask> = {}): GroundTask {
  return {
    id: 0,
    turnaround_id: 0,
    task_type: "CLEANING",
    team_id: "",
    planned_start: "",
    planned_end: "",
    deadline: "",
    actual_finish: null,
    status: "PLANNED",
    blocker_note: "",
    signed_by: "",
    ...overrides
  };
}

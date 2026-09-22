import type { GroundTask } from "../types/entities";

// Builds an empty task row used as a fallback skeleton by task pages.
export const createEmptyTask = (): GroundTask => ({
  id: 0,
  turnaround_id: 0,
  task_type: "CLEANING",
  team_id: "",
  planned_start: "",
  deadline: "",
  actual_finish: null,
  status: "PLANNED",
  blocker_note: "",
  signed_at: null,
});

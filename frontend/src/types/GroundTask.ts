import type { GroundTaskTypeValue } from "../constants/GroundTaskType";
import type { GroundTaskStatusValue } from "../constants/GroundTaskStatus";

export interface GroundTask {
  id: number;
  turnaround_id: number;
  task_type: GroundTaskTypeValue;
  team_id: string;
  planned_start: string;
  planned_end: string;
  deadline: string;
  actual_finish: string | null;
  status: GroundTaskStatusValue;
  blocker_note: string;
  signed_by: string;
  created_at?: string;
  updated_at?: string;
}

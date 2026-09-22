import type { DelayTypeValue } from "../constants/DelayType";

export interface DelayEvent {
  id: number;
  turnaround_id: number;
  delay_type: DelayTypeValue;
  minutes: number;
  root_cause: string;
  responsibility_team: string;
  resolved_at: string | null;
  created_at?: string;
}

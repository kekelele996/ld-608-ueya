export interface DelayEvent {
  id: number;
  turnaround_id: number;
  delay_type: string;
  minutes: number;
  root_cause: string;
  responsibility_team: string;
  resolved_at: string;
}

import type { DelayEvent } from "../types/entities";

export const createEmptyDelay = (): DelayEvent => ({
  id: 0,
  turnaround_id: 0,
  delay_type: "LATE_ARRIVAL",
  minutes: 0,
  root_cause: "",
  responsibility_team: "",
  resolved_at: null,
  created_at: "",
});

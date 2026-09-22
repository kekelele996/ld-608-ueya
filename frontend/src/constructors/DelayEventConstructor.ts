import type { DelayEvent } from "../types/DelayEvent";

export function createDefaultDelayEvent(overrides: Partial<DelayEvent> = {}): DelayEvent {
  return {
    id: 0,
    turnaround_id: 0,
    delay_type: "ATC",
    minutes: 30,
    root_cause: "",
    responsibility_team: "",
    resolved_at: null,
    ...overrides
  };
}

export function createDelayForm(turnaroundId: number) {
  return {
    turnaround_id: turnaroundId,
    delay_type: "ATC" as DelayEvent["delay_type"],
    minutes: 30,
    root_cause: "",
    responsibility_team: ""
  };
}

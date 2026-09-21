import type { DelayEvent } from "../types/DelayEvent";

export const createDefaultDelayEvent = (overrides: Partial<DelayEvent> = {}): DelayEvent => ({
  id: 1 as never,
  turnaround_id: 1 as never,
  delay_type: "CATERING" as never,
  minutes: "minutes 1" as never,
  root_cause: "root cause 1" as never,
  responsibility_team: "responsibility team 1" as never,
  resolved_at: "2026-06-11T09:00:00Z" as never,
  ...overrides
});

export const createDelayEventForm = createDefaultDelayEvent;
export const createDelayEventResponse = createDefaultDelayEvent;

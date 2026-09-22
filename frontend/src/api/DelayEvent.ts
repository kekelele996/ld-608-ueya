import { request } from "./client";
import type { DelayEvent } from "../types/DelayEvent";

export interface DelayImpact {
  delay_type: string;
  event_count: number;
  total_minutes: number;
  responsibility_team: string;
}

export const delayEventApi = {
  list: () => request<DelayEvent[]>("/delay-events"),
  impacts: () => request<DelayImpact[]>("/delay-events/impacts"),
  resolve: (id: number) =>
    request<DelayEvent>(`/delay-events/${id}/resolve`, { method: "POST", body: {} })
};

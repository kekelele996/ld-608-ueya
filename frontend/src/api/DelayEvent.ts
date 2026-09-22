import { request, type ListEnvelope, RequestError } from "./client";
import type { DelayEvent, PlanResult } from "../types/entities";

export interface RegisterDelayPayload {
  delay_type: string;
  minutes: number;
  root_cause?: string;
  responsibility_team?: string;
}

export const listDelays = (params?: { page?: number; page_size?: number }) =>
  request<ListEnvelope<DelayEvent>>("/delay-events", { query: params });

export const registerDelay = (turnaroundId: number, payload: RegisterDelayPayload) =>
  request<PlanResult>(`/turnarounds/${turnaroundId}/delays`, {
    method: "POST",
    body: payload,
  });

export const resolveDelay = (id: number) =>
  request<{ id: number; resolved_at: string }>(`/delay-events/${id}/resolve`, { method: "POST" });

// isPlanConflict narrows the 409 envelope returned by plan/delay/adjust.
export const isPlanConflict = (error: unknown): error is RequestError =>
  error instanceof RequestError &&
  (error.code === "PLAN_CONFLICT" || error.code === "BOOKING_CONFLICT");

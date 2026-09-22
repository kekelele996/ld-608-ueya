import { request, type ListEnvelope } from "./client";
import type { FlightTurnaround, TurnaroundDetail, PlanResult } from "../types/entities";

export interface RegisterTurnaroundPayload {
  flight_no: string;
  aircraft_reg: string;
  stand_no: string;
  arrival_time: string;
  departure_time: string;
}

export const listTurnarounds = (params?: { page?: number; page_size?: number }) =>
  request<ListEnvelope<FlightTurnaround>>("/turnarounds", { query: params });

export const getTurnaround = (id: number) =>
  request<TurnaroundDetail>(`/turnarounds/${id}`);

export const registerTurnaround = (payload: RegisterTurnaroundPayload) =>
  request<FlightTurnaround>("/turnarounds", { method: "POST", body: payload });

// On conflict the backend returns 409 with { result } carrying the unsaved
// preview; the caller catches RequestError to display itemized reasons.
export const generatePlan = (id: number, taskTypes?: string[]) =>
  request<PlanResult>(`/turnarounds/${id}/plan`, {
    method: "POST",
    body: { task_types: taskTypes ?? [] },
  });

import { request } from "./client";
import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { GroundTask } from "../types/GroundTask";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { DelayEvent } from "../types/DelayEvent";
import type {
  PlanConflictResult,
  ReplanResult,
  TurnaroundDetail
} from "../types/api";

export const turnaroundApi = {
  list: () => request<FlightTurnaround[]>("/turnarounds"),
  detail: (id: number) => request<TurnaroundDetail>(`/turnarounds/${id}`),
  create: (body: {
    flight_no: string;
    aircraft_reg: string;
    stand_no: string;
    arrival_time: string;
    departure_time: string;
  }) => request<FlightTurnaround>("/turnarounds", { method: "POST", body }),
  arrive: (id: number) =>
    request<FlightTurnaround>(`/turnarounds/${id}/arrive`, { method: "POST", body: {} }),
  // 整批生成：冲突时 ApiError.data 携带 PlanConflictResult（逐项原因）。
  generatePlan: (id: number) =>
    request<PlanConflictResult>(`/turnarounds/${id}/generate-plan`, {
      method: "POST",
      body: {},
      silent: true
    }),
  release: (id: number) =>
    request<FlightTurnaround>(`/turnarounds/${id}/release`, { method: "POST", body: {} }),
  registerDelay: (
    id: number,
    body: { delay_type: string; minutes: number; root_cause: string; responsibility_team: string }
  ) =>
    request<ReplanResult>(`/turnarounds/${id}/delays`, {
      method: "POST",
      body
    })
};

export { type GroundTask, type ResourceBooking, type DelayEvent };

import { request } from "./client";
import type { GroundTask } from "../types/GroundTask";

export const groundTaskApi = {
  list: () => request<GroundTask[]>("/tasks"),
  sign: (id: number, operator?: string) =>
    request<GroundTask>(`/tasks/${id}/sign`, { method: "POST", body: { operator } }),
  finish: (id: number, note?: string) =>
    request<GroundTask>(`/tasks/${id}/finish`, { method: "POST", body: { note } }),
  block: (id: number, note: string) =>
    request<GroundTask>(`/tasks/${id}/block`, { method: "POST", body: { note } })
};

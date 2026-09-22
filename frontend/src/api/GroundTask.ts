import { request, type ListEnvelope } from "./client";
import type { GroundTask } from "../types/entities";

export const listTasks = (params?: { page?: number; page_size?: number; status?: string }) =>
  request<ListEnvelope<GroundTask>>("/ground-tasks", { query: params });

export const acceptTask = (id: number) =>
  request<GroundTask>(`/ground-tasks/${id}/accept`, { method: "POST" });

export const completeTask = (id: number) =>
  request<GroundTask>(`/ground-tasks/${id}/complete`, { method: "POST" });

export const blockTask = (id: number, blocker_note: string) =>
  request<GroundTask>(`/ground-tasks/${id}/block`, { method: "POST", body: { blocker_note } });

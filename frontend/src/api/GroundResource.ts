import { request } from "./client";
import type { GroundResource } from "../types/GroundResource";

export const groundResourceApi = {
  list: () => request<GroundResource[]>("/resources"),
  setStatus: (id: number, body: { status: string; maintenance_due_at?: string | null }) =>
    request<GroundResource>(`/resources/${id}/status`, { method: "POST", body })
};

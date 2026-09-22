import { request, type ListEnvelope } from "./client";
import type { GroundResource, ResourceBooking } from "../types/entities";

export const listResources = (params?: { status?: string; type?: string }) =>
  request<ListEnvelope<GroundResource>>("/ground-resources", { query: params });

export const updateResourceStatus = (
  id: number,
  payload: { status?: string; maintenance_start?: string; maintenance_end?: string },
) => request<GroundResource>(`/ground-resources/${id}/status`, { method: "PATCH", body: payload });

export interface CalendarResponse {
  items: ResourceBooking[];
  resources: GroundResource[];
}

export const getCalendar = (from?: string, to?: string) =>
  request<CalendarResponse>("/resource-calendar", {
    query: { from, to },
  });

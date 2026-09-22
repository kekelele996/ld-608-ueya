import { request } from "./client";
import type { ResourceBooking } from "../types/ResourceBooking";

export const resourceBookingApi = {
  list: () => request<ResourceBooking[]>("/resource-bookings"),
  listPending: () => request<ResourceBooking[]>("/resource-bookings/pending"),
  // 调整时段；冲突时后端不回滚动作，而是把预约置为 PENDING 并返回冲突原因。
  adjust: (
    id: number,
    body: { resource_id: number; start_time: string; end_time: string }
  ) =>
    request<ResourceBooking>(`/resource-bookings/${id}/adjust`, {
      method: "POST",
      body
    }),
  release: (id: number) =>
    request<ResourceBooking>(`/resource-bookings/${id}/release`, {
      method: "POST",
      body: {}
    })
};

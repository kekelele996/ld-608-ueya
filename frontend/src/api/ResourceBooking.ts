import { request, type ListEnvelope } from "./client";
import type { ResourceBooking } from "../types/entities";

export const listBookings = (params?: { status?: string }) =>
  request<ListEnvelope<ResourceBooking>>("/resource-bookings", { query: params });

export interface AdjustResult {
  saved: boolean;
  booking_id: number;
  conflicts: unknown[];
}

// Adjust rejects with 409 + details (itemized conflicts) when the new window
// is unusable; the booking stays PENDING in that case.
export const adjustBooking = (id: number, start_time: string, end_time: string) =>
  request<AdjustResult>(`/resource-bookings/${id}/adjust`, {
    method: "POST",
    body: { start_time, end_time },
  });

export const releaseBooking = (id: number) =>
  request<{ id: number; booking_status: string }>(`/resource-bookings/${id}/release`, {
    method: "POST",
  });

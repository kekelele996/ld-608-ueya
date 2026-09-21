import type { ResourceBooking } from "../types/ResourceBooking";

export const createDefaultResourceBooking = (overrides: Partial<ResourceBooking> = {}): ResourceBooking => ({
  id: 1 as never,
  resource_id: 1 as never,
  turnaround_id: 1 as never,
  task_id: 1 as never,
  start_time: "2026-06-11T09:00:00Z" as never,
  end_time: "2026-06-11T09:00:00Z" as never,
  booking_status: "ON_STAND" as never,
  conflict_reason: "conflict reason 1" as never,
  ...overrides
});

export const createResourceBookingForm = createDefaultResourceBooking;
export const createResourceBookingResponse = createDefaultResourceBooking;

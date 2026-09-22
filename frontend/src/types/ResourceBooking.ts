import type { BookingStatusValue } from "../constants/BookingStatus";

export interface ResourceBooking {
  id: number;
  resource_id: number;
  turnaround_id: number;
  task_id: number | null;
  start_time: string;
  end_time: string;
  booking_status: BookingStatusValue;
  conflict_reason: string;
}

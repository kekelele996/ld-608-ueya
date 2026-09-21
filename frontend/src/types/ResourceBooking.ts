export interface ResourceBooking {
  id: number;
  resource_id: number;
  turnaround_id: number;
  task_id: number;
  start_time: string;
  end_time: string;
  booking_status: string;
  conflict_reason: string;
}

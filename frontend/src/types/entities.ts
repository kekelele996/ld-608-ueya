export interface FlightTurnaround {
  id: number;
  flight_no: string;
  aircraft_reg: string;
  stand_no: string;
  arrival_time: string;
  departure_time: string;
  turnaround_status: string;
  status_text?: string;
  delay_reason: string;
  total_delay_minutes: number;
}

export interface GroundTask {
  id: number;
  turnaround_id: number;
  task_type: string;
  task_type_text?: string;
  team_id: string;
  planned_start: string;
  deadline: string;
  actual_finish: string | null;
  status: string;
  status_text?: string;
  blocker_note: string;
  signed_at: string | null;
}

export interface GroundResource {
  id: number;
  resource_code: string;
  resource_type: string;
  location: string;
  availability_status: string;
  status_text?: string;
  maintenance_due_at: string | null;
  maintenance_end: string | null;
  owner_team: string;
}

export interface ResourceBooking {
  id: number;
  resource_id: number;
  resource_code?: string;
  resource_type?: string;
  turnaround_id: number;
  task_id: number;
  start_time: string;
  end_time: string;
  booking_status: string;
  status_text?: string;
  conflict_code: string;
  conflict_reason: string;
}

export interface DelayEvent {
  id: number;
  turnaround_id: number;
  flight_no?: string;
  delay_type: string;
  delay_type_text?: string;
  minutes: number;
  root_cause: string;
  responsibility_team: string;
  resolved_at: string | null;
  created_at: string;
}

export interface ConflictItem {
  code: string;
  message: string;
  task_type?: string;
  resource_id?: number;
  resource_code?: string;
  booking_id?: number;
  start_time?: string;
  end_time?: string;
}

export interface PlanResult {
  saved: boolean;
  turnaround_id: number;
  tasks: GroundTask[];
  bookings: ResourceBooking[];
  conflicts: ConflictItem[];
  rescheduled_tasks?: number;
  frozen_tasks?: number;
  pending_bookings?: number;
  added_delay_minutes?: number;
}

export interface TurnaroundDetail extends FlightTurnaround {
  task_count: number;
  completed_count: number;
  accepted_count: number;
  pending_bookings: number;
  tasks: GroundTask[];
  bookings: ResourceBooking[];
}

export interface DashboardStats {
  active_turnarounds: number;
  open_tasks: number;
  pending_bookings: number;
  total_delay_minutes: number;
}

export interface AuditLogEntry {
  id: number;
  actor: string;
  role: string;
  action: string;
  target_type: string;
  target_id: string;
  created_at: string;
}

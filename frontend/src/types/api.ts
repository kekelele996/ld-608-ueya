import type { FlightTurnaround } from "./FlightTurnaround";
import type { GroundTask } from "./GroundTask";
import type { ResourceBooking } from "./ResourceBooking";
import type { DelayEvent } from "./DelayEvent";

// 单项冲突说明：整批拒绝时逐项返回，延误挤出时同样复用该结构。
export interface ConflictItem {
  code: string;
  message: string;
  task_type?: string;
  task_id?: number;
  booking_id?: number;
}

// 保障计划整批生成结果。
export interface PlanConflictResult {
  saved: boolean;
  items: ConflictItem[];
  task_count: number;
}

// 延误重排结果：保持原时点 / 被移动 / 被挤出的三类明细。
export interface ReplanResult {
  delay_event_id: number;
  shift_minutes: number;
  kept_tasks: Array<{
    task_id: number;
    task_type: string;
    status: string;
    planned_start: string;
    signed_by: string;
    reason: string;
  }>;
  moved_tasks: GroundTask[];
  moved_bookings: ResourceBooking[];
  bumped: ConflictItem[];
}

export interface TurnaroundDetail extends FlightTurnaround {
  tasks: GroundTask[];
  bookings: ResourceBooking[];
  delays: DelayEvent[];
}

export interface DashboardStats {
  turnaround_total: number;
  in_service: number;
  delayed: number;
  task_total: number;
  task_finished: number;
  task_overdue: number;
  task_blocked: number;
  pending_bookings: number;
  delay_minutes: number;
}

export interface AuditLog {
  id: number;
  actor: string;
  role: string;
  action: string;
  target_type: string;
  target_id: string;
  detail: string;
  created_at: string;
}

export interface UserInfo {
  id: number;
  username: string;
  name: string;
  role: UserRoleValue;
  team_id: string;
}

export type UserRoleValue = "DISPATCHER" | "CREW" | "RESOURCE" | "SUPERVISOR";

export interface LoginResponse {
  token: string;
  user: UserInfo;
}

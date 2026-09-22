export type ResourceStatus = "AVAILABLE" | "BOOKED" | "MAINTENANCE" | "OFFLINE";

export const RESOURCE_STATUSES: ResourceStatus[] = [
  "AVAILABLE",
  "BOOKED",
  "MAINTENANCE",
  "OFFLINE",
];

export const RESOURCE_STATUS_TEXT: Record<ResourceStatus, string> = {
  AVAILABLE: "可用",
  BOOKED: "占用中",
  MAINTENANCE: "维护中",
  OFFLINE: "离线",
};

export type BookingStatus = "CONFIRMED" | "PENDING" | "RELEASED";

export const BOOKING_STATUS_TEXT: Record<BookingStatus, string> = {
  CONFIRMED: "已确认",
  PENDING: "待处理",
  RELEASED: "已释放",
};

export type GroundTaskStatus = "PLANNED" | "ACCEPTED" | "COMPLETED" | "BLOCKED";

export const GROUND_TASK_STATUS_TEXT: Record<GroundTaskStatus, string> = {
  PLANNED: "待签收",
  ACCEPTED: "已签收",
  COMPLETED: "已完成",
  BLOCKED: "阻塞",
};

// ResourceStatus 与后端 backend/src/constants/ResourceStatus.go 重复定义。
export const RESOURCE_STATUSES = [
  "AVAILABLE",
  "BOOKED",
  "MAINTENANCE",
  "OFFLINE"
] as const;

export type ResourceStatusValue = (typeof RESOURCE_STATUSES)[number];

export const RESOURCE_STATUS_TEXT: Record<ResourceStatusValue, string> = {
  AVAILABLE: "可用",
  BOOKED: "已占用",
  MAINTENANCE: "维护中",
  OFFLINE: "离线"
};

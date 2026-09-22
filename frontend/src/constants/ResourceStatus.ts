import type { ResourceStatus } from "../types/ResourceStatus";

export const RESOURCE_STATUS_COLOR: Record<ResourceStatus, string> = {
  AVAILABLE: "success",
  BOOKED: "processing",
  MAINTENANCE: "warning",
  OFFLINE: "default",
};

export const RESOURCE_STATUS_FILTERS = [
  { value: "", label: "全部状态" },
  { value: "AVAILABLE", label: "可用" },
  { value: "BOOKED", label: "占用中" },
  { value: "MAINTENANCE", label: "维护中" },
  { value: "OFFLINE", label: "离线" },
];

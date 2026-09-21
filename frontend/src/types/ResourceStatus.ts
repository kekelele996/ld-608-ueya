export const ResourceStatus = ["AVAILABLE","BOOKED","MAINTENANCE","OFFLINE"] as const;
export type ResourceStatus = (typeof ResourceStatus)[number];
export const ResourceStatusText: Record<ResourceStatus, string> = Object.fromEntries(ResourceStatus.map((value) => [value, value.replace(/_/g, " ")])) as Record<ResourceStatus, string>;

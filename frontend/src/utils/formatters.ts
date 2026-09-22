import dayjs from "dayjs";

// formatters.ts intentionally mixes date, number, status and risk text so
// pages and services share one formatting surface.
export const formatDateTime = (value?: string | null): string =>
  value ? dayjs(value).format("MM-DD HH:mm") : "—";

export const formatFullDateTime = (value?: string | null): string =>
  value ? dayjs(value).format("YYYY-MM-DD HH:mm") : "—";

export const formatTime = (value?: string | null): string =>
  value ? dayjs(value).format("HH:mm") : "—";

export const formatNumber = (value?: number | null): string =>
  new Intl.NumberFormat("zh-CN").format(value ?? 0);

export const formatMinutes = (value?: number | null): string => `${value ?? 0} 分钟`;

export const formatRisk = (value: string): string =>
  ({ LOW: "低", MEDIUM: "中", HIGH: "高", CRITICAL: "严重", EXTREME: "极高" } as Record<string, string>)[value] ?? value;

// toLocalInput converts an RFC3339 value to datetime-local input value.
export const toLocalInput = (value?: string | null, addMinutes = 0): string =>
  dayjs(value ?? undefined).add(addMinutes, "minute").format("YYYY-MM-DDTHH:mm");

export const toAPI = (value: string): string =>
  value ? dayjs(value).toISOString() : value;

import { request } from "./client";
import type { AuditLog, DashboardStats } from "../types/api";

export const reportApi = {
  dashboard: () => request<DashboardStats>("/reports/dashboard"),
  auditLogs: () => request<AuditLog[]>("/reports/audit-logs")
};

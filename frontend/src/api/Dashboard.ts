import { request } from "./client";
import type { DashboardStats, AuditLogEntry, ResourceBooking } from "../types/entities";
import type { ListEnvelope } from "./client";

export const getDashboardStats = () => request<DashboardStats>("/dashboard/stats");

export const getPendingBookings = () =>
  request<{ items: ResourceBooking[]; total: number }>("/dashboard/pending-bookings");

export const getAuditLogs = (params?: { page?: number; page_size?: number }) =>
  request<ListEnvelope<AuditLogEntry>>("/audit-logs", { query: params });

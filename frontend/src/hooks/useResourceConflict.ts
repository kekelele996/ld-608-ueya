import { useMemo } from "react";
import type { ResourceBooking } from "../types/ResourceBooking";

// useResourceConflict：从预约列表中筛出待处理（被挤出/改期失败）项，
// 并按资源聚合；资源页与过站详情共用。
export function useResourceConflict(bookings: ResourceBooking[] = []) {
  return useMemo(() => {
    const pending = bookings.filter((b) => b.booking_status === "PENDING");
    const byResource = pending.reduce<Record<number, ResourceBooking[]>>((acc, b) => {
      (acc[b.resource_id] ??= []).push(b);
      return acc;
    }, {});
    const hasConflict = pending.length > 0;
    return { pending, byResource, hasConflict, count: pending.length };
  }, [bookings]);
}

import { useMemo } from "react";
import type { GroundResource, ResourceBooking } from "../types/entities";
import { formatTime } from "../utils/formatters";

// useResourceConflict computes per-resource / per-booking conflict hints
// for ResourceCalendar and the pending booking worklist.
export function useResourceConflict(
  bookings: ResourceBooking[],
  resources?: GroundResource[],
) {
  return useMemo(() => {
    const pending = bookings.filter((b) => b.booking_status === "PENDING");
    const pendingIds = new Set(pending.map((b) => b.id));
    const resourceMap = new Map<number, GroundResource>((resources ?? []).map((r) => [r.id, r]));
    const summaries = pending.map((b) => ({
      booking: b,
      resource: resourceMap.get(b.resource_id),
      window: `${formatTime(b.start_time)} ~ ${formatTime(b.end_time)}`,
    }));
    return { pending, pendingIds, summaries };
  }, [bookings, resources]);
}

import { useMemo } from "react";
import type { TurnaroundDetail } from "../types/entities";

// useTurnaroundProgress computes task completion / pending counts used by
// StatCard rows on dashboard and turnaround detail.
export function useTurnaroundProgress(detail: TurnaroundDetail | null) {
  return useMemo(() => {
    if (!detail) {
      return { total: 0, completed: 0, accepted: 0, percent: 0, pendingBookings: 0, open: 0 };
    }
    const total = detail.tasks.length;
    const completed = detail.tasks.filter((t) => t.status === "COMPLETED").length;
    const accepted = detail.tasks.filter((t) => t.status === "ACCEPTED").length;
    const open = detail.tasks.filter((t) => t.status !== "COMPLETED").length;
    const percent = total ? Math.round(((completed + accepted * 0.5) / total) * 100) : 0;
    return { total, completed, accepted, percent, pendingBookings: detail.pending_bookings, open };
  }, [detail]);
}

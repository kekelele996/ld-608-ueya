import React from "react";
import { Timeline, Tag } from "antd";
import type { TurnaroundDetail } from "../../types/entities";
import { StatusBadge } from "./StatusBadge";
import { ConflictBadge } from "./ConflictBadge";
import { formatTime } from "../../utils/formatters";

interface TurnaroundTimelineProps {
  detail: TurnaroundDetail;
}

// TurnaroundTimeline renders ordered tasks with planned time, status and
// booking conflict badges; shared by turnarounds detail and tasks page.
export const TurnaroundTimeline: React.FC<TurnaroundTimelineProps> = ({ detail }) => {
  const bookingByTask = new Map(detail.bookings.map((b) => [b.task_id, b]));
  return (
    <Timeline
      items={[...detail.tasks]
        .sort((a, b) => a.planned_start.localeCompare(b.planned_start))
        .map((task) => {
          const booking = bookingByTask.get(task.id);
          return {
            color:
              task.status === "COMPLETED"
                ? "green"
                : task.status === "BLOCKED" || booking?.booking_status === "PENDING"
                  ? "red"
                  : task.status === "ACCEPTED"
                    ? "blue"
                    : "gray",
            children: (
              <div>
                <strong>{task.task_type_text ?? task.task_type}</strong>{" "}
                <StatusBadge kind="task" value={task.status} />
                {task.signed_at && (
                  <Tag color="blue" title="已签收任务在延误重排时保持原时点">
                    已签收·时点冻结
                  </Tag>
                )}
                {booking && <StatusBadge kind="booking" value={booking.booking_status} />}
                {booking?.conflict_code && (
                  <ConflictBadge code={booking.conflict_code} message={booking.conflict_reason} />
                )}
                <div style={{ color: "rgba(0,0,0,0.55)" }}>
                  计划 {formatTime(task.planned_start)} ~ {formatTime(task.deadline)} · 班组 {task.team_id}
                  {task.blocker_note ? ` · 阻塞：${task.blocker_note}` : ""}
                </div>
              </div>
            ),
          };
        })}
    />
  );
};

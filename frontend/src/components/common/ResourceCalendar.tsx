import React, { useMemo } from "react";
import { Table } from "antd";
import type { GroundResource, ResourceBooking } from "../../types/entities";
import { StatusBadge } from "./StatusBadge";
import { ConflictBadge } from "./ConflictBadge";
import { formatTime, formatDateTime } from "../../utils/formatters";

interface ResourceCalendarProps {
  resources: GroundResource[];
  bookings: ResourceBooking[];
  onSelectBooking?: (booking: ResourceBooking) => void;
}

// ResourceCalendar shows each resource with its time-window bookings as
// timeline-like rows. Shared by resources page and turnaround detail.
export const ResourceCalendar: React.FC<ResourceCalendarProps> = ({
  resources,
  bookings,
  onSelectBooking,
}) => {
  const byResource = useMemo(() => {
    const map = new Map<number, ResourceBooking[]>();
    [...bookings]
      .sort((a, b) => a.start_time.localeCompare(b.start_time))
      .forEach((b) => {
        const list = map.get(b.resource_id) ?? [];
        list.push(b);
        map.set(b.resource_id, list);
      });
    return map;
  }, [bookings]);

  return (
    <Table
      rowKey="id"
      size="small"
      pagination={false}
      dataSource={resources}
      columns={[
        { title: "资源", dataIndex: "resource_code", width: 120 },
        {
          title: "类型 / 位置",
          width: 200,
          render: (_, r) => (
            <div>
              <div>{r.resource_type}</div>
              <small style={{ color: "rgba(0,0,0,0.45)" }}>{r.location ?? "—"}</small>
            </div>
          ),
        },
        {
          title: "状态",
          width: 110,
          render: (_, r) => <StatusBadge kind="resource" value={r.availability_status} />,
        },
        {
          title: "预约时段",
          render: (_, r) => {
            const list = byResource.get(r.id) ?? [];
            if (!list.length) return <span style={{ color: "rgba(0,0,0,0.35)" }}>无预约</span>;
            return (
              <div style={{ display: "flex", flexDirection: "column", gap: 4 }}>
                {list.map((b) => (
                  <a
                    key={b.id}
                    onClick={() => onSelectBooking?.(b)}
                    style={{
                      display: "flex",
                      gap: 8,
                      alignItems: "center",
                      padding: "2px 8px",
                      borderRadius: 4,
                      background:
                        b.booking_status === "PENDING"
                          ? "#fffbe6"
                          : b.booking_status === "RELEASED"
                            ? "rgba(0,0,0,0.03)"
                            : "#f6ffed",
                      border:
                        b.booking_status === "PENDING"
                          ? "1px solid #ffe58f"
                          : "1px solid transparent",
                    }}
                  >
                    <span style={{ whiteSpace: "nowrap" }}>
                      {formatDateTime(b.start_time)} ~ {formatTime(b.end_time)}
                    </span>
                    <StatusBadge kind="booking" value={b.booking_status} />
                    <ConflictBadge code={b.conflict_code} message={b.conflict_reason} />
                    <span style={{ marginLeft: "auto", color: "rgba(0,0,0,0.45)" }}>
                      航班 #{b.turnaround_id}
                    </span>
                  </a>
                ))}
              </div>
            );
          },
        },
      ]}
    />
  );
};

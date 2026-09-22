import { Table, Tag } from "antd";
import type { ColumnsType } from "antd/es/table";
import dayjs from "dayjs";
import type { GroundResource } from "../../types/GroundResource";
import type { ResourceBooking } from "../../types/ResourceBooking";
import { formatClock, formatDateTime } from "../../utils/formatters";
import { RESOURCE_STATUS_TEXT } from "../../constants/ResourceStatus";
import { BOOKING_STATUS_TEXT } from "../../constants/BookingStatus";

// ResourceCalendar：资源调度页与过站详情共用的按资源分组预约视图。
export function ResourceCalendar({
  resources,
  bookings,
  date,
  onSelectBooking
}: {
  resources: GroundResource[];
  bookings: ResourceBooking[];
  date?: dayjs.Dayjs;
  onSelectBooking?: (booking: ResourceBooking) => void;
}) {
  const day = date ?? dayjs();
  const dayBookings = bookings.filter((b) => dayjs(b.start_time).isSame(day, "day"));

  const columns: ColumnsType<GroundResource> = [
    {
      title: "资源",
      dataIndex: "resource_code",
      width: 130,
      render: (code: string, row) => (
        <div>
          <strong>{code}</strong>
          <div className="muted">{row.location}</div>
        </div>
      )
    },
    {
      title: "状态",
      dataIndex: "availability_status",
      width: 100,
      render: (status: string) => (
        <Tag color={status === "AVAILABLE" ? "green" : status === "OFFLINE" ? "red" : "orange"}>
          {RESOURCE_STATUS_TEXT[status as keyof typeof RESOURCE_STATUS_TEXT] ?? status}
        </Tag>
      )
    },
    {
      title: `${day.format("MM-DD")} 当日预约`,
      render: (_, row) => {
        const mine = dayBookings
          .filter((b) => b.resource_id === row.id)
          .sort((a, b) => a.start_time.localeCompare(b.start_time));
        if (!mine.length) return <span className="muted">空闲</span>;
        return (
          <div className="booking-strip">
            {mine.map((b) => (
              <Tag
                key={b.id}
                className={`booking-chip booking-${b.booking_status.toLowerCase()}`}
                onClick={() => onSelectBooking?.(b)}
                title={b.conflict_reason || undefined}
              >
                {formatClock(b.start_time)}~{formatClock(b.end_time)}
                {" · "}
                {BOOKING_STATUS_TEXT[b.booking_status]}
              </Tag>
            ))}
          </div>
        );
      }
    },
    {
      title: "维护到期",
      dataIndex: "maintenance_due_at",
      width: 150,
      render: (value: string | null) => formatDateTime(value)
    }
  ];

  return (
    <Table<GroundResource>
      rowKey="id"
      size="small"
      pagination={false}
      columns={columns}
      dataSource={resources}
    />
  );
}

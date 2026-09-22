import React, { useState } from "react";
import { Table, Button, Modal, DatePicker, Space, App } from "antd";
import dayjs, { Dayjs } from "dayjs";
import type { ResourceBooking, ConflictItem } from "../../types/entities";
import { StatusBadge } from "./StatusBadge";
import { ConflictBadge } from "./ConflictBadge";
import { ConflictPanel } from "./ConflictPanel";
import { adjustBooking } from "../../api/ResourceBooking";
import { RequestError } from "../../api/client";
import { formatDateTime, toAPI } from "../../utils/formatters";
import { useRole } from "../../hooks/useRole";

interface PendingBookingListProps {
  rows: ResourceBooking[];
  flightByTurn?: Map<number, string>;
  onChanged: () => void;
}

// PendingBookingList renders bookings squeezed into PENDING by delay
// reschedule, with an adjust-window modal re-running conflict checks.
export const PendingBookingList: React.FC<PendingBookingListProps> = ({
  rows,
  flightByTurn,
  onChanged,
}) => {
  const { message } = App.useApp();
  const { isResource } = useRole();
  const [target, setTarget] = useState<ResourceBooking | null>(null);
  const [start, setStart] = useState<Dayjs | null>(null);
  const [end, setEnd] = useState<Dayjs | null>(null);
  const [conflicts, setConflicts] = useState<ConflictItem[]>([]);
  const [submitting, setSubmitting] = useState(false);

  const openAdjust = (booking: ResourceBooking) => {
    setTarget(booking);
    setStart(dayjs(booking.start_time));
    setEnd(dayjs(booking.end_time));
    setConflicts([]);
  };

  const submit = async () => {
    if (!target || !start || !end) return;
    setSubmitting(true);
    setConflicts([]);
    try {
      await adjustBooking(target.id, toAPI(start.format("YYYY-MM-DDTHH:mm")), toAPI(end.format("YYYY-MM-DDTHH:mm")));
      message.success(`预约 ${target.id} 已调整并确认`);
      setTarget(null);
      onChanged();
    } catch (err) {
      if (err instanceof RequestError) {
        message.error(err.message);
        setConflicts((err.details as ConflictItem[]) ?? []);
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Table
        rowKey="id"
        size="small"
        pagination={false}
        dataSource={rows}
        locale={{ emptyText: "没有待处理预约" }}
        columns={[
          { title: "预约", dataIndex: "id", width: 70, render: (v) => `#${v}` },
          {
            title: "资源",
            render: (_, b) => b.resource_code ?? `#${b.resource_id}`,
          },
          {
            title: "航班",
            render: (_, b) => flightByTurn?.get(b.turnaround_id) ?? `#${b.turnaround_id}`,
          },
          {
            title: "当前时段",
            render: (_, b) => `${formatDateTime(b.start_time)} ~ ${dayjs(b.end_time).format("HH:mm")}`,
          },
          {
            title: "状态 / 原因",
            render: (_, b) => (
              <Space>
                <StatusBadge kind="booking" value={b.booking_status} />
                <ConflictBadge code={b.conflict_code} message={b.conflict_reason} />
              </Space>
            ),
          },
          {
            title: "操作",
            width: 120,
            render: (_, b) =>
              isResource ? (
                <Button type="link" size="small" onClick={() => openAdjust(b)}>
                  调整时段
                </Button>
              ) : (
                <span style={{ color: "rgba(0,0,0,0.35)" }}>资源管理员可调整</span>
              ),
          },
        ]}
      />
      <Modal
        open={!!target}
        title={target ? `调整预约 #${target.id}（${target.resource_code ?? target.resource_id}）` : ""}
        onCancel={() => setTarget(null)}
        onOk={submit}
        confirmLoading={submitting}
        okText="保存并校验"
        cancelText="取消"
        destroyOnClose
      >
        <Space direction="vertical" style={{ width: "100%" }}>
          <span>
            原时段：{target ? formatDateTime(target.start_time) : ""} ~{" "}
            {target ? dayjs(target.end_time).format("HH:mm") : ""}
          </span>
          <Space>
            <span>新开始</span>
            <DatePicker
              showTime={{ format: "HH:mm" }}
              format="MM-DD HH:mm"
              value={start}
              onChange={setStart}
            />
            <span>新结束</span>
            <DatePicker
              showTime={{ format: "HH:mm" }}
              format="MM-DD HH:mm"
              value={end}
              onChange={setEnd}
            />
          </Space>
          <ConflictPanel conflicts={conflicts} title="调整失败：预约保持待处理" />
        </Space>
      </Modal>
    </>
  );
};

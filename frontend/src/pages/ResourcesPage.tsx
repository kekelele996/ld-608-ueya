import {
  Alert,
  Button,
  Card,
  DatePicker,
  Form,
  Modal,
  Select,
  Space,
  Table,
  Tabs,
  Tag,
  message
} from "antd";
import type { ColumnsType } from "antd/es/table";
import dayjs from "dayjs";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useAppDispatch, useAppSelector } from "../stores/store";
import { fetchResources } from "../stores/resourceSlice";
import { fetchBookings, fetchPendingBookings } from "../stores/bookingSlice";
import { fetchTurnarounds } from "../stores/turnaroundSlice";
import { groundResourceApi } from "../api/GroundResource";
import { resourceBookingApi } from "../api/ResourceBooking";
import { ResourceCalendar } from "../components/common/ResourceCalendar";
import { StatusBadge } from "../components/common/StatusBadge";
import { ConflictBadge } from "../components/common/ConflictBadge";
import { AdjustBookingModal } from "../components/common/AdjustBookingModal";
import { useResourceConflict } from "../hooks/useResourceConflict";
import { formatDateTime, toApiTime } from "../utils/formatters";
import {
  RESOURCE_STATUSES,
  RESOURCE_STATUS_TEXT,
  type ResourceStatusValue
} from "../constants/ResourceStatus";
import { BOOKING_STATUS_TEXT } from "../constants/BookingStatus";
import { GROUND_TASK_TYPE_TEXT } from "../constants/GroundTaskType";
import { canRole } from "../constants/UserRole";
import { teamText } from "../constants/teams";
import type { GroundResource } from "../types/GroundResource";
import type { ResourceBooking } from "../types/ResourceBooking";
import { renderLog, LOG_TEMPLATES } from "../constants/logTemplates";

export function ResourcesPage() {
  const dispatch = useAppDispatch();
  const resources = useAppSelector((s) => s.resources.rows);
  const bookings = useAppSelector((s) => s.bookings.rows);
  const pending = useAppSelector((s) => s.bookings.pending);
  const turnarounds = useAppSelector((s) => s.turnarounds.rows);
  const user = useAppSelector((s) => s.auth.user);

  const [day, setDay] = useState(dayjs());
  const [adjustTarget, setAdjustTarget] = useState<ResourceBooking | null>(null);
  const [maintenanceTarget, setMaintenanceTarget] = useState<GroundResource | null>(null);
  const [maintenanceForm] = Form.useForm();

  const refresh = useCallback(() => {
    dispatch(fetchResources());
    dispatch(fetchBookings());
    dispatch(fetchPendingBookings());
    dispatch(fetchTurnarounds());
  }, [dispatch]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const conflict = useResourceConflict(bookings);
  const canMaintain = canRole(user?.role, "RESOURCE_SET_MAINTENANCE");
  const canAdjust = canRole(user?.role, "BOOKING_ADJUST");

  const onAdjustSubmit = async (values: {
    resource_id: number;
    start_time: string;
    end_time: string;
  }) => {
    if (!adjustTarget) return;
    const updated = await resourceBookingApi.adjust(adjustTarget.id, values);
    if (updated.booking_status === "CONFIRMED") {
      message.success(
        renderLog(LOG_TEMPLATES.ResourceBooking.confirm, { bookingId: updated.id })
      );
    } else {
      message.warning(`仍存在冲突，保持待处理：${updated.conflict_reason}`);
    }
    setAdjustTarget(null);
    refresh();
  };

  const resourceById = useMemo(
    () => new Map(resources.map((r) => [r.id, r])),
    [resources]
  );
  const turnMap = useMemo(() => new Map(turnarounds.map((t) => [t.id, t])), [turnarounds]);

  const resourceColumns: ColumnsType<GroundResource> = [
    { title: "资源编号", dataIndex: "resource_code", width: 120 },
    {
      title: "类型",
      dataIndex: "resource_type",
      width: 110,
      render: (v: string) => GROUND_TASK_TYPE_TEXT[v as keyof typeof GROUND_TASK_TYPE_TEXT] ?? v
    },
    { title: "位置", dataIndex: "location" },
    { title: "归属班组", dataIndex: "owner_team", width: 130, render: teamText },
    {
      title: "状态",
      dataIndex: "availability_status",
      width: 100,
      render: (v: string) => <StatusBadge value={v} />
    },
    { title: "维护到期", dataIndex: "maintenance_due_at", width: 150, render: formatDateTime },
    {
      title: "操作",
      width: 130,
      render: (_, r) =>
        canMaintain ? (
          <Button
            size="small"
            onClick={() => {
              maintenanceForm.setFieldsValue({
                status: r.availability_status,
                due: r.maintenance_due_at ? dayjs(r.maintenance_due_at) : null
              });
              setMaintenanceTarget(r);
            }}
          >
            维护/离线
          </Button>
        ) : (
          <span className="muted">仅资源管理员</span>
        )
    }
  ];

  const pendingColumns: ColumnsType<ResourceBooking> = [
    {
      title: "资源",
      render: (_, b) => resourceById.get(b.resource_id)?.resource_code ?? `#${b.resource_id}`
    },
    {
      title: "航班",
      render: (_, b) => turnMap.get(b.turnaround_id)?.flight_no ?? `#${b.turnaround_id}`
    },
    { title: "窗口", render: (_, b) => `${formatDateTime(b.start_time)} ~ ${formatDateTime(b.end_time)}` },
    {
      title: "状态",
      dataIndex: "booking_status",
      render: (v: string) => <Tag color="orange">{BOOKING_STATUS_TEXT[v as keyof typeof BOOKING_STATUS_TEXT] ?? v}</Tag>
    },
    {
      title: "冲突原因",
      dataIndex: "conflict_reason",
      render: (v: string) => <ConflictBadge message={v} />
    },
    {
      title: "操作",
      width: 150,
      render: (_, b) => (
        <Space>
          {canAdjust && (
            <Button size="small" type="primary" ghost onClick={() => setAdjustTarget(b)}>
              调整时段
            </Button>
          )}
          {canAdjust && (
            <Button size="small" danger ghost onClick={async () => {
              await resourceBookingApi.release(b.id);
              message.success("预约已释放");
              refresh();
            }}>
              释放
            </Button>
          )}
        </Space>
      )
    }
  ];

  return (
    <div className="page-grid">
      {conflict.hasConflict && (
        <Alert
          type="warning"
          showIcon
          message={`${conflict.count} 个预约处于待处理状态（被延误挤出或改期冲突），请调整时段或释放资源。`}
        />
      )}

      <Tabs
        defaultActiveKey="calendar"
        items={[
          {
            key: "calendar",
            label: "预约日历",
            children: (
              <Card
                title="资源预约日历"
                extra={
                  <Space>
                    <span className="muted">日期</span>
                    <DatePicker value={day} onChange={(d) => d && setDay(d)} allowClear={false} />
                  </Space>
                }
              >
                <ResourceCalendar
                  resources={resources}
                  bookings={bookings}
                  date={day}
                  onSelectBooking={(b) => canAdjust && setAdjustTarget(b)}
                />
              </Card>
            )
          },
          {
            key: "ledger",
            label: `资源台账 (${resources.length})`,
            children: (
              <Card>
                <Table<GroundResource>
                  rowKey="id"
                  size="small"
                  columns={resourceColumns}
                  dataSource={resources}
                  pagination={{ pageSize: 8 }}
                />
              </Card>
            )
          },
          {
            key: "pending",
            label: `待处理预约 (${pending.length})`,
            children: (
              <Card>
                <Table<ResourceBooking>
                  rowKey="id"
                  size="small"
                  columns={pendingColumns}
                  dataSource={pending}
                  pagination={{ pageSize: 8 }}
                  locale={{ emptyText: "没有待处理预约" }}
                />
              </Card>
            )
          }
        ]}
      />

      <AdjustBookingModal
        open={!!adjustTarget}
        booking={adjustTarget}
        resources={resources}
        onClose={() => setAdjustTarget(null)}
        onSubmit={onAdjustSubmit}
      />

      <Modal
        title={`资源状态维护：${maintenanceTarget?.resource_code ?? ""}`}
        open={!!maintenanceTarget}
        onCancel={() => setMaintenanceTarget(null)}
        onOk={async () => {
          const values = await maintenanceForm.validateFields();
          if (maintenanceTarget) {
            await groundResourceApi.setStatus(maintenanceTarget.id, {
              status: values.status,
              maintenance_due_at: values.due ? toApiTime(values.due) : null
            });
            message.success(
              renderLog(LOG_TEMPLATES.GroundResource.maintenance, {
                resource: maintenanceTarget.resource_code,
                status: RESOURCE_STATUS_TEXT[values.status as ResourceStatusValue] ?? values.status
              })
            );
            setMaintenanceTarget(null);
            refresh();
          }
        }}
        destroyOnClose
      >
        <Form form={maintenanceForm} layout="vertical">
          <Form.Item name="status" label="可用状态" rules={[{ required: true }]}>
            <Select
              options={RESOURCE_STATUSES.map((s) => ({ value: s, label: RESOURCE_STATUS_TEXT[s] }))}
            />
          </Form.Item>
          <Form.Item name="due" label="维护到期时间">
            <DatePicker showTime format="MM-DD HH:mm" style={{ width: "100%" }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

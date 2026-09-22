import { Button, Card, Form, Input, Modal, Select, Space, Table, Tag, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../stores/store";
import { fetchTasks } from "../stores/taskSlice";
import { fetchTurnarounds } from "../stores/turnaroundSlice";
import { fetchBookings } from "../stores/bookingSlice";
import { groundTaskApi } from "../api/GroundTask";
import { StatusBadge } from "../components/common/StatusBadge";
import { ConflictBadge } from "../components/common/ConflictBadge";
import { TeamTag } from "../components/common/TeamTag";
import { usePagination } from "../hooks/usePagination";
import { formatClock } from "../utils/formatters";
import { GROUND_TASK_TYPE_TEXT } from "../constants/GroundTaskType";
import {
  GROUND_TASK_STATUSES,
  GROUND_TASK_STATUS_TEXT,
  type GroundTaskStatusValue
} from "../constants/GroundTaskStatus";
import { canRole } from "../constants/UserRole";
import type { GroundTask } from "../types/GroundTask";

export function TasksPage() {
  const dispatch = useAppDispatch();
  const tasks = useAppSelector((s) => s.tasks.rows);
  const turnarounds = useAppSelector((s) => s.turnarounds.rows);
  const bookings = useAppSelector((s) => s.bookings.rows);
  const user = useAppSelector((s) => s.auth.user);
  const [typeFilter, setTypeFilter] = useState<string>("ALL");
  const [statusFilter, setStatusFilter] = useState<GroundTaskStatusValue | "ALL">("ALL");
  const [blockTarget, setBlockTarget] = useState<GroundTask | null>(null);
  const [blockForm] = Form.useForm();

  useEffect(() => {
    dispatch(fetchTasks());
    dispatch(fetchTurnarounds());
    dispatch(fetchBookings());
  }, [dispatch]);

  const filtered = useMemo(
    () =>
      tasks.filter(
        (t) =>
          (typeFilter === "ALL" || t.task_type === typeFilter) &&
          (statusFilter === "ALL" || t.status === statusFilter)
      ),
    [tasks, typeFilter, statusFilter]
  );
  const { page, setPage, pageRows, pageSize, total } = usePagination(filtered, 8);

  const turnMap = useMemo(() => new Map(turnarounds.map((t) => [t.id, t])), [turnarounds]);
  const canSign = canRole(user?.role, "TASK_SIGN");

  const refresh = () => dispatch(fetchTasks());

  const columns: ColumnsType<GroundTask> = [
    {
      title: "任务",
      dataIndex: "task_type",
      width: 110,
      render: (v: string) => GROUND_TASK_TYPE_TEXT[v as keyof typeof GROUND_TASK_TYPE_TEXT] ?? v
    },
    {
      title: "航班",
      render: (_, r) => {
        const turn = turnMap.get(r.turnaround_id);
        return (
          <Link to={`/turnarounds/${r.turnaround_id}`}>{turn?.flight_no ?? `#${r.turnaround_id}`}</Link>
        );
      }
    },
    { title: "班组", dataIndex: "team_id", width: 130, render: (v: string) => <TeamTag teamId={v} /> },
    { title: "计划", width: 130, render: (_, r) => `${formatClock(r.planned_start)}~${formatClock(r.planned_end)}` },
    { title: "截止", dataIndex: "deadline", width: 80, render: formatClock },
    {
      title: "状态",
      dataIndex: "status",
      width: 90,
      render: (v: string) => <StatusBadge value={v} />,
      filters: GROUND_TASK_STATUSES.map((s) => ({ text: GROUND_TASK_STATUS_TEXT[s], value: s }))
    },
    {
      title: "签收/阻塞",
      render: (_, r) => (
        <Space direction="vertical" size={0}>
          {r.signed_by && <span className="muted">签收：{r.signed_by}</span>}
          {r.blocker_note && <ConflictBadge code="BLOCKED" message={r.blocker_note} />}
          {bookings.find((b) => b.task_id === r.id && b.booking_status === "PENDING")?.conflict_reason && (
            <Tag color="orange">关联预约待处理</Tag>
          )}
        </Space>
      )
    },
    {
      title: "操作",
      width: 170,
      render: (_, r) =>
        canSign ? (
          <Space size={2}>
            {r.status === "PLANNED" && (
              <Button size="small" type="primary" ghost onClick={async () => {
                await groundTaskApi.sign(r.id, user?.name);
                message.success("签收成功，延误重排时将保持此时点");
                refresh();
              }}>
                签收
              </Button>
            )}
            {(r.status === "SIGNED" || r.status === "BLOCKED") && (
              <Button size="small" type="primary" onClick={async () => {
                await groundTaskApi.finish(r.id);
                message.success("任务已完成");
                refresh();
              }}>
                完成
              </Button>
            )}
            {r.status !== "FINISHED" && (
              <Button size="small" danger ghost onClick={() => setBlockTarget(r)}>
                阻塞
              </Button>
            )}
          </Space>
        ) : (
          <span className="muted">仅班组/调度可操作</span>
        )
    }
  ];

  return (
    <Card
      title="地勤任务"
      extra={
        <Space>
          <Select
            value={typeFilter}
            style={{ width: 130 }}
            onChange={(v) => { setTypeFilter(v); setPage(1); }}
            options={[
              { value: "ALL", label: "全部类型" },
              ...Object.entries(GROUND_TASK_TYPE_TEXT).map(([value, label]) => ({ value, label }))
            ]}
          />
          <Select
            value={statusFilter}
            style={{ width: 120 }}
            onChange={(v) => { setStatusFilter(v); setPage(1); }}
            options={[
              { value: "ALL", label: "全部状态" },
              ...GROUND_TASK_STATUSES.map((s) => ({ value: s, label: GROUND_TASK_STATUS_TEXT[s] }))
            ]}
          />
        </Space>
      }
    >
      <Table<GroundTask>
        rowKey="id"
        size="small"
        columns={columns}
        dataSource={pageRows}
        pagination={{ current: page, pageSize, total, onChange: setPage }}
      />
      <Modal
        title="登记阻塞原因"
        open={!!blockTarget}
        onCancel={() => setBlockTarget(null)}
        onOk={async () => {
          const values = await blockForm.validateFields();
          if (blockTarget) {
            await groundTaskApi.block(blockTarget.id, values.note);
            message.success("已登记阻塞");
            setBlockTarget(null);
            blockForm.resetFields();
            refresh();
          }
        }}
        destroyOnClose
      >
        <Form form={blockForm} layout="vertical">
          <Form.Item name="note" label="阻塞原因" rules={[{ required: true, message: "请填写阻塞原因" }]}>
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}

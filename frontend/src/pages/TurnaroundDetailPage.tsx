import {
  Alert,
  Button,
  Card,
  Col,
  Descriptions,
  Form,
  Input,
  InputNumber,
  Modal,
  Row,
  Select,
  Space,
  Table,
  Tag,
  Timeline,
  message
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../stores/store";
import { fetchTurnaroundDetail } from "../stores/turnaroundSlice";
import { fetchResources } from "../stores/resourceSlice";
import { fetchTasks } from "../stores/taskSlice";
import { turnaroundApi } from "../api/FlightTurnaround";
import { groundTaskApi } from "../api/GroundTask";
import { resourceBookingApi } from "../api/ResourceBooking";
import { ApiError } from "../api/client";
import { StatusBadge } from "../components/common/StatusBadge";
import { ConflictPanel, ConflictBadge } from "../components/common/ConflictBadge";
import { TurnaroundTimeline } from "../components/common/TurnaroundTimeline";
import { AdjustBookingModal } from "../components/common/AdjustBookingModal";
import { useResourceConflict } from "../hooks/useResourceConflict";
import { formatClock, formatDateTime, formatFullDateTime } from "../utils/formatters";
import { GROUND_TASK_TYPE_TEXT } from "../constants/GroundTaskType";
import { GROUND_TASK_STATUS_TEXT } from "../constants/GroundTaskStatus";
import { DELAY_TYPES, DELAY_TYPE_TEXT } from "../constants/DelayType";
import { TEAMS, teamText, taskTypeText } from "../constants/teams";
import { canRole } from "../constants/UserRole";
import { renderLog, LOG_TEMPLATES } from "../constants/logTemplates";
import type { ConflictItem, PlanConflictResult, ReplanResult } from "../types/api";
import type { GroundTask } from "../types/GroundTask";
import type { ResourceBooking } from "../types/ResourceBooking";

export function TurnaroundDetailPage() {
  const { id } = useParams();
  const turnaroundId = Number(id);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const detail = useAppSelector((s) => s.turnarounds.detail);
  const resources = useAppSelector((s) => s.resources.rows);
  const user = useAppSelector((s) => s.auth.user);

  const [conflicts, setConflicts] = useState<ConflictItem[]>([]);
  const [replan, setReplan] = useState<ReplanResult | null>(null);
  const [delayOpen, setDelayOpen] = useState(false);
  const [blockTarget, setBlockTarget] = useState<GroundTask | null>(null);
  const [adjustTarget, setAdjustTarget] = useState<ResourceBooking | null>(null);
  const [busy, setBusy] = useState(false);
  const [delayForm] = Form.useForm();
  const [blockForm] = Form.useForm();

  const refresh = useCallback(() => {
    dispatch(fetchTurnaroundDetail(turnaroundId));
    dispatch(fetchResources());
    dispatch(fetchTasks());
  }, [dispatch, turnaroundId]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const tasks = detail?.tasks ?? [];
  const bookings = detail?.bookings ?? [];
  const delays = detail?.delays ?? [];
  const conflict = useResourceConflict(bookings);
  const pendingTaskIds = useMemo(
    () => conflict.pending.map((b) => b.task_id).filter((v): v is number => !!v),
    [conflict.pending]
  );

  const resourceById = useMemo(
    () => new Map(resources.map((r) => [r.id, r])),
    [resources]
  );

  const canGenerate = canRole(user?.role, "PLAN_GENERATE");
  const canSign = canRole(user?.role, "TASK_SIGN");
  const canAdjust = canRole(user?.role, "BOOKING_ADJUST");
  const canDelay = canRole(user?.role, "DELAY_REGISTER");
  const canRelease = canRole(user?.role, "TURNAROUND_RELEASE");

  // 1) 到站后按任务类型/班组/资源时段整批生成；冲突整批不保存并逐项展示。
  const onGenerate = async () => {
    setBusy(true);
    setConflicts([]);
    try {
      const result: PlanConflictResult = await turnaroundApi.generatePlan(turnaroundId);
      if (result.saved) {
        message.success(
          renderLog(LOG_TEMPLATES.GroundTask.dispatch, {
            flightNo: detail?.flight_no ?? "",
            count: result.task_count
          })
        );
      }
      refresh();
    } catch (err) {
      if (err instanceof ApiError && err.code === "PLAN_CONFLICT") {
        const result = err.data as PlanConflictResult;
        setConflicts(result?.items ?? []);
        message.error(`整批未保存：共 ${result?.items.length ?? 0} 项冲突`);
      }
    } finally {
      setBusy(false);
    }
  };

  const onSign = async (task: GroundTask) => {
    await groundTaskApi.sign(task.id, user?.name);
    message.success(`已签收任务：${taskTypeText(task.task_type)}，计划时点保持 ${formatClock(task.planned_start)}`);
    refresh();
  };

  const onFinish = async (task: GroundTask) => {
    await groundTaskApi.finish(task.id);
    message.success(`任务已完成：${taskTypeText(task.task_type)}`);
    refresh();
  };

  // 2) 延误登记增加分钟数：未完成任务重排、已签收保持原时点、挤出预约待处理。
  const onRegisterDelay = async () => {
    const values = await delayForm.validateFields();
    setBusy(true);
    try {
      const result = await turnaroundApi.registerDelay(turnaroundId, {
        delay_type: values.delay_type,
        minutes: values.minutes,
        root_cause: values.root_cause,
        responsibility_team: values.responsibility_team
      });
      setReplan(result);
      message.success(
        renderLog(LOG_TEMPLATES.DelayEvent.replan, {
          flightNo: detail?.flight_no ?? "",
          moved: result.moved_tasks.length,
          bumped: result.bumped.length
        })
      );
      setDelayOpen(false);
      refresh();
    } finally {
      setBusy(false);
    }
  };

  // 3) 页面调整预约时段；保存后刷新看到同一结果（CONFIRMED 或 PENDING+原因）。
  const onAdjustSubmit = async (values: {
    resource_id: number;
    start_time: string;
    end_time: string;
  }) => {
    if (!adjustTarget) return;
    const updated = await resourceBookingApi.adjust(adjustTarget.id, values);
    if (updated.booking_status === "CONFIRMED") {
      message.success(`预约 #${updated.id} 已确认新时段`);
    } else {
      message.warning(`预约 #${updated.id} 仍为待处理：${updated.conflict_reason}`);
    }
    setAdjustTarget(null);
    refresh();
  };

  const onRelease = async () => {
    await turnaroundApi.release(turnaroundId);
    message.success("航班已放行");
    navigate("/turnarounds");
  };

  const taskColumns: ColumnsType<GroundTask> = [
    {
      title: "任务类型",
      dataIndex: "task_type",
      width: 120,
      render: (v: string) => GROUND_TASK_TYPE_TEXT[v as keyof typeof GROUND_TASK_TYPE_TEXT] ?? v
    },
    { title: "班组", dataIndex: "team_id", width: 130, render: teamText },
    { title: "计划开始", dataIndex: "planned_start", width: 90, render: formatClock },
    { title: "计划结束", dataIndex: "planned_end", width: 90, render: formatClock },
    { title: "截止", dataIndex: "deadline", width: 90, render: formatClock },
    {
      title: "状态",
      dataIndex: "status",
      width: 90,
      render: (v: string) => <StatusBadge value={v} />
    },
    {
      title: "资源预约 / 冲突",
      render: (_, task) => {
        const mine = bookings.filter((b) => b.task_id === task.id);
        if (!mine.length) return <span className="muted">—</span>;
        return (
          <Space direction="vertical" size={2}>
            {mine.map((b) => (
              <Space key={b.id} size={4} wrap>
                <Tag>{resourceById.get(b.resource_id)?.resource_code ?? `R${b.resource_id}`}</Tag>
                <span className="muted">
                  {formatClock(b.start_time)}~{formatClock(b.end_time)}
                </span>
                <StatusBadge value={b.booking_status} />
                {b.conflict_reason && <ConflictBadge code="RESOURCE_TIME_OVERLAP" message={b.conflict_reason} />}
                {b.booking_status === "PENDING" && canAdjust && (
                  <Button size="small" type="link" onClick={() => setAdjustTarget(b)}>
                    调整时段
                  </Button>
                )}
              </Space>
            ))}
          </Space>
        );
      }
    },
    {
      title: "操作",
      width: 200,
      render: (_, task) => (
        <Space size={2}>
          {task.status === "PLANNED" && canSign && (
            <Button size="small" type="primary" ghost onClick={() => onSign(task)}>
              签收
            </Button>
          )}
          {(task.status === "SIGNED" || task.status === "BLOCKED") && canSign && (
            <Button size="small" type="primary" onClick={() => onFinish(task)}>
              完成
            </Button>
          )}
          {task.status !== "FINISHED" && canSign && (
            <Button size="small" danger ghost onClick={() => setBlockTarget(task)}>
              阻塞
            </Button>
          )}
        </Space>
      )
    }
  ];

  const bookingColumns: ColumnsType<ResourceBooking> = [
    {
      title: "资源",
      render: (_, b) => resourceById.get(b.resource_id)?.resource_code ?? `资源#${b.resource_id}`
    },
    { title: "开始", dataIndex: "start_time", width: 110, render: formatDateTime },
    { title: "结束", dataIndex: "end_time", width: 110, render: formatDateTime },
    {
      title: "状态",
      dataIndex: "booking_status",
      width: 100,
      render: (v: string) => <StatusBadge value={v} />
    },
    {
      title: "冲突原因",
      dataIndex: "conflict_reason",
      render: (v: string) => (v ? <ConflictBadge message={v} /> : <span className="muted">—</span>)
    },
    {
      title: "操作",
      width: 120,
      render: (_, b) =>
        canAdjust ? (
          <Button size="small" type="link" onClick={() => setAdjustTarget(b)}>
            调整时段
          </Button>
        ) : null
    }
  ];

  if (!detail) {
    return <Card loading />;
  }

  return (
    <div className="page-grid detail-page">
      <Card
        title={
          <Space wrap>
            <Button size="small" onClick={() => navigate("/turnarounds")}>
              ← 返回列表
            </Button>
            <strong style={{ fontSize: 16 }}>
              {detail.flight_no} · {detail.aircraft_reg}
            </strong>
            <StatusBadge value={detail.turnaround_status} />
            {detail.accumulated_delay > 0 && <Tag color="volcano">累计延误 {detail.accumulated_delay} 分钟</Tag>}
          </Space>
        }
        extra={
          <Space wrap>
            {!detail.plan_generated && canGenerate && (
              <Button type="primary" loading={busy} onClick={onGenerate}>
                生成保障任务与资源预约
              </Button>
            )}
            {canDelay && (
              <Button danger onClick={() => setDelayOpen(true)}>
                延误登记并重排
              </Button>
            )}
            {canRelease && (
              <Button onClick={onRelease} disabled={tasks.length === 0}>
                放行航班
              </Button>
            )}
          </Space>
        }
      >
        <Descriptions size="small" column={3}>
          <Descriptions.Item label="机位">{detail.stand_no}</Descriptions.Item>
          <Descriptions.Item label="到站时间">{formatFullDateTime(detail.arrival_time)}</Descriptions.Item>
          <Descriptions.Item label="离港时间">{formatFullDateTime(detail.departure_time)}</Descriptions.Item>
          <Descriptions.Item label="延误原因">{detail.delay_reason || "—"}</Descriptions.Item>
          <Descriptions.Item label="计划状态">
            {detail.plan_generated ? "已派发" : "未派发"}
          </Descriptions.Item>
          <Descriptions.Item label="班组签收人">
            {tasks.filter((t) => t.signed_by).map((t) => t.signed_by).filter((v, i, a) => a.indexOf(v) === i).join("、") || "—"}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      {conflicts.length > 0 && <ConflictPanel items={conflicts} />}

      {replan && (
        <Alert
          type="warning"
          showIcon
          className="replan-alert"
          message={`延误 +${replan.shift_minutes} 分钟重排结果（刷新后保持一致）`}
          description={
            <div>
              <Space wrap>
                <Tag color="blue">保持原时点 {replan.kept_tasks.length} 项（已签收）</Tag>
                <Tag color="orange">重排任务 {replan.moved_tasks.length} 项</Tag>
                <Tag color="red">挤出预约 {replan.bumped.length} 项 → 待处理</Tag>
              </Space>
              <ul className="replan-keep-list">
                {replan.kept_tasks.map((k) => (
                  <li key={k.task_id}>
                    {taskTypeText(k.task_type)}（{GROUND_TASK_STATUS_TEXT[k.status as keyof typeof GROUND_TASK_STATUS_TEXT] ?? k.status}）
                    保持 {formatClock(k.planned_start)} —— {k.reason}
                  </li>
                ))}
              </ul>
              {replan.bumped.length > 0 && <ConflictPanel title="被挤出的预约（待处理）" items={replan.bumped} />}
            </div>
          }
        />
      )}

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={14}>
          <Card title="地勤任务（按类型/班组派发）" className="full-height">
            <Table<GroundTask>
              rowKey="id"
              size="small"
              columns={taskColumns}
              dataSource={tasks}
              pagination={false}
              scroll={{ x: 980 }}
            />
          </Card>
        </Col>
        <Col xs={24} lg={10}>
          <Card title="保障时间轴" className="full-height">
            <TurnaroundTimeline tasks={tasks} pendingTaskIds={pendingTaskIds} />
          </Card>
        </Col>
      </Row>

      <Card
        title={
          <Space>
            资源预约
            {conflict.count > 0 && <Tag color="red">{conflict.count} 项待处理</Tag>}
          </Space>
        }
      >
        <Table<ResourceBooking>
          rowKey="id"
          size="small"
          columns={bookingColumns}
          dataSource={bookings}
          pagination={false}
        />
      </Card>

      <Card title="延误事件">
        <Timeline
          items={delays.map((d) => ({
            color: d.resolved_at ? "green" : "red",
            children: (
              <Space wrap>
                <Tag color="volcano">{DELAY_TYPE_TEXT[d.delay_type as keyof typeof DELAY_TYPE_TEXT] ?? d.delay_type}</Tag>
                <strong>+{d.minutes} 分钟</strong>
                <span className="muted">{d.root_cause || "无原因说明"}</span>
                {d.responsibility_team && <span>责任：{teamText(d.responsibility_team)}</span>}
                {d.resolved_at && <Tag color="green">已关闭 {formatDateTime(d.resolved_at)}</Tag>}
              </Space>
            )
          }))}
        />
      </Card>

      {/* 延误登记 */}
      <Modal
        title="延误登记并触发重排"
        open={delayOpen}
        onCancel={() => setDelayOpen(false)}
        onOk={onRegisterDelay}
        confirmLoading={busy}
        destroyOnClose
      >
        <Alert
          type="info"
          showIcon
          className="modal-alert"
          message="登记分钟数后：未签收任务与预约整体顺延，已签收任务保持原时点，撞车预约进入待处理。"
        />
        <Form form={delayForm} layout="vertical" initialValues={{ delay_type: "ATC", minutes: 40 }}>
          <Form.Item name="delay_type" label="延误类型" rules={[{ required: true }]}>
            <Select options={DELAY_TYPES.map((t) => ({ value: t, label: DELAY_TYPE_TEXT[t] }))} />
          </Form.Item>
          <Form.Item name="minutes" label="延误分钟数" rules={[{ required: true }]}>
            <InputNumber min={1} max={720} style={{ width: "100%" }} addonAfter="分钟" />
          </Form.Item>
          <Form.Item name="responsibility_team" label="责任班组">
            <Select allowClear options={TEAMS.map((t) => ({ value: t.id, label: t.name }))} />
          </Form.Item>
          <Form.Item name="root_cause" label="原因说明">
            <Input.TextArea rows={2} placeholder="如：出港流控等待" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 阻塞登记 */}
      <Modal
        title={`阻塞任务：${blockTarget ? taskTypeText(blockTarget.task_type) : ""}`}
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

      <AdjustBookingModal
        open={!!adjustTarget}
        booking={adjustTarget}
        resources={resources}
        requiredType={
          adjustTarget?.task_id
            ? tasks.find((t) => t.id === adjustTarget.task_id)?.task_type
            : undefined
        }
        onClose={() => setAdjustTarget(null)}
        onSubmit={onAdjustSubmit}
      />
    </div>
  );
}

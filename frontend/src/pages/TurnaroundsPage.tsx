import React, { useCallback, useEffect, useMemo, useState } from "react";
import {
  Card,
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  DatePicker,
  Drawer,
  Progress,
  App,
  Row,
  Col,
  Divider,
} from "antd";
import dayjs from "dayjs";
import { useSearchParams } from "react-router-dom";
import type { FlightTurnaround, TurnaroundDetail, PlanResult, ConflictItem } from "../types/entities";
import {
  listTurnarounds,
  getTurnaround,
  registerTurnaround,
  generatePlan,
} from "../api/FlightTurnaround";
import { registerDelay } from "../api/DelayEvent";
import { acceptTask, completeTask, blockTask } from "../api/GroundTask";
import { RequestError } from "../api/client";
import { StatusBadge } from "../components/common/StatusBadge";
import { DelayTag } from "../components/common/DelayTag";
import { ConflictPanel } from "../components/common/ConflictPanel";
import { TurnaroundTimeline } from "../components/common/TurnaroundTimeline";
import { TeamTag } from "../components/common/TeamTag";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { useRole } from "../hooks/useRole";
import { DELAY_TYPES } from "../types/auth";
import { formatDateTime, formatFullDateTime, toAPI } from "../utils/formatters";
import { createTurnaroundForm, createDelayForm } from "../constructors/FlightTurnaroundConstructor";

export const TurnaroundsPage: React.FC = () => {
  const { message } = App.useApp();
  const { isDispatcher, isTeam } = useRole();
  const [searchParams] = useSearchParams();
  const focusId = Number(searchParams.get("focus") ?? 0);

  const [rows, setRows] = useState<FlightTurnaround[]>([]);
  const [loading, setLoading] = useState(false);
  const [registerOpen, setRegisterOpen] = useState(false);
  const [detail, setDetail] = useState<TurnaroundDetail | null>(null);
  const [detailOpen, setDetailOpen] = useState(false);
  const [delayOpen, setDelayOpen] = useState(false);
  const [busy, setBusy] = useState(false);
  const [planConflicts, setPlanConflicts] = useState<ConflictItem[]>([]);
  const [delayConflicts, setDelayConflicts] = useState<ConflictItem[]>([]);
  const [lastResult, setLastResult] = useState<PlanResult | null>(null);
  const [blockTarget, setBlockTarget] = useState<number | null>(null);
  const [blockNote, setBlockNote] = useState("");

  const [registerForm] = Form.useForm();
  const [delayForm] = Form.useForm();

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const res = await listTurnarounds({ page: 1, page_size: 100 });
      setRows(res.items);
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const openDetail = useCallback(
    async (id: number) => {
      try {
        const d = await getTurnaround(id);
        setDetail(d);
        setDetailOpen(true);
        setPlanConflicts([]);
        setDelayConflicts([]);
        setLastResult(null);
      } catch (err) {
        message.error((err as Error).message);
      }
    },
    [message],
  );

  useEffect(() => {
    if (focusId) openDetail(focusId);
  }, [focusId, openDetail]);

  const progress = useTurnaroundProgress(detail);

  const submitRegister = async () => {
    const values = registerForm.getFieldsValue();
    setBusy(true);
    try {
      await registerTurnaround({
        flight_no: values.flight_no,
        aircraft_reg: values.aircraft_reg,
        stand_no: values.stand_no,
        arrival_time: toAPI(values.range[0].format("YYYY-MM-DDTHH:mm")),
        departure_time: toAPI(values.range[1].format("YYYY-MM-DDTHH:mm")),
      });
      message.success("航班到站已登记，可生成保障计划");
      setRegisterOpen(false);
      registerForm.resetFields();
      refresh();
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setBusy(false);
    }
  };

  // Whole-batch generation. On ANY conflict the backend rolls the batch back
  // and returns per-item reasons; the drawer shows them without refreshing.
  const onGeneratePlan = async () => {
    if (!detail) return;
    setBusy(true);
    setPlanConflicts([]);
    try {
      const result = await generatePlan(detail.id);
      setLastResult(result);
      message.success(`已整批保存 ${result.tasks.length} 个任务和 ${result.bookings.length} 条预约`);
      const refreshed = await getTurnaround(detail.id);
      setDetail(refreshed);
      refresh();
    } catch (err) {
      if (err instanceof RequestError) {
        message.error(err.message);
        setPlanConflicts((err.details as ConflictItem[]) ?? []);
        setLastResult((err.result as PlanResult) ?? null);
      }
    } finally {
      setBusy(false);
    }
  };

  const onRegisterDelay = async () => {
    if (!detail) return;
    const values = delayForm.getFieldsValue();
    setBusy(true);
    setDelayConflicts([]);
    try {
      const result = await registerDelay(detail.id, {
        delay_type: values.delay_type,
        minutes: Number(values.minutes),
        root_cause: values.root_cause,
        responsibility_team: values.responsibility_team,
      });
      setLastResult(result);
      setDelayConflicts(result.conflicts ?? []);
      message.success(
        `延误 +${values.minutes} 分钟：重排 ${result.rescheduled_tasks ?? 0} 个任务，` +
          `${result.pending_bookings ?? 0} 条预约待处理，已签收任务保持原时点`,
      );
      const refreshed = await getTurnaround(detail.id);
      setDetail(refreshed);
      setDelayOpen(false);
      refresh();
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setBusy(false);
    }
  };

  const reloadDetail = async () => {
    if (!detail) return;
    const d = await getTurnaround(detail.id);
    setDetail(d);
    refresh();
  };

  const onAccept = async (taskId: number) => {
    try {
      await acceptTask(taskId);
      message.success(`任务 ${taskId} 已签收，计划时点已冻结`);
      await reloadDetail();
    } catch (err) {
      message.error((err as Error).message);
    }
  };

  const onComplete = async (taskId: number) => {
    try {
      await completeTask(taskId);
      message.success(`任务 ${taskId} 已完成，资源预约已释放`);
      await reloadDetail();
    } catch (err) {
      message.error((err as Error).message);
    }
  };

  const onBlock = async () => {
    if (!blockTarget) return;
    try {
      await blockTask(blockTarget, blockNote);
      message.success("已标记阻塞");
      setBlockTarget(null);
      setBlockNote("");
      await reloadDetail();
    } catch (err) {
      message.error((err as Error).message);
    }
  };

  const columns = useMemo(
    () => [
      { title: "航班号", dataIndex: "flight_no" },
      { title: "机号", dataIndex: "aircraft_reg", width: 90 },
      { title: "机位", dataIndex: "stand_no", width: 80 },
      {
        title: "到站",
        dataIndex: "arrival_time",
        render: (v: string) => formatDateTime(v),
      },
      {
        title: "离港",
        dataIndex: "departure_time",
        render: (v: string) => formatDateTime(v),
      },
      {
        title: "状态",
        render: (_: unknown, t: FlightTurnaround) => (
          <StatusBadge kind="turnaround" value={t.turnaround_status} />
        ),
      },
      {
        title: "延误",
        render: (_: unknown, t: FlightTurnaround) => (
          <DelayTag minutes={t.total_delay_minutes} reason={t.delay_reason} />
        ),
      },
      {
        title: "操作",
        width: 90,
        render: (_: unknown, t: FlightTurnaround) => (
          <Button type="link" size="small" onClick={() => openDetail(t.id)}>
            保障详情
          </Button>
        ),
      },
    ],
    [openDetail],
  );

  return (
    <div>
      <Card
        size="small"
        title="航班过站"
        extra={
          isDispatcher && (
            <Button type="primary" onClick={() => setRegisterOpen(true)}>
              航班到站登记
            </Button>
          )
        }
      >
        <Table rowKey="id" loading={loading} dataSource={rows} columns={columns} pagination={false} />
      </Card>

      <Modal
        title="航班到站登记"
        open={registerOpen}
        onCancel={() => setRegisterOpen(false)}
        onOk={submitRegister}
        confirmLoading={busy}
        destroyOnClose
        afterOpenChange={(open) => {
          if (open) {
            const form = createTurnaroundForm();
            registerForm.setFieldsValue({
              flight_no: "",
              aircraft_reg: "",
              stand_no: "",
              range: [dayjs(form.arrival_time), dayjs(form.departure_time)],
            });
          }
        }}
      >
        <Form form={registerForm} layout="vertical">
          <Row gutter={8}>
            <Col span={12}>
              <Form.Item name="flight_no" label="航班号" rules={[{ required: true }]}>
                <Input placeholder="如 CA101" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="aircraft_reg" label="机号" rules={[{ required: true }]}>
                <Input placeholder="如 B-6271" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="stand_no" label="机位" rules={[{ required: true }]}>
            <Input placeholder="如 203" />
          </Form.Item>
          <Form.Item name="range" label="到站 ~ 离港时间" rules={[{ required: true }]}>
            <DatePicker.RangePicker
              showTime={{ format: "HH:mm" }}
              format="MM-DD HH:mm"
              style={{ width: "100%" }}
            />
          </Form.Item>
        </Form>
      </Modal>

      <Drawer
        width={760}
        open={detailOpen}
        onClose={() => setDetailOpen(false)}
        title={
          detail ? (
            <Space>
              航班 {detail.flight_no}（{detail.stand_no} 机位）
              <StatusBadge kind="turnaround" value={detail.turnaround_status} />
              <DelayTag minutes={detail.total_delay_minutes} />
            </Space>
          ) : (
            ""
          )
        }
        extra={
          <Space>
            {isDispatcher && detail && (
              <>
                <Button onClick={() => {
                  const form = createDelayForm();
                  delayForm.setFieldsValue(form);
                  setDelayOpen(true);
                }}>
                  登记延误
                </Button>
                <Button type="primary" loading={busy} onClick={onGeneratePlan}>
                  生成保障任务与预约
                </Button>
              </>
            )}
          </Space>
        }
      >
        {detail && (
          <>
            <Row gutter={8} style={{ marginBottom: 8 }}>
              <Col span={6}>
                <Card size="small">
                  <div>到站</div>
                  <strong>{formatFullDateTime(detail.arrival_time)}</strong>
                </Card>
              </Col>
              <Col span={6}>
                <Card size="small">
                  <div>离港</div>
                  <strong>{formatFullDateTime(detail.departure_time)}</strong>
                </Card>
              </Col>
              <Col span={6}>
                <Card size="small">
                  <div>任务进度</div>
                  <Progress percent={progress.percent} size="small" style={{ marginBottom: 0 }} />
                </Card>
              </Col>
              <Col span={6}>
                <Card size="small">
                  <div>待处理预约</div>
                  <strong style={{ color: detail.pending_bookings ? "#cf1322" : undefined }}>
                    {detail.pending_bookings}
                  </strong>
                </Card>
              </Col>
            </Row>
            {detail.delay_reason && (
              <Card size="small" style={{ marginBottom: 8 }}>
                延误原因：{detail.delay_reason}（累计 +{detail.total_delay_minutes} 分钟）
              </Card>
            )}

            <ConflictPanel
              conflicts={planConflicts}
              title={`整批未保存 · ${planConflicts.length} 项冲突，任务和预约均未写入`}
              banner="请处理资源占用（释放/上线/调整维护窗口）后重新生成；每条冲突逐项列出原因。"
            />
            <ConflictPanel
              conflicts={delayConflicts}
              title="重排产生待处理预约"
            />
            {lastResult?.saved && (lastResult.rescheduled_tasks !== undefined) && (
              <Card size="small" style={{ marginBottom: 8, marginTop: planConflicts.length ? 12 : 0 }}>
                重排 {lastResult.rescheduled_tasks} 个未签收任务 · 冻结 {lastResult.frozen_tasks}{" "}
                个已签收任务 · {lastResult.pending_bookings} 条预约进入待处理
              </Card>
            )}

            <Card size="small" title="任务时间轴与资源预约">
              <TurnaroundTimeline detail={detail} />
              <Divider style={{ margin: "8px 0" }} />
              <Table
                rowKey="id"
                size="small"
                pagination={false}
                dataSource={[...detail.tasks].sort((a, b) =>
                  a.planned_start.localeCompare(b.planned_start),
                )}
                columns={[
                  {
                    title: "任务",
                    render: (_, t) => (
                      <Space>
                        <StatusBadge kind="taskType" value={t.task_type} />
                        <TeamTag teamId={t.team_id} />
                      </Space>
                    ),
                  },
                  {
                    title: "计划/截止",
                    render: (_, t) => `${formatDateTime(t.planned_start)} ~ ${dayjs(t.deadline).format("HH:mm")}`,
                  },
                  { title: "状态", render: (_, t) => <StatusBadge kind="task" value={t.status} /> },
                  {
                    title: "资源预约",
                    render: (_, t) => {
                      const b = detail.bookings.find((x) => x.task_id === t.id);
                      if (!b) return "—";
                      return (
                        <Space size={4}>
                          <span>{b.resource_code ?? `#${b.resource_id}`}</span>
                          <StatusBadge kind="booking" value={b.booking_status} />
                        </Space>
                      );
                    },
                  },
                  {
                    title: "签收/完成",
                    width: 190,
                    render: (_, t) => (
                      <Space size={4}>
                        {t.status === "PLANNED" && isTeam && (
                          <Button type="link" size="small" onClick={() => onAccept(t.id)}>
                            签收
                          </Button>
                        )}
                        {t.status === "ACCEPTED" && isTeam && (
                          <>
                            <Button type="link" size="small" onClick={() => onComplete(t.id)}>
                              完成
                            </Button>
                            <Button
                              type="link"
                              size="small"
                              danger
                              onClick={() => setBlockTarget(t.id)}
                            >
                              阻塞
                            </Button>
                          </>
                        )}
                        {t.signed_at && (
                          <span style={{ color: "rgba(0,0,0,0.45)", fontSize: 12 }}>
                            已冻结 {dayjs(t.signed_at).format("MM-DD HH:mm")}
                          </span>
                        )}
                      </Space>
                    ),
                  },
                ]}
              />
            </Card>
          </>
        )}
      </Drawer>

      <Modal
        title="登记延误并重排未完成任务"
        open={delayOpen}
        onCancel={() => setDelayOpen(false)}
        onOk={onRegisterDelay}
        confirmLoading={busy}
        destroyOnClose
        okText="登记并重排"
      >
        <Form form={delayForm} layout="vertical">
          <Row gutter={8}>
            <Col span={12}>
              <Form.Item name="delay_type" label="延误类型" rules={[{ required: true }]}>
                <select
                  style={{ width: "100%", height: 32, borderRadius: 6, border: "1px solid #d9d9d9", paddingInline: 8 }}
                  defaultValue="LATE_ARRIVAL"
                >
                  {DELAY_TYPES.map((d) => (
                    <option key={d.value} value={d.value}>
                      {d.label}
                    </option>
                  ))}
                </select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="minutes" label="增加分钟数" rules={[{ required: true }]}>
                <Input type="number" min={1} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="root_cause" label="根因说明">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="responsibility_team" label="责任班组/方">
            <Input placeholder="如 AIRLINE / TEAM-FUEL" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={`标记任务 ${blockTarget ?? ""} 阻塞`}
        open={!!blockTarget}
        onOk={onBlock}
        onCancel={() => setBlockTarget(null)}
        okText="确认阻塞"
      >
        <Input.TextArea
          rows={2}
          placeholder="阻塞原因，如：航食车未到位"
          value={blockNote}
          onChange={(e) => setBlockNote(e.target.value)}
        />
      </Modal>
    </div>
  );
};

import { Card, Col, Row, List, Tag, Button, Space, Progress, Alert } from "antd";
import {
  ClockCircleOutlined,
  ExclamationCircleOutlined,
  FieldTimeOutlined,
  CheckCircleOutlined,
  WarningOutlined
} from "@ant-design/icons";
import { useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../stores/store";
import { fetchTurnarounds } from "../stores/turnaroundSlice";
import { fetchTasks } from "../stores/taskSlice";
import { fetchPendingBookings } from "../stores/bookingSlice";
import { fetchDelays } from "../stores/delaySlice";
import { reportApi } from "../api/report";
import { useState } from "react";
import type { DashboardStats, AuditLog } from "../types/api";
import { StatCard } from "../components/common/StatCard";
import { StatusBadge } from "../components/common/StatusBadge";
import { DelayTag } from "../components/common/DelayTag";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { formatDateTime, formatMinutes } from "../utils/formatters";
import { TURNAROUND_STATUS_TEXT } from "../constants/TurnaroundStatus";

function ProgressRow({ turnaroundId, label }: { turnaroundId: number; label: string }) {
  const tasks = useAppSelector((s) => s.tasks.rows.filter((t) => t.turnaround_id === turnaroundId));
  const progress = useTurnaroundProgress(tasks);
  return (
    <div>
      <div className="muted">{label}</div>
      <Progress percent={progress.percent} size="small" status={progress.ready ? "success" : "active"} />
      <Space size={4} wrap>
        {progress.overdue > 0 && <Tag color="red">超时 {progress.overdue}</Tag>}
        {progress.blocked > 0 && <Tag color="orange">阻塞 {progress.blocked}</Tag>}
        <Tag>{progress.finished}/{progress.total} 完成</Tag>
      </Space>
    </div>
  );
}

export function DashboardPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const turnarounds = useAppSelector((s) => s.turnarounds.rows);
  const pendingBookings = useAppSelector((s) => s.bookings.pending);
  const delays = useAppSelector((s) => s.delays.rows);
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [logs, setLogs] = useState<AuditLog[]>([]);

  useEffect(() => {
    dispatch(fetchTurnarounds());
    dispatch(fetchTasks());
    dispatch(fetchPendingBookings());
    dispatch(fetchDelays());
    reportApi.dashboard().then(setStats).catch(() => undefined);
    reportApi.auditLogs().then(setLogs).catch(() => undefined);
  }, [dispatch]);

  return (
    <div className="page-grid">
      <Row gutter={[16, 16]}>
        <Col xs={12} md={6}>
          <StatCard label="在保航班" value={stats?.in_service ?? "-"} suffix={`/ ${stats?.turnaround_total ?? 0}`} />
        </Col>
        <Col xs={12} md={6}>
          <StatCard label="延误航班" value={stats?.delayed ?? 0} tone="danger" icon={<FieldTimeOutlined />} />
        </Col>
        <Col xs={12} md={6}>
          <StatCard label="超时任务" value={stats?.task_overdue ?? 0} tone="warning" icon={<ExclamationCircleOutlined />} />
        </Col>
        <Col xs={12} md={6}>
          <StatCard label="待处理预约" value={stats?.pending_bookings ?? 0} tone="danger" icon={<ClockCircleOutlined />} />
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={15}>
          <Card
            title="航班进度"
            extra={<Link to="/turnarounds">全部航班 →</Link>}
          >
            <List
              dataSource={turnarounds.slice(0, 5)}
              locale={{ emptyText: <span className="muted">暂无航班</span> }}
              renderItem={(row) => (
                <List.Item
                  className="clickable"
                  onClick={() => navigate(`/turnarounds/${row.id}`)}
                  extra={<StatusBadge value={row.turnaround_status} />}
                >
                  <List.Item.Meta
                    title={
                      <Space>
                        <strong>{row.flight_no}</strong>
                        <span className="muted">{row.aircraft_reg} · 机位 {row.stand_no}</span>
                        {row.accumulated_delay > 0 && <DelayTag type={row.delay_reason === "空管流控" ? "ATC" : "OTHER"} minutes={row.accumulated_delay} />}
                      </Space>
                    }
                    description={
                      <div>
                        <div className="muted">
                          到站 {formatDateTime(row.arrival_time)} · 离港 {formatDateTime(row.departure_time)} ·{" "}
                          {TURNAROUND_STATUS_TEXT[row.turnaround_status]}
                        </div>
                        <ProgressRow turnaroundId={row.id} label="保障任务完成率" />
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
        <Col xs={24} lg={9}>
          <Card title="异常与待处理" className="side-card">
            {pendingBookings.length > 0 && (
              <Alert
                type="error"
                showIcon
                className="mb-12"
                icon={<WarningOutlined />}
                message={`${pendingBookings.length} 个资源预约被挤出待处理`}
                description={
                  <Button size="small" type="primary" danger onClick={() => navigate("/resources")}>
                    去资源调度改期
                  </Button>
                }
              />
            )}
            <List
              size="small"
              header={<strong>最近延误事件</strong>}
              dataSource={delays.slice(0, 5)}
              locale={{ emptyText: <span className="muted">暂无延误登记</span> }}
              renderItem={(d) => (
                <List.Item>
                  <Space wrap>
                    <DelayTag type={d.delay_type} minutes={d.minutes} />
                    <span>航班 #{d.turnaround_id}</span>
                    <span className="muted">{d.root_cause || "—"}</span>
                  </Space>
                </List.Item>
              )}
            />
            <List
              size="small"
              className="mt-12"
              header={<strong>操作日志</strong>}
              dataSource={logs.slice(0, 6)}
              renderItem={(log) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={<CheckCircleOutlined style={{ color: "#52c41a" }} />}
                    title={<span className="log-detail">{log.detail || log.action}</span>}
                    description={
                      <span className="muted">
                        {log.actor} · {formatDateTime(log.created_at)}
                      </span>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>

      {stats && stats.delay_minutes > 0 && (
        <Alert
          type="info"
          showIcon
          message={`今日累计延误 ${formatMinutes(stats.delay_minutes)}，已完成任务 ${stats.task_finished}/${stats.task_total}，阻塞 ${stats.task_blocked} 项。`}
        />
      )}
    </div>
  );
}

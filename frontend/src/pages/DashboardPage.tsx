import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Card, Row, Col, Typography, Button, Table, Tag, App } from "antd";
import { useNavigate } from "react-router-dom";
import { StatCard } from "../components/common/StatCard";
import { StatusBadge } from "../components/common/StatusBadge";
import { DelayTag } from "../components/common/DelayTag";
import { PendingBookingList } from "../components/common/PendingBookingList";
import { getDashboardStats, getPendingBookings, getAuditLogs } from "../api/Dashboard";
import { listTurnarounds } from "../api/FlightTurnaround";
import type {
  DashboardStats as Stats,
  ResourceBooking,
  FlightTurnaround,
  AuditLogEntry,
} from "../types/entities";
import { formatDateTime } from "../utils/formatters";

export const DashboardPage: React.FC = () => {
  const { message } = App.useApp();
  const navigate = useNavigate();
  const [stats, setStats] = useState<Stats | null>(null);
  const [turns, setTurns] = useState<FlightTurnaround[]>([]);
  const [pending, setPending] = useState<ResourceBooking[]>([]);
  const [logs, setLogs] = useState<AuditLogEntry[]>([]);
  const [loading, setLoading] = useState(false);

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const [s, t, p, l] = await Promise.all([
        getDashboardStats(),
        listTurnarounds({ page: 1, page_size: 100 }),
        getPendingBookings(),
        getAuditLogs({ page: 1, page_size: 12 }),
      ]);
      setStats(s);
      setTurns(t.items);
      setPending(p.items);
      setLogs(l.items);
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const flightByTurn = useMemo(
    () => new Map(turns.map((t) => [t.id, t.flight_no])),
    [turns],
  );

  return (
    <div>
      <Row gutter={12} style={{ marginBottom: 12 }}>
        <Col span={6}>
          <StatCard label="在站保障航班" value={stats?.active_turnarounds ?? 0} loading={loading} />
        </Col>
        <Col span={6}>
          <StatCard label="未完成任务" value={stats?.open_tasks ?? 0} loading={loading} />
        </Col>
        <Col span={6}>
          <StatCard
            label="待处理资源预约"
            value={stats?.pending_bookings ?? 0}
            loading={loading}
            danger={(stats?.pending_bookings ?? 0) > 0}
          />
        </Col>
        <Col span={6}>
          <StatCard label="累计延误分钟" value={stats?.total_delay_minutes ?? 0} suffix="分钟" loading={loading} />
        </Col>
      </Row>

      <Row gutter={12}>
        <Col span={14}>
          <Card
            size="small"
            title="在站航班"
            extra={<Button type="link" onClick={() => navigate("/turnarounds")}>全部航班</Button>}
          >
            <Table
              rowKey="id"
              size="small"
              pagination={false}
              dataSource={turns}
              columns={[
                { title: "航班", dataIndex: "flight_no" },
                { title: "机位", dataIndex: "stand_no", width: 80 },
                {
                  title: "到站 / 离港",
                  render: (_, t) => (
                    <span>
                      {formatDateTime(t.arrival_time)} / {formatDateTime(t.departure_time)}
                    </span>
                  ),
                },
                {
                  title: "状态",
                  render: (_, t) => <StatusBadge kind="turnaround" value={t.turnaround_status} />,
                },
                {
                  title: "延误",
                  render: (_, t) => <DelayTag minutes={t.total_delay_minutes} reason={t.delay_reason} />,
                },
                {
                  title: "",
                  width: 90,
                  render: (_, t) => (
                    <Button type="link" size="small" onClick={() => navigate(`/turnarounds?focus=${t.id}`)}>
                      详情
                    </Button>
                  ),
                },
              ]}
            />
          </Card>

          <Card size="small" title="最近操作日志" style={{ marginTop: 12 }}>
            <Table
              rowKey="id"
              size="small"
              pagination={false}
              dataSource={logs}
              columns={[
                { title: "时间", dataIndex: "created_at", width: 130, render: formatDateTime },
                { title: "操作人", dataIndex: "actor", width: 110 },
                { title: "内容", dataIndex: "action" },
              ]}
            />
          </Card>
        </Col>
        <Col span={10}>
          <Card
            size="small"
            title={
              <span>
                <Tag color="warning">待处理</Tag> 被挤出的资源预约
              </span>
            }
            extra={
              <Button type="link" size="small" onClick={refresh}>
                刷新
              </Button>
            }
          >
            <Typography.Paragraph type="secondary" style={{ marginBottom: 8 }}>
              延误重排后与维护/离线/其它航班冲突的预约会进入此列表；调整时段并通过校验后恢复已确认。
            </Typography.Paragraph>
            <PendingBookingList rows={pending} flightByTurn={flightByTurn} onChanged={refresh} />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

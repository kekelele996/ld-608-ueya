import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Card, Table, Tag, Button, Select, Space, Statistic, Row, Col, App } from "antd";
import { useNavigate } from "react-router-dom";
import type { DelayEvent } from "../types/entities";
import { listDelays, resolveDelay } from "../api/DelayEvent";
import { listTurnarounds } from "../api/FlightTurnaround";
import type { FlightTurnaround } from "../types/entities";
import { DelayTag } from "../components/common/DelayTag";
import { useRole } from "../hooks/useRole";
import { formatDateTime } from "../utils/formatters";
import { DELAY_TYPES } from "../types/auth";

export const DelaysPage: React.FC = () => {
  const { message } = App.useApp();
  const navigate = useNavigate();
  const { isDispatcher } = useRole();
  const [rows, setRows] = useState<DelayEvent[]>([]);
  const [turns, setTurns] = useState<FlightTurnaround[]>([]);
  const [typeFilter, setTypeFilter] = useState("");
  const [loading, setLoading] = useState(false);

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const [d, t] = await Promise.all([
        listDelays({ page: 1, page_size: 200 }),
        listTurnarounds({ page: 1, page_size: 200 }),
      ]);
      setRows(d.items);
      setTurns(t.items);
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const filtered = useMemo(
    () => rows.filter((r) => (typeFilter ? r.delay_type === typeFilter : true)),
    [rows, typeFilter],
  );
  const totalMinutes = rows.reduce((sum, r) => sum + r.minutes, 0);
  const openCount = rows.filter((r) => !r.resolved_at).length;

  const onResolve = async (id: number) => {
    try {
      await resolveDelay(id);
      message.success(`延误事件 ${id} 已关闭归因`);
      refresh();
    } catch (err) {
      message.error((err as Error).message);
    }
  };

  return (
    <Card
      size="small"
      title="延误归因"
      extra={
        <Space>
          <Select
            size="small"
            value={typeFilter}
            onChange={setTypeFilter}
            style={{ width: 140 }}
            options={[{ value: "", label: "全部类型" }, ...DELAY_TYPES]}
          />
          {isDispatcher && (
            <Button type="primary" size="small" onClick={() => navigate("/turnarounds")}>
              去登记延误
            </Button>
          )}
        </Space>
      }
    >
      <Row gutter={12} style={{ marginBottom: 12 }}>
        <Col span={8}>
          <Card size="small">
            <Statistic title="延误事件数" value={rows.length} />
          </Card>
        </Col>
        <Col span={8}>
          <Card size="small">
            <Statistic title="累计影响分钟" value={totalMinutes} suffix="分钟" />
          </Card>
        </Col>
        <Col span={8}>
          <Card size="small">
            <Statistic title="待归因" value={openCount} valueStyle={{ color: openCount ? "#cf1322" : undefined }} />
          </Card>
        </Col>
      </Row>
      <Table
        rowKey="id"
        size="small"
        loading={loading}
        dataSource={filtered}
        pagination={{ pageSize: 12 }}
        columns={[
          { title: "事件#", dataIndex: "id", width: 70, render: (v) => `#${v}` },
          {
            title: "航班",
            render: (_, d) => d.flight_no ?? `#${d.turnaround_id}`,
          },
          {
            title: "类型",
            render: (_, d) => (
              <Tag>{DELAY_TYPES.find((t) => t.value === d.delay_type)?.label ?? d.delay_type}</Tag>
            ),
          },
          { title: "分钟", render: (_, d) => <DelayTag minutes={d.minutes} /> },
          { title: "根因", dataIndex: "root_cause", render: (v: string) => v || "—" },
          { title: "责任方", dataIndex: "responsibility_team", render: (v: string) => v || "—" },
          { title: "登记时间", dataIndex: "created_at", render: formatDateTime },
          {
            title: "状态",
            render: (_, d) => (d.resolved_at ? <Tag color="success">已归因</Tag> : <Tag color="processing">待归因</Tag>),
          },
          {
            title: "操作",
            width: 100,
            render: (_, d) =>
              !d.resolved_at && isDispatcher ? (
                <Button type="link" size="small" onClick={() => onResolve(d.id)}>
                  关闭归因
                </Button>
              ) : (
                "—"
              ),
          },
        ]}
      />
    </Card>
  );
};

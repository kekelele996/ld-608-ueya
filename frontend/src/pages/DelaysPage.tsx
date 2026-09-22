import { Button, Card, Col, Progress, Row, Space, Statistic, Table, Tag, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useEffect, useMemo } from "react";
import { Link } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../stores/store";
import { fetchDelays, fetchDelayImpacts } from "../stores/delaySlice";
import { fetchTurnarounds } from "../stores/turnaroundSlice";
import { delayEventApi, type DelayImpact } from "../api/DelayEvent";
import { DelayTag } from "../components/common/DelayTag";
import { TimelineList } from "../components/common/TimelineList";
import { usePagination } from "../hooks/usePagination";
import { formatDateTime, formatMinutes } from "../utils/formatters";
import { DELAY_TYPE_TEXT } from "../constants/DelayType";
import { canRole } from "../constants/UserRole";
import { teamText } from "../constants/teams";
import type { DelayEvent } from "../types/DelayEvent";

// ChartPanel：用纯 CSS 条形展示分钟数（不引第三方图表，全部本地数据）。
function ChartPanel({ impacts }: { impacts: DelayImpact[] }) {
  const max = Math.max(1, ...impacts.map((i) => i.total_minutes));
  if (!impacts.length) return <span className="muted">暂无延误统计</span>;
  return (
    <div className="chart-panel">
      {impacts.map((row) => (
        <div key={row.delay_type} className="chart-row">
          <div className="chart-label">{DELAY_TYPE_TEXT[row.delay_type as keyof typeof DELAY_TYPE_TEXT] ?? row.delay_type}</div>
          <Progress
            percent={Math.round((row.total_minutes / max) * 100)}
            size="small"
            strokeColor="#ff4d4f"
            format={() => `${formatMinutes(row.total_minutes)} / ${row.event_count} 次`}
          />
        </div>
      ))}
    </div>
  );
}

export function DelaysPage() {
  const dispatch = useAppDispatch();
  const delays = useAppSelector((s) => s.delays.rows);
  const impacts = useAppSelector((s) => s.delays.impacts);
  const turnarounds = useAppSelector((s) => s.turnarounds.rows);
  const user = useAppSelector((s) => s.auth.user);

  useEffect(() => {
    dispatch(fetchDelays());
    dispatch(fetchDelayImpacts());
    dispatch(fetchTurnarounds());
  }, [dispatch]);

  const turnMap = useMemo(() => new Map(turnarounds.map((t) => [t.id, t])), [turnarounds]);
  const { pageRows, pageSize, page, setPage, total } = usePagination(delays, 8);
  const totalMinutes = delays.reduce((sum, d) => sum + d.minutes, 0);
  const canResolve = canRole(user?.role, "DELAY_RESOLVE");

  const columns: ColumnsType<DelayEvent> = [
    {
      title: "类型",
      dataIndex: "delay_type",
      width: 130,
      render: (v: string, r) => <DelayTag type={v} minutes={r.minutes} />
    },
    {
      title: "航班",
      render: (_, r) => (
        <Link to={`/turnarounds/${r.turnaround_id}`}>
          {turnMap.get(r.turnaround_id)?.flight_no ?? `#${r.turnaround_id}`}
        </Link>
      )
    },
    { title: "分钟数", dataIndex: "minutes", width: 100, render: (v: number) => formatMinutes(v) },
    { title: "责任班组", dataIndex: "responsibility_team", width: 140, render: (v: string) => (v ? teamText(v) : "—") },
    { title: "原因", dataIndex: "root_cause", render: (v: string) => v || "—" },
    {
      title: "登记时间",
      dataIndex: "created_at",
      width: 140,
      render: formatDateTime
    },
    {
      title: "状态",
      width: 110,
      render: (_, r) =>
        r.resolved_at ? <Tag color="green">已关闭 {formatDateTime(r.resolved_at)}</Tag> : <Tag color="red">处理中</Tag>
    },
    {
      title: "操作",
      width: 90,
      render: (_, r) =>
        !r.resolved_at && canResolve ? (
          <Button
            size="small"
            type="link"
            onClick={async () => {
              await delayEventApi.resolve(r.id);
              message.success("延误事件已关闭");
              dispatch(fetchDelays());
            }}
          >
            关闭归因
          </Button>
        ) : null
    }
  ];

  return (
    <div className="page-grid">
      <Row gutter={[16, 16]}>
        <Col xs={12} md={6}>
          <Card>
            <Statistic title="延误事件" value={delays.length} suffix="起" />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic title="累计影响" value={totalMinutes} suffix="分钟" valueStyle={{ color: "#cf1322" }} />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic title="处理中" value={delays.filter((d) => !d.resolved_at).length} suffix="起" />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic title="已关闭" value={delays.filter((d) => d.resolved_at).length} suffix="起" />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={15}>
          <Card title="延误登记列表">
            <Table<DelayEvent>
              rowKey="id"
              size="small"
              columns={columns}
              dataSource={pageRows}
              pagination={{ current: page, pageSize, total, onChange: setPage }}
            />
          </Card>
        </Col>
        <Col xs={24} lg={9}>
          <Space direction="vertical" style={{ width: "100%" }} size={16}>
            <Card title="影响分钟数统计">
              <ChartPanel impacts={impacts} />
            </Card>
            <Card title="延误时间线">
              <TimelineList
                items={delays.slice(0, 8).map((d) => ({
                  color: d.resolved_at ? "green" : "red",
                  children: (
                    <Space size={4} wrap>
                      <DelayTag type={d.delay_type} minutes={d.minutes} />
                      <span>{turnMap.get(d.turnaround_id)?.flight_no ?? `#${d.turnaround_id}`}</span>
                    </Space>
                  )
                }))}
              />
            </Card>
          </Space>
        </Col>
      </Row>
    </div>
  );
}

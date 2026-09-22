import {
  Button,
  Card,
  Col,
  DatePicker,
  Form,
  Input,
  Modal,
  Progress,
  Row,
  Select,
  Space,
  Table,
  Tag
} from "antd";
import type { ColumnsType } from "antd/es/table";
import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import dayjs from "dayjs";
import { useAppDispatch, useAppSelector } from "../stores/store";
import { fetchTurnarounds } from "../stores/turnaroundSlice";
import { fetchTasks } from "../stores/taskSlice";
import { turnaroundApi } from "../api/FlightTurnaround";
import { StatusBadge } from "../components/common/StatusBadge";
import { DelayTag } from "../components/common/DelayTag";
import { usePagination } from "../hooks/usePagination";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { formatDateTime, toApiTime } from "../utils/formatters";
import {
  TURNAROUND_STATUSES,
  TURNAROUND_STATUS_TEXT,
  type TurnaroundStatusValue
} from "../constants/TurnaroundStatus";
import { canRole } from "../constants/UserRole";
import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { GroundTask } from "../types/GroundTask";
import { message } from "antd";
import { createTurnaroundForm } from "../constructors/FlightTurnaroundConstructor";

function RowProgress({ turnaroundId }: { turnaroundId: number }) {
  const tasks = useAppSelector((s) => s.tasks.rows.filter((t) => t.turnaround_id === turnaroundId));
  const p = useTurnaroundProgress(tasks);
  return <Progress percent={p.percent} size="small" status={p.ready ? "success" : "active"} />;
}

export function TurnaroundsPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const rows = useAppSelector((s) => s.turnarounds.rows);
  const loading = useAppSelector((s) => s.turnarounds.loading);
  const user = useAppSelector((s) => s.auth.user);
  const [statusFilter, setStatusFilter] = useState<TurnaroundStatusValue | "ALL">("ALL");
  const [createOpen, setCreateOpen] = useState(false);
  const [form] = Form.useForm();
  const { page, setPage, pageRows, pageSize, total } = usePagination(
    useMemo(
      () => (statusFilter === "ALL" ? rows : rows.filter((r) => r.turnaround_status === statusFilter)),
      [rows, statusFilter]
    ),
    6
  );

  useEffect(() => {
    dispatch(fetchTurnarounds());
    dispatch(fetchTasks());
  }, [dispatch]);

  const canCreate = canRole(user?.role, "TURNAROUND_CREATE");

  const columns: ColumnsType<FlightTurnaround> = [
    {
      title: "航班",
      render: (_, r) => (
        <Space direction="vertical" size={0}>
          <strong>{r.flight_no}</strong>
          <span className="muted">
            {r.aircraft_reg} · 机位 {r.stand_no}
          </span>
        </Space>
      )
    },
    { title: "到站", dataIndex: "arrival_time", width: 130, render: formatDateTime },
    { title: "离港", dataIndex: "departure_time", width: 130, render: formatDateTime },
    {
      title: "状态",
      dataIndex: "turnaround_status",
      width: 110,
      render: (value: string) => <StatusBadge value={value} />
    },
    {
      title: "延误",
      width: 140,
      render: (_, r) =>
        r.accumulated_delay > 0 ? <DelayTag type="OTHER" minutes={r.accumulated_delay} /> : <span className="muted">—</span>
    },
    { title: "任务完成率", width: 170, render: (_, r) => <RowProgress turnaroundId={r.id} /> },
    {
      title: "操作",
      width: 110,
      render: (_, r) => (
        <Button type="link" size="small" onClick={() => navigate(`/turnarounds/${r.id}`)}>
          详情/保障
        </Button>
      )
    }
  ];

  return (
    <div className="page-grid">
      <Card
        title="航班过站"
        extra={
          <Space>
            <Select
              value={statusFilter}
              style={{ width: 140 }}
              onChange={(v) => {
                setStatusFilter(v);
                setPage(1);
              }}
              options={[
                { value: "ALL", label: "全部状态" },
                ...TURNAROUND_STATUSES.map((s) => ({ value: s, label: TURNAROUND_STATUS_TEXT[s] }))
              ]}
            />
            {canCreate && (
              <Button type="primary" onClick={() => {
                form.setFieldsValue(createTurnaroundForm());
                setCreateOpen(true);
              }}>
                过站登记
              </Button>
            )}
          </Space>
        }
      >
        <Table<FlightTurnaround>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={pageRows}
          pagination={{ current: page, pageSize, total, onChange: setPage }}
        />
      </Card>

      <Modal
        title="航班过站登记"
        open={createOpen}
        onCancel={() => setCreateOpen(false)}
        onOk={async () => {
          const values = await form.validateFields();
          await turnaroundApi.create({
            flight_no: values.flight_no,
            aircraft_reg: values.aircraft_reg,
            stand_no: values.stand_no,
            arrival_time: toApiTime(values.range[0]),
            departure_time: toApiTime(values.range[1])
          });
          message.success("过站登记成功");
          setCreateOpen(false);
          dispatch(fetchTurnarounds());
        }}
        destroyOnClose
      >
        <Form form={form} layout="vertical" initialValues={createTurnaroundForm()}>
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="flight_no" label="航班号" rules={[{ required: true }]}>
                <Input placeholder="CA1831" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="aircraft_reg" label="机号" rules={[{ required: true }]}>
                <Input placeholder="B-5821" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="stand_no" label="机位号" rules={[{ required: true }]}>
            <Input placeholder="203" />
          </Form.Item>
          <Form.Item name="range" label="到站 ~ 离港" rules={[{ required: true }]}>
            <DatePicker.RangePicker showTime format="MM-DD HH:mm" style={{ width: "100%" }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

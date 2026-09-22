import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Card, Row, Col, Button, Space, Modal, DatePicker, Select, Form, Tag, App } from "antd";
import dayjs from "dayjs";
import type { GroundResource } from "../types/entities";
import { listResources, updateResourceStatus } from "../api/GroundResource";
import { listBookings } from "../api/ResourceBooking";
import { listTurnarounds } from "../api/FlightTurnaround";
import type { ResourceBooking, FlightTurnaround } from "../types/entities";
import { ResourceCalendar } from "../components/common/ResourceCalendar";
import { PendingBookingList } from "../components/common/PendingBookingList";
import { StatusBadge } from "../components/common/StatusBadge";
import { useResourceConflict } from "../hooks/useResourceConflict";
import { useRole } from "../hooks/useRole";
import { formatDateTime, toAPI } from "../utils/formatters";

export const ResourcesPage: React.FC = () => {
  const { message } = App.useApp();
  const { isResource } = useRole();
  const [resources, setResources] = useState<GroundResource[]>([]);
  const [bookings, setBookings] = useState<ResourceBooking[]>([]);
  const [turns, setTurns] = useState<FlightTurnaround[]>([]);
  const [loading, setLoading] = useState(false);
  const [maintTarget, setMaintTarget] = useState<GroundResource | null>(null);
  const [maintForm] = Form.useForm();

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const [resRes, turnRes] = await Promise.all([
        listResources({}),
        listTurnarounds({ page: 1, page_size: 200 }),
      ]);
      const all = await listBookings({});
      setResources(resRes.items);
      setBookings(all.items);
      setTurns(turnRes.items);
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const { summaries } = useResourceConflict(bookings, resources);
  const flightByTurn = useMemo(() => new Map(turns.map((t) => [t.id, t.flight_no])), [turns]);

  const submitMaintenance = async () => {
    if (!maintTarget) return;
    const values = maintForm.getFieldsValue();
    try {
      await updateResourceStatus(maintTarget.id, {
        status: values.status,
        maintenance_start: values.window ? toAPI(values.window[0].format("YYYY-MM-DDTHH:mm")) : undefined,
        maintenance_end: values.window ? toAPI(values.window[1].format("YYYY-MM-DDTHH:mm")) : undefined,
      });
      message.success(`资源 ${maintTarget.resource_code} 已更新`);
      setMaintTarget(null);
      refresh();
    } catch (err) {
      message.error((err as Error).message);
    }
  };

  return (
    <div>
      <Row gutter={12}>
        <Col span={15}>
          <Card size="small" title="资源台账与预约日历" loading={loading} extra={<Button size="small" onClick={refresh}>刷新</Button>}>
            <ResourceCalendar resources={resources} bookings={bookings} onSelectBooking={() => undefined} />
          </Card>

          <Card size="small" title="资源状态维护" style={{ marginTop: 12 }} loading={loading}>
            <Row gutter={[8, 8]}>
              {resources.map((r) => (
                <Col span={12} key={r.id}>
                  <Card size="small" styles={{ body: { padding: 10 } }}>
                    <Space style={{ justifyContent: "space-between", width: "100%" }}>
                      <div>
                        <strong>{r.resource_code}</strong> <Tag>{r.resource_type}</Tag>
                        <div>
                          <StatusBadge kind="resource" value={r.availability_status} />
                          <span style={{ color: "rgba(0,0,0,0.55)", marginLeft: 8 }}>{r.location ?? "—"}</span>
                        </div>
                        {r.maintenance_due_at && (
                          <small style={{ color: "#d46b08" }}>
                            维护 {formatDateTime(r.maintenance_due_at)} ~{" "}
                            {r.maintenance_end ? dayjs(r.maintenance_end).format("MM-DD HH:mm") : ""}
                          </small>
                        )}
                      </div>
                      {isResource && (
                        <Button size="small" onClick={() => {
                          setMaintTarget(r);
                          maintForm.setFieldsValue({
                            status: r.availability_status,
                            window: r.maintenance_due_at
                              ? [dayjs(r.maintenance_due_at), dayjs(r.maintenance_end)]
                              : undefined,
                          });
                        }}>
                          维护/状态
                        </Button>
                      )}
                    </Space>
                  </Card>
                </Col>
              ))}
            </Row>
          </Card>
        </Col>
        <Col span={9}>
          <Card
            size="small"
            title={<Space><Tag color="warning">{summaries.length}</Tag>待处理预约（被挤出）</Space>}
          >
            <PendingBookingList
              rows={summaries.map((s) => s.booking)}
              flightByTurn={flightByTurn}
              onChanged={refresh}
            />
          </Card>
        </Col>
      </Row>

      <Modal
        title={maintTarget ? `资源 ${maintTarget.resource_code} 状态/维护窗口` : ""}
        open={!!maintTarget}
        onCancel={() => setMaintTarget(null)}
        onOk={submitMaintenance}
        okText="保存"
        destroyOnClose
      >
        <Form form={maintForm} layout="vertical">
          <Form.Item name="status" label="可用状态" rules={[{ required: true }]}>
            <Select
              options={[
                { value: "AVAILABLE", label: "可用" },
                { value: "BOOKED", label: "占用中" },
                { value: "MAINTENANCE", label: "维护中" },
                { value: "OFFLINE", label: "离线" },
              ]}
            />
          </Form.Item>
          <Form.Item name="window" label="维护时段（可选）">
            <DatePicker.RangePicker
              showTime={{ format: "HH:mm" }}
              format="MM-DD HH:mm"
              style={{ width: "100%" }}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

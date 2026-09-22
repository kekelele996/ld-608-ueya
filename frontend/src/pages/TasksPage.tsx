import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Card, Table, Select, Space, Button, Modal, Input, Tag, App } from "antd";
import type { GroundTask } from "../types/entities";
import { listTasks, acceptTask, completeTask, blockTask } from "../api/GroundTask";
import { listTurnarounds } from "../api/FlightTurnaround";
import type { FlightTurnaround } from "../types/entities";
import { StatusBadge } from "../components/common/StatusBadge";
import { TeamTag } from "../components/common/TeamTag";
import { useRole } from "../hooks/useRole";
import { formatDateTime } from "../utils/formatters";
import { GROUND_TASK_STATUS_TEXT } from "../types/ResourceStatus";

export const TasksPage: React.FC = () => {
  const { message } = App.useApp();
  const { isTeam } = useRole();
  const [rows, setRows] = useState<GroundTask[]>([]);
  const [turns, setTurns] = useState<FlightTurnaround[]>([]);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [blockTarget, setBlockTarget] = useState<GroundTask | null>(null);
  const [note, setNote] = useState("");

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      const [tasks, turnRes] = await Promise.all([
        listTasks({ page: 1, page_size: 200, status: statusFilter }),
        listTurnarounds({ page: 1, page_size: 200 }),
      ]);
      setRows(tasks.items);
      setTurns(turnRes.items);
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  }, [message, statusFilter]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const flightMap = useMemo(() => new Map(turns.map((t) => [t.id, t.flight_no])), [turns]);

  const act = async (fn: () => Promise<unknown>, ok: string) => {
    try {
      await fn();
      message.success(ok);
      refresh();
    } catch (err) {
      message.error((err as Error).message);
    }
  };

  return (
    <Card
      size="small"
      title="地勤任务"
      extra={
        <Select
          size="small"
          value={statusFilter}
          onChange={setStatusFilter}
          style={{ width: 140 }}
          options={[
            { value: "", label: "全部状态" },
            ...Object.entries(GROUND_TASK_STATUS_TEXT).map(([value, label]) => ({ value, label })),
          ]}
        />
      }
    >
      <Table
        rowKey="id"
        size="small"
        loading={loading}
        dataSource={rows}
        pagination={{ pageSize: 12 }}
        columns={[
          { title: "任务#", dataIndex: "id", width: 70, render: (v) => `#${v}` },
          {
            title: "航班",
            render: (_, t) => flightMap.get(t.turnaround_id) ?? `#${t.turnaround_id}`,
          },
          { title: "任务", render: (_, t) => <StatusBadge kind="taskType" value={t.task_type} /> },
          { title: "班组", render: (_, t) => <TeamTag teamId={t.team_id} /> },
          {
            title: "计划开始",
            dataIndex: "planned_start",
            render: formatDateTime,
          },
          { title: "截止", dataIndex: "deadline", render: formatDateTime },
          {
            title: "状态",
            render: (_, t) => (
              <Space size={4}>
                <StatusBadge kind="task" value={t.status} />
                {t.signed_at && <Tag color="blue">时点已冻结</Tag>}
              </Space>
            ),
          },
          {
            title: "阻塞原因",
            dataIndex: "blocker_note",
            render: (v: string) => v ?? "—",
          },
          {
            title: "操作",
            width: 210,
            render: (_, t) =>
              isTeam ? (
                <Space size={2}>
                  {t.status === "PLANNED" && (
                    <Button type="link" size="small" onClick={() => act(() => acceptTask(t.id), `任务 ${t.id} 已签收`)}>
                      签收
                    </Button>
                  )}
                  {t.status === "ACCEPTED" && (
                    <>
                      <Button type="link" size="small" onClick={() => act(() => completeTask(t.id), `任务 ${t.id} 已完成`)}>
                        完成
                      </Button>
                      <Button danger type="link" size="small" onClick={() => setBlockTarget(t)}>
                        阻塞
                      </Button>
                    </>
                  )}
                </Space>
              ) : (
                <span style={{ color: "rgba(0,0,0,0.35)" }}>班组/调度可操作</span>
              ),
          },
        ]}
      />

      <Modal
        title={blockTarget ? `阻塞任务 #${blockTarget.id}` : ""}
        open={!!blockTarget}
        onCancel={() => setBlockTarget(null)}
        onOk={() =>
          blockTarget &&
          act(
            () => blockTask(blockTarget.id, note),
            "已标记阻塞",
          ).then(() => {
            setBlockTarget(null);
            setNote("");
          })
        }
        okText="确认阻塞"
      >
        <Input.TextArea rows={3} value={note} onChange={(e) => setNote(e.target.value)} placeholder="阻塞原因" />
      </Modal>
    </Card>
  );
};

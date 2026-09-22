import { Steps } from "antd";
import {
  ClockCircleOutlined,
  StopOutlined,
  ToolOutlined,
  CheckCircleOutlined,
  RocketOutlined
} from "@ant-design/icons";
import { useMemo } from "react";
import type { GroundTask } from "../../types/GroundTask";
import { formatClock } from "../../utils/formatters";
import { GROUND_TASK_TYPE_TEXT } from "../../constants/GroundTaskType";

// TurnaroundTimeline：过站详情/过站列表共用的任务时间轴。
// 已签收任务固定标注签收时点；待处理预约任务显示警告。
export function TurnaroundTimeline({
  tasks,
  pendingTaskIds
}: {
  tasks: GroundTask[];
  pendingTaskIds?: number[];
}) {
  const ordered = useMemo(
    () => [...tasks].sort((a, b) => a.planned_start.localeCompare(b.planned_start)),
    [tasks]
  );
  const pendingSet = useMemo(() => new Set(pendingTaskIds ?? []), [pendingTaskIds]);

  const items = ordered.map((task) => {
    const finished = task.status === "FINISHED";
    const signed = task.status === "SIGNED";
    const blocked = task.status === "BLOCKED";
    return {
      title: `${GROUND_TASK_TYPE_TEXT[task.task_type] ?? task.task_type}
        ${pendingSet.has(task.id) ? "（资源待处理）" : ""}`,
      status: (finished ? "finish" : blocked ? "error" : signed ? "process" : "wait") as
        | "finish"
        | "error"
        | "process"
        | "wait",
      icon: finished ? (
        <CheckCircleOutlined />
      ) : blocked ? (
        <StopOutlined />
      ) : signed ? (
        <ToolOutlined />
      ) : (
        <ClockCircleOutlined />
      ),
      description: (
        <div>
          <div>
            {formatClock(task.planned_start)} ~ {formatClock(task.planned_end)} · 截止{" "}
            {formatClock(task.deadline)}
          </div>
          {task.signed_by && <div className="muted">签收：{task.signed_by}</div>}
          {task.blocker_note && <div className="danger-text">{task.blocker_note}</div>}
        </div>
      )
    };
  });

  return (
    <Steps
      direction="vertical"
      size="small"
      current={ordered.findIndex((t) => t.status !== "FINISHED")}
      items={items}
    />
  );
}

export { RocketOutlined };

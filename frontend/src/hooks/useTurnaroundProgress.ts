import { useMemo } from "react";
import type { GroundTask } from "../types/GroundTask";
import { GROUND_TASK_STATUS_TEXT } from "../constants/GroundTaskStatus";

// useTurnaroundProgress：计算单个过站的任务完成率与超时任务，
// Dashboard / 过站列表 / 过站详情三处共用同一口径。
export function useTurnaroundProgress(tasks: GroundTask[] = []) {
  return useMemo(() => {
    const total = tasks.length;
    const finished = tasks.filter((t) => t.status === "FINISHED").length;
    const signed = tasks.filter((t) => t.status === "SIGNED").length;
    const blocked = tasks.filter((t) => t.status === "BLOCKED").length;
    const planned = tasks.filter((t) => t.status === "PLANNED").length;
    const now = Date.now();
    const overdue = tasks.filter(
      (t) => t.status !== "FINISHED" && new Date(t.deadline).getTime() < now
    ).length;
    const percent = total === 0 ? 0 : Math.round((finished / total) * 100);
    const nextTask = [...tasks]
      .filter((t) => t.status === "PLANNED" || t.status === "SIGNED")
      .sort((a, b) => a.planned_start.localeCompare(b.planned_start))[0];

    return {
      total,
      finished,
      signed,
      blocked,
      planned,
      overdue,
      percent,
      nextTask,
      statusText: GROUND_TASK_STATUS_TEXT,
      ready: total > 0 && finished === total
    };
  }, [tasks]);
}

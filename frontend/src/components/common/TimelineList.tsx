import { Timeline } from "antd";
import type { TimelineItemProps } from "antd";

// TimelineList：延误页/任务页共用的轻量时间线壳。
export function TimelineList({ items }: { items: TimelineItemProps[] }) {
  if (!items.length) return <span className="muted">暂无记录</span>;
  return <Timeline items={items} />;
}

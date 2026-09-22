import React from "react";
import { Timeline, Empty } from "antd";

interface TimelineItem {
  id: number;
  color?: string;
  title: React.ReactNode;
  description?: React.ReactNode;
}

interface TimelineListProps {
  items: TimelineItem[];
  emptyText?: string;
}

// TimelineList is shared by tasks page and delays page.
export const TimelineList: React.FC<TimelineListProps> = ({ items, emptyText = "暂无记录" }) => {
  if (!items.length) return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={emptyText} />;
  return (
    <Timeline
      items={items.map((item) => ({
        key: item.id,
        color: item.color ?? "gray",
        children: (
          <div>
            <div>{item.title}</div>
            {item.description && <div style={{ color: "rgba(0,0,0,0.55)" }}>{item.description}</div>}
          </div>
        ),
      }))}
    />
  );
};

import { Empty } from "antd";
import { InboxOutlined } from "@ant-design/icons";

export function EmptyState({ title = "暂无数据", hint }: { title?: string; hint?: string }) {
  return (
    <div className="empty-state">
      <Empty image={<InboxOutlined style={{ fontSize: 36 }} />} description={title} />
      {hint && <p className="muted">{hint}</p>}
    </div>
  );
}

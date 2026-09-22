import { Alert } from "antd";
import { WarningOutlined } from "@ant-design/icons";
import { CONFLICT_CODE_TEXT } from "../../constants/errorCodes";
import type { ConflictItem } from "../../types/api";

// ConflictBadge / 冲突面板：整批生成被拒或预约被挤出时逐项展示原因。
export function ConflictBadge({ code, message }: { code?: string; message?: string }) {
  if (!message) return null;
  return (
    <span className="conflict-badge" title={message}>
      <WarningOutlined style={{ color: "#cf1322", marginRight: 4 }} />
      {code ? CONFLICT_CODE_TEXT[code] ?? code : "冲突"}
    </span>
  );
}

export function ConflictPanel({
  title = "冲突明细（整批未保存）",
  items
}: {
  title?: string;
  items: ConflictItem[];
}) {
  if (!items.length) return null;
  return (
    <Alert
      type="error"
      showIcon
      className="conflict-panel"
      message={title}
      description={
        <ul className="conflict-list">
          {items.map((item, index) => (
            <li key={`${item.code}-${index}`}>
              <ConflictBadge code={item.code} />
              <span className="conflict-msg">{item.message}</span>
            </li>
          ))}
        </ul>
      }
    />
  );
}

import { Tag } from "antd";
import { DELAY_TYPE_TEXT } from "../../constants/DelayType";
import { formatMinutes } from "../../utils/formatters";

// DelayTag 用于看板/过站列表/延误页：延误类型 + 分钟数。
export function DelayTag({ type, minutes }: { type: string; minutes?: number }) {
  const label = DELAY_TYPE_TEXT[type as keyof typeof DELAY_TYPE_TEXT] ?? type;
  return (
    <Tag color="volcano" className="delay-tag">
      {label}
      {typeof minutes === "number" && minutes > 0 ? ` · ${formatMinutes(minutes)}` : ""}
    </Tag>
  );
}

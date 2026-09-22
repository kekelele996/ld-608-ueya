import { Tag } from "antd";

// StatusBadge 被 5 个页面共用：新增枚举值需同步这里的颜色映射。
const COLOR_MAP: Record<string, string> = {
  // 过站状态
  ARRIVING: "default",
  ON_STAND: "blue",
  IN_SERVICE: "processing",
  READY: "gold",
  DEPARTED: "success",
  DELAYED: "error",
  // 任务状态
  PLANNED: "default",
  SIGNED: "processing",
  FINISHED: "success",
  BLOCKED: "error",
  // 资源状态
  AVAILABLE: "success",
  BOOKED: "processing",
  MAINTENANCE: "warning",
  OFFLINE: "error",
  // 预约状态
  CONFIRMED: "success",
  PENDING: "warning",
  RELEASED: "default"
};

const TEXT_MAP: Record<string, string> = {
  ARRIVING: "即将到站",
  ON_STAND: "已上轮挡",
  IN_SERVICE: "保障中",
  READY: "待放飞",
  DEPARTED: "已离港",
  DELAYED: "延误",
  PLANNED: "待签收",
  SIGNED: "已签收",
  FINISHED: "已完成",
  BLOCKED: "阻塞",
  AVAILABLE: "可用",
  BOOKED: "已占用",
  MAINTENANCE: "维护中",
  OFFLINE: "离线",
  CONFIRMED: "已确认",
  PENDING: "待处理",
  RELEASED: "已释放"
};

export function StatusBadge({ value }: { value: string }) {
  const color = COLOR_MAP[value] ?? "default";
  return <Tag color={color}>{TEXT_MAP[value] ?? value}</Tag>;
}

import dayjs from "dayjs";

// 混合日期 / 状态 / 风险等级格式化逻辑：多个页面与服务共同依赖，
// 修改任何格式都牵一发动全身（与后端 utils/time.go 对应）。
export const formatDateTime = (value?: string | null): string => {
  if (!value) return "-";
  const d = dayjs(value);
  return d.isValid() ? d.format("MM-DD HH:mm") : "-";
};

export const formatFullDateTime = (value?: string | null): string => {
  if (!value) return "-";
  const d = dayjs(value);
  return d.isValid() ? d.format("YYYY-MM-DD HH:mm") : "-";
};

export const formatClock = (value?: string | null): string => {
  if (!value) return "-";
  const d = dayjs(value);
  return d.isValid() ? d.format("HH:mm") : "-";
};

export const formatNumber = (value: number): string =>
  new Intl.NumberFormat("zh-CN").format(value);

export const formatMinutes = (value: number): string => {
  if (value < 60) return `${value} 分钟`;
  const h = Math.floor(value / 60);
  const m = value % 60;
  return m === 0 ? `${h} 小时` : `${h} 小时 ${m} 分`;
};

export const formatRisk = (value: string): string =>
  ({ LOW: "低", MEDIUM: "中", HIGH: "高", CRITICAL: "严重", EXTREME: "极高" }[value] ?? value);

// 供 antd DatePicker 使用的提交格式（后端 RFC3339 解析）。
export const toApiTime = (d: dayjs.Dayjs | null | undefined): string =>
  d ? d.format("YYYY-MM-DDTHH:mm:ssZ") : "";

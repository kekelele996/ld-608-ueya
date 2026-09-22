import { Card, Statistic } from "antd";

// StatCard 看板指标卡，Dashboard 与多个页面共用。
export function StatCard({
  label,
  value,
  suffix,
  tone = "default",
  icon
}: {
  label: string;
  value: number | string;
  suffix?: string;
  tone?: "default" | "danger" | "warning" | "success";
  icon?: React.ReactNode;
}) {
  return (
    <Card className={`stat-card stat-${tone}`}>
      <Statistic title={label} value={value} suffix={suffix} prefix={icon} />
    </Card>
  );
}

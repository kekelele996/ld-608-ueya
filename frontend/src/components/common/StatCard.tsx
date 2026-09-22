import React from "react";
import { Card, Statistic } from "antd";

interface StatCardProps {
  label: string;
  value: number | string;
  suffix?: string;
  loading?: boolean;
  danger?: boolean;
}

// StatCard is shared by dashboard headers on every page.
export const StatCard: React.FC<StatCardProps> = ({ label, value, suffix, loading, danger }) => (
  <Card loading={loading} size="small" styles={{ body: { padding: 16 } }}>
    <Statistic
      title={label}
      value={value}
      suffix={suffix}
      valueStyle={danger ? { color: "#cf1322" } : undefined}
    />
  </Card>
);

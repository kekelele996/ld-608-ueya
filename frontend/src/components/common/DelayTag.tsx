import React from "react";
import { Tag } from "antd";

interface DelayTagProps {
  minutes: number;
  reason?: string;
}

// DelayTag is shared by dashboard, turnarounds list and delays page.
export const DelayTag: React.FC<DelayTagProps> = ({ minutes, reason }) => {
  if (!minutes) return <Tag color="success">正点</Tag>;
  const level = minutes >= 30 ? "error" : minutes >= 15 ? "warning" : "processing";
  return (
    <Tag color={level} title={reason}>
      延误 +{minutes}′
    </Tag>
  );
};

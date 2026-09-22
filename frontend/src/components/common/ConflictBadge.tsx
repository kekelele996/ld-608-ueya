import React from "react";
import { Tag, Tooltip } from "antd";
import { CONFLICT_REASON_COLOR, CONFLICT_REASON_TEXT } from "../../constants/statusText";

interface ConflictBadgeProps {
  code?: string;
  message?: string;
}

// ConflictBadge renders a machine-readable conflict reason with the human
// explanation tooltip. Used on resources page, turnaround detail, pending
// booking list and dashboard.
export const ConflictBadge: React.FC<ConflictBadgeProps> = ({ code, message }) => {
  if (!code) return null;
  const text = CONFLICT_REASON_TEXT[code] ?? code;
  const color = CONFLICT_REASON_COLOR[code] ?? "error";
  return (
    <Tooltip title={message ?? text}>
      <Tag color={color} style={{ marginInlineEnd: 0 }}>
        {text}
      </Tag>
    </Tooltip>
  );
};

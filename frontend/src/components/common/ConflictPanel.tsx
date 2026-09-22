import React from "react";
import { Alert } from "antd";
import type { ConflictItem } from "../../types/entities";
import { ConflictBadge } from "./ConflictBadge";

interface ConflictPanelProps {
  conflicts: ConflictItem[];
  title?: string;
  banner?: string;
}

// ConflictPanel lists every per-item conflict reason returned by a rejected
// batch (plan generation) or a failed booking adjustment. Nothing is saved.
export const ConflictPanel: React.FC<ConflictPanelProps> = ({ conflicts, title, banner }) => {
  if (!conflicts?.length) return null;
  return (
    <Alert
      type="error"
      showIcon
      style={{ marginTop: 12 }}
      message={title ?? `整批未保存 · ${conflicts.length} 项冲突`}
      description={
        <div>
          {banner && <p style={{ marginTop: 0 }}>{banner}</p>}
          <ul style={{ marginBottom: 0, paddingLeft: 18 }}>
            {conflicts.map((item, idx) => (
              <li key={`${item.code}-${idx}`} style={{ marginBottom: 6 }}>
                <ConflictBadge code={item.code} />
                <span>{item.message}</span>
                {item.resource_code ? (
                  <span style={{ color: "rgba(0,0,0,0.55)" }}>（资源 {item.resource_code}）</span>
                ) : null}
              </li>
            ))}
          </ul>
        </div>
      }
    />
  );
};

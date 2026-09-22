import React from "react";
import { Tag } from "antd";

interface TeamTagProps {
  teamId: string;
}

const TEAM_TEXT: Record<string, string> = {
  "TEAM-CLEAN": "保洁班",
  "TEAM-CATER": "航食班",
  "TEAM-BAG": "行李班",
  "TEAM-FUEL": "加油班",
  "TEAM-WATER": "水务班",
  "TEAM-RAMP": "机坪班",
};

// TeamTag is shared by tasks page and turnaround timeline.
export const TeamTag: React.FC<TeamTagProps> = ({ teamId }) => (
  <Tag color="purple">{TEAM_TEXT[teamId] ?? teamId}</Tag>
);

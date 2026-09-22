import { Tag } from "antd";
import { teamText } from "../../constants/teams";

// TeamTag 班组标签：任务页/详情页/延误页共用。
export function TeamTag({ teamId }: { teamId: string }) {
  const palette: Record<string, string> = {
    "TEAM-CLEAN": "cyan",
    "TEAM-CATER": "gold",
    "TEAM-BAG": "geekblue",
    "TEAM-FUEL": "orange",
    "TEAM-WATER": "blue",
    "TEAM-RAMP": "purple"
  };
  return <Tag color={palette[teamId] ?? "default"}>{teamText(teamId)}</Tag>;
}

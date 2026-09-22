import { GROUND_TASK_TYPE_TEXT, type GroundTaskTypeValue } from "./GroundTaskType";

export interface TeamInfo {
  id: string;
  name: string;
  taskKinds: GroundTaskTypeValue[];
}

export const TEAMS: TeamInfo[] = [
  { id: "TEAM-CLEAN", name: "客舱清洁班", taskKinds: ["CLEANING"] },
  { id: "TEAM-CATER", name: "航机配餐班", taskKinds: ["CATERING"] },
  { id: "TEAM-BAG", name: "行李装卸班", taskKinds: ["BAGGAGE"] },
  { id: "TEAM-FUEL", name: "航油保障班", taskKinds: ["REFUEL"] },
  { id: "TEAM-WATER", name: "清水污水班", taskKinds: ["WATER_SERVICE"] },
  { id: "TEAM-RAMP", name: "机坪牵引班", taskKinds: ["PUSHBACK"] }
];

export function teamText(teamId: string): string {
  return TEAMS.find((t) => t.id === teamId)?.name ?? teamId;
}

export function taskTypeText(value: string): string {
  return GROUND_TASK_TYPE_TEXT[value as GroundTaskTypeValue] ?? value;
}

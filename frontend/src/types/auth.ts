export type Role = "DISPATCHER" | "TEAM" | "RESOURCE_MANAGER" | "SUPERVISOR";

export const ROLE_TEXT: Record<Role, string> = {
  DISPATCHER: "地勤调度",
  TEAM: "保障班组",
  RESOURCE_MANAGER: "资源管理员",
  SUPERVISOR: "运行督导",
};

export const ROLE_USERNAMES: { username: string; role: Role; team?: string }[] = [
  { username: "dispatcher", role: "DISPATCHER" },
  { username: "team-bag", role: "TEAM", team: "TEAM-BAG" },
  { username: "team-fuel", role: "TEAM", team: "TEAM-FUEL" },
  { username: "team-ramp", role: "TEAM", team: "TEAM-RAMP" },
  { username: "resource", role: "RESOURCE_MANAGER" },
  { username: "supervisor", role: "SUPERVISOR" },
];

export const DELAY_TYPES: { value: string; label: string }[] = [
  { value: "LATE_ARRIVAL", label: "晚到" },
  { value: "WEATHER", label: "天气" },
  { value: "AIRLINE", label: "航司原因" },
  { value: "SECURITY", label: "安检" },
  { value: "GROUND", label: "地勤原因" },
];

export const USER_ROLES = ["DISPATCHER", "CREW", "RESOURCE", "SUPERVISOR"] as const;
export type UserRoleValue = (typeof USER_ROLES)[number];

export const USER_ROLE_TEXT: Record<UserRoleValue, string> = {
  DISPATCHER: "地勤调度",
  CREW: "班组",
  RESOURCE: "资源管理员",
  SUPERVISOR: "运行督导"
};

// 与后端 constants/rbac.go 的 RBACMatrix 保持一致，用于按钮显隐。
export const RBAC_ACTIONS = [
  "TURNAROUND_CREATE",
  "TURNAROUND_ARRIVE",
  "TURNAROUND_RELEASE",
  "PLAN_GENERATE",
  "TASK_SIGN",
  "TASK_FINISH",
  "TASK_BLOCK",
  "RESOURCE_CREATE",
  "RESOURCE_UPDATE",
  "RESOURCE_SET_MAINTENANCE",
  "BOOKING_ADJUST",
  "BOOKING_CONFIRM",
  "BOOKING_RELEASE",
  "DELAY_REGISTER",
  "DELAY_REPLAN",
  "DELAY_RESOLVE"
] as const;
export type RbacAction = (typeof RBAC_ACTIONS)[number];

export const RBAC_MATRIX: Record<RbacAction, UserRoleValue[]> = {
  TURNAROUND_CREATE: ["DISPATCHER", "SUPERVISOR"],
  TURNAROUND_ARRIVE: ["DISPATCHER", "SUPERVISOR"],
  TURNAROUND_RELEASE: ["DISPATCHER", "SUPERVISOR"],
  PLAN_GENERATE: ["DISPATCHER", "SUPERVISOR"],
  TASK_SIGN: ["CREW", "DISPATCHER", "SUPERVISOR"],
  TASK_FINISH: ["CREW", "DISPATCHER", "SUPERVISOR"],
  TASK_BLOCK: ["CREW", "DISPATCHER", "SUPERVISOR"],
  RESOURCE_CREATE: ["RESOURCE", "SUPERVISOR"],
  RESOURCE_UPDATE: ["RESOURCE", "SUPERVISOR"],
  RESOURCE_SET_MAINTENANCE: ["RESOURCE", "SUPERVISOR"],
  BOOKING_ADJUST: ["RESOURCE", "DISPATCHER", "SUPERVISOR"],
  BOOKING_CONFIRM: ["RESOURCE", "DISPATCHER", "SUPERVISOR"],
  BOOKING_RELEASE: ["RESOURCE", "DISPATCHER", "SUPERVISOR"],
  DELAY_REGISTER: ["DISPATCHER", "SUPERVISOR", "CREW"],
  DELAY_REPLAN: ["DISPATCHER", "SUPERVISOR"],
  DELAY_RESOLVE: ["SUPERVISOR", "DISPATCHER"]
};

export function canRole(role: UserRoleValue | undefined, action: RbacAction): boolean {
  return !!role && RBAC_MATRIX[action].includes(role);
}

export const SEED_ACCOUNTS = [
  { username: "dispatcher", password: "dispatch123", role: "DISPATCHER" as const, hint: "地勤调度" },
  { username: "crew", password: "crew123", role: "CREW" as const, hint: "班组" },
  { username: "resource", password: "resource123", role: "RESOURCE" as const, hint: "资源管理员" },
  { username: "supervisor", password: "super123", role: "SUPERVISOR" as const, hint: "运行督导" }
];

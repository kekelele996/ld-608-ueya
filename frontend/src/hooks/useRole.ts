import { useAppSelector } from "../stores/hooks";
import type { Role } from "../types/auth";

// useRole gates buttons across pages; mirrors backend rbacMiddleware roles.
export function useRole() {
  const { role } = useAppSelector((s) => s.auth);
  const can = (...roles: Role[]) => roles.includes(role as Role);
  return { role, can, isDispatcher: can("DISPATCHER", "SUPERVISOR"), isResource: can("RESOURCE_MANAGER", "SUPERVISOR"), isTeam: can("TEAM", "DISPATCHER", "SUPERVISOR") };
}

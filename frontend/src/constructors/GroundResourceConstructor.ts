import type { GroundResource } from "../types/GroundResource";

export const createDefaultGroundResource = (overrides: Partial<GroundResource> = {}): GroundResource => ({
  id: 1 as never,
  resource_code: "resource code 1" as never,
  resource_type: "CATERING" as never,
  location: "location 1" as never,
  availability_status: "ON_STAND" as never,
  maintenance_due_at: "2026-06-11T09:00:00Z" as never,
  owner_team: "owner team 1" as never,
  ...overrides
});

export const createGroundResourceForm = createDefaultGroundResource;
export const createGroundResourceResponse = createDefaultGroundResource;

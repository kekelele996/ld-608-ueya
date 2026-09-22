import type { GroundResource } from "../types/GroundResource";

export function createDefaultGroundResource(
  overrides: Partial<GroundResource> = {}
): GroundResource {
  return {
    id: 0,
    resource_code: "",
    resource_type: "CLEANING",
    location: "",
    availability_status: "AVAILABLE",
    maintenance_due_at: null,
    owner_team: "",
    available_from: null,
    available_to: null,
    ...overrides
  };
}

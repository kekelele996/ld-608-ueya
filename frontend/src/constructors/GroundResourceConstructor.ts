import type { GroundResource } from "../types/entities";

export const createEmptyResource = (): GroundResource => ({
  id: 0,
  resource_code: "",
  resource_type: "",
  location: "",
  availability_status: "AVAILABLE",
  maintenance_due_at: null,
  maintenance_end: null,
  owner_team: "",
});

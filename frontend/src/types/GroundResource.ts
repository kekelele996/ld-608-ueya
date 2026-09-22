import type { ResourceStatusValue } from "../constants/ResourceStatus";

export interface GroundResource {
  id: number;
  resource_code: string;
  resource_type: string;
  location: string;
  availability_status: ResourceStatusValue;
  maintenance_due_at: string | null;
  owner_team: string;
  available_from: string | null;
  available_to: string | null;
}

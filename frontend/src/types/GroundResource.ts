export interface GroundResource {
  id: number;
  resource_code: string;
  resource_type: string;
  location: string;
  availability_status: string;
  maintenance_due_at: string;
  owner_team: string;
}

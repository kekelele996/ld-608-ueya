// 本地种子数据快照：全部来自本地数据库，禁止接入第三方 API。
// 仅在后端完全不可用时作为只读兜底（不参与写操作闭环）。
import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { GroundTask } from "../types/GroundTask";
import type { GroundResource } from "../types/GroundResource";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { DelayEvent } from "../types/DelayEvent";

const BASE = "2026-09-21T08:00:00+08:00";

export const mockSeed: {
  turnarounds: FlightTurnaround[];
  tasks: GroundTask[];
  resources: GroundResource[];
  bookings: ResourceBooking[];
  delays: DelayEvent[];
} = {
  turnarounds: [
    {
      id: 1,
      flight_no: "CA1831",
      aircraft_reg: "B-5821",
      stand_no: "203",
      arrival_time: BASE,
      departure_time: "2026-09-21T09:15:00+08:00",
      turnaround_status: "ARRIVING",
      delay_reason: "",
      accumulated_delay: 0,
      plan_generated: false
    }
  ],
  tasks: [],
  resources: [
    {
      id: 1,
      resource_code: "BELT-01",
      resource_type: "BAGGAGE",
      location: "T1-远机位",
      availability_status: "AVAILABLE",
      maintenance_due_at: null,
      owner_team: "TEAM-BAG",
      available_from: null,
      available_to: null
    }
  ],
  bookings: [],
  delays: []
};

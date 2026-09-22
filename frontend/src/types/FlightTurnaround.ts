import type { TurnaroundStatusValue } from "../constants/TurnaroundStatus";

export interface FlightTurnaround {
  id: number;
  flight_no: string;
  aircraft_reg: string;
  stand_no: string;
  arrival_time: string;
  departure_time: string;
  turnaround_status: TurnaroundStatusValue;
  delay_reason: string;
  accumulated_delay: number;
  plan_generated: boolean;
  created_at?: string;
  updated_at?: string;
}

import type { FlightTurnaround } from "../types/FlightTurnaround";

export const createDefaultFlightTurnaround = (overrides: Partial<FlightTurnaround> = {}): FlightTurnaround => ({
  id: 1 as never,
  flight_no: "flight no 1" as never,
  aircraft_reg: "aircraft reg 1" as never,
  stand_no: "stand no 1" as never,
  arrival_time: "2026-06-11T09:00:00Z" as never,
  departure_time: "2026-06-11T09:00:00Z" as never,
  turnaround_status: "ON_STAND" as never,
  delay_reason: "delay reason 1" as never,
  ...overrides
});

export const createFlightTurnaroundForm = createDefaultFlightTurnaround;
export const createFlightTurnaroundResponse = createDefaultFlightTurnaround;

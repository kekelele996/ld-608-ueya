import dayjs from "dayjs";
import type { FlightTurnaround } from "../types/FlightTurnaround";

// 页面/store/service 禁止散写默认结构，统一由构造器产出。
export function createDefaultFlightTurnaround(
  overrides: Partial<FlightTurnaround> = {}
): FlightTurnaround {
  return {
    id: 0,
    flight_no: "",
    aircraft_reg: "",
    stand_no: "",
    arrival_time: dayjs().add(1, "hour").format(),
    departure_time: dayjs().add(2, "hour").format(),
    turnaround_status: "ARRIVING",
    delay_reason: "",
    accumulated_delay: 0,
    plan_generated: false,
    ...overrides
  };
}

export function createTurnaroundForm() {
  return {
    flight_no: "",
    aircraft_reg: "",
    stand_no: "",
    range: [dayjs().add(1, "hour"), dayjs().add(2, "hour").add(15, "minute")] as [
      dayjs.Dayjs,
      dayjs.Dayjs
    ]
  };
}

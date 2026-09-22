import dayjs from "dayjs";
import type { RegisterTurnaroundPayload } from "../api/FlightTurnaround";
import type { RegisterDelayPayload } from "../api/DelayEvent";

// Default turnaround form used by 过站登记 modal.
export const createTurnaroundForm = (overrides: Partial<RegisterTurnaroundPayload> = {}): RegisterTurnaroundPayload => ({
  flight_no: "",
  aircraft_reg: "",
  stand_no: "",
  arrival_time: dayjs().add(5, "minute").toISOString(),
  departure_time: dayjs().add(65, "minute").toISOString(),
  ...overrides,
});

// Default delay registration form.
export const createDelayForm = (overrides: Partial<RegisterDelayPayload> = {}): RegisterDelayPayload => ({
  delay_type: "LATE_ARRIVAL",
  minutes: 30,
  root_cause: "",
  responsibility_team: "",
  ...overrides,
});

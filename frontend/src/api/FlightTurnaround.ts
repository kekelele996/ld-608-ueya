import { mockData } from "../mocks/seedData";
import type { FlightTurnaround } from "../types/FlightTurnaround";

const endpoint = "/api/flight-turnaround";

export async function listFlightTurnaround(): Promise<FlightTurnaround[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.flightTurnaround as unknown as FlightTurnaround[])];
}

export async function saveFlightTurnaround(payload: FlightTurnaround) {
  console.info("save FlightTurnaround", payload);
  return payload;
}

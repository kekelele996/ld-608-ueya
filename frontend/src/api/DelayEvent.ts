import { mockData } from "../mocks/seedData";
import type { DelayEvent } from "../types/DelayEvent";

const endpoint = "/api/delay-event";

export async function listDelayEvent(): Promise<DelayEvent[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.delayEvent as unknown as DelayEvent[])];
}

export async function saveDelayEvent(payload: DelayEvent) {
  console.info("save DelayEvent", payload);
  return payload;
}

import { mockData } from "../mocks/seedData";
import type { ResourceBooking } from "../types/ResourceBooking";

const endpoint = "/api/resource-booking";

export async function listResourceBooking(): Promise<ResourceBooking[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.resourceBooking as unknown as ResourceBooking[])];
}

export async function saveResourceBooking(payload: ResourceBooking) {
  console.info("save ResourceBooking", payload);
  return payload;
}

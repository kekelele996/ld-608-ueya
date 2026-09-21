import { mockData } from "../mocks/seedData";
import type { GroundResource } from "../types/GroundResource";

const endpoint = "/api/ground-resource";

export async function listGroundResource(): Promise<GroundResource[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.groundResource as unknown as GroundResource[])];
}

export async function saveGroundResource(payload: GroundResource) {
  console.info("save GroundResource", payload);
  return payload;
}

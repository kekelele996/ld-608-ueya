import { mockData } from "../mocks/seedData";
import type { GroundTask } from "../types/GroundTask";

const endpoint = "/api/ground-task";

export async function listGroundTask(): Promise<GroundTask[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.groundTask as unknown as GroundTask[])];
}

export async function saveGroundTask(payload: GroundTask) {
  console.info("save GroundTask", payload);
  return payload;
}

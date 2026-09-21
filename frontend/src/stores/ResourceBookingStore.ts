import { create } from "zustand";
import { listResourceBooking } from "../api/ResourceBooking";
import type { ResourceBooking } from "../types/ResourceBooking";

type State = { rows: ResourceBooking[]; loading: boolean; load: () => Promise<void> };

export const useResourceBookingStore = create<State>((set) => ({
  rows: [],
  loading: false,
  async load() {
    set({ loading: true });
    set({ rows: await listResourceBooking(), loading: false });
  }
}));

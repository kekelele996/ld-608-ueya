import { create } from "zustand";
import { listDelayEvent } from "../api/DelayEvent";
import type { DelayEvent } from "../types/DelayEvent";

type State = { rows: DelayEvent[]; loading: boolean; load: () => Promise<void> };

export const useDelayEventStore = create<State>((set) => ({
  rows: [],
  loading: false,
  async load() {
    set({ loading: true });
    set({ rows: await listDelayEvent(), loading: false });
  }
}));

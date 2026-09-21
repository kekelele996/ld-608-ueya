import { create } from "zustand";
import { listFlightTurnaround } from "../api/FlightTurnaround";
import type { FlightTurnaround } from "../types/FlightTurnaround";

type State = { rows: FlightTurnaround[]; loading: boolean; load: () => Promise<void> };

export const useFlightTurnaroundStore = create<State>((set) => ({
  rows: [],
  loading: false,
  async load() {
    set({ loading: true });
    set({ rows: await listFlightTurnaround(), loading: false });
  }
}));

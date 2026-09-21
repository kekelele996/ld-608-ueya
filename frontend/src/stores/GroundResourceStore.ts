import { create } from "zustand";
import { listGroundResource } from "../api/GroundResource";
import type { GroundResource } from "../types/GroundResource";

type State = { rows: GroundResource[]; loading: boolean; load: () => Promise<void> };

export const useGroundResourceStore = create<State>((set) => ({
  rows: [],
  loading: false,
  async load() {
    set({ loading: true });
    set({ rows: await listGroundResource(), loading: false });
  }
}));

import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { Role } from "../types/auth";

export interface AuthState {
  token: string;
  username: string;
  displayName: string;
  role: Role | "";
  teamId: string;
}

const initialState: AuthState = {
  token: localStorage.getItem("ground-turn-token") ?? "",
  username: localStorage.getItem("ground-turn-username") ?? "",
  displayName: localStorage.getItem("ground-turn-display") ?? "",
  role: (localStorage.getItem("ground-turn-role") as Role | "") ?? "",
  teamId: localStorage.getItem("ground-turn-team") ?? "",
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setSession(state, action: PayloadAction<Omit<AuthState, "">>) {
      const p = action.payload;
      state.token = p.token;
      state.username = p.username;
      state.displayName = p.displayName;
      state.role = p.role;
      state.teamId = p.teamId;
      localStorage.setItem("ground-turn-token", p.token);
      localStorage.setItem("ground-turn-username", p.username);
      localStorage.setItem("ground-turn-display", p.displayName);
      localStorage.setItem("ground-turn-role", p.role);
      localStorage.setItem("ground-turn-team", p.teamId);
    },
    logout(state) {
      state.token = "";
      state.username = "";
      state.displayName = "";
      state.role = "";
      state.teamId = "";
      ["ground-turn-token", "ground-turn-username", "ground-turn-display", "ground-turn-role", "ground-turn-team"]
        .forEach((k) => localStorage.removeItem(k));
    },
  },
});

export const { setSession, logout } = authSlice.actions;
export default authSlice.reducer;

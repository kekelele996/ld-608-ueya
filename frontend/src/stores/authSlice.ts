import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import { authApi } from "../api/auth";
import { tokenStore } from "../api/client";
import type { UserInfo } from "../types/api";
import { canRole, type RbacAction } from "../constants/UserRole";

interface AuthState {
  token: string;
  user: UserInfo | null;
}

const initialState: AuthState = {
  token: tokenStore.get(),
  user: null
};

const slice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    loginSucceeded(state, action: PayloadAction<{ token: string; user: UserInfo }>) {
      state.token = action.payload.token;
      state.user = action.payload.user;
      tokenStore.set(action.payload.token);
    },
    logout(state) {
      state.token = "";
      state.user = null;
      tokenStore.clear();
    },
    hydrate(state, action: PayloadAction<UserInfo>) {
      state.user = action.payload;
    }
  }
});

export const { loginSucceeded, logout, hydrate } = slice.actions;

export default slice.reducer;

// store 选择器：组件按钮显隐统一通过 can(role, action)
export const selectCan =
  (action: RbacAction) =>
  (state: { auth: AuthState }): boolean =>
    canRole(state.auth.user?.role, action);

export { authApi };

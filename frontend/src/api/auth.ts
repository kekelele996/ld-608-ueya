import { request, tokenStore } from "./client";
import type { Role } from "../types/auth";

export interface LoginResponse {
  token: string;
  username: string;
  display_name: string;
  role: Role;
  team_id: string;
}

export interface MeResponse {
  username: string;
  role: Role;
  team_id: string;
}

export const login = async (username: string): Promise<LoginResponse> => {
  const res = await request<LoginResponse>("/auth/login", {
    method: "POST",
    body: { username },
  });
  tokenStore.set(res.token);
  return res;
};

export const me = () => request<MeResponse>("/auth/me");

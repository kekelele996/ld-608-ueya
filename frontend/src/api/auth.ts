import { request } from "./client";
import type { LoginResponse } from "../types/api";

export const authApi = {
  login: (username: string, password: string) =>
    request<LoginResponse>("/auth/login", { method: "POST", body: { username, password } }),
  me: () => request<LoginResponse["user"]>("/auth/me")
};

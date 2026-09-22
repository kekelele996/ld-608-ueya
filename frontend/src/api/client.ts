import { message } from "antd";
import { ERROR_CODES } from "../constants/errorCodes";
import { ERROR_MESSAGES } from "../constants/errorMessages";

// 统一响应信封，与后端 middlewares/response.go 一致。
export interface ApiEnvelope<T> {
  ok: boolean;
  data?: T;
  error?: { code: string; message: string };
}

export class ApiError extends Error {
  code: string;
  status: number;
  data?: unknown;
  constructor(code: string, message: string, status: number, data?: unknown) {
    super(message);
    this.code = code;
    this.status = status;
    this.data = data;
  }
}

const TOKEN_KEY = "ground-turn-token";

export const tokenStore = {
  get: () => localStorage.getItem(TOKEN_KEY) ?? "",
  set: (token: string) => localStorage.setItem(TOKEN_KEY, token),
  clear: () => localStorage.removeItem(TOKEN_KEY)
};

// 前端统一请求 /api，禁止硬编码 localhost（nginx 反代到 backend:3000）。
const BASE = "/api";

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: unknown;
  fallback?: unknown; // GET 失败时的本地 mock 回退
  silent?: boolean;
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, fallback, silent } = options;
  let res: Response;
  try {
    res = await fetch(`${BASE}${path}`, {
      method,
      headers: {
        "Content-Type": "application/json",
        ...(tokenStore.get() ? { Authorization: `Bearer ${tokenStore.get()}` } : {})
      },
      body: body !== undefined ? JSON.stringify(body) : undefined
    });
  } catch (networkError) {
    if (fallback !== undefined) return fallback as T;
    throw new ApiError("NETWORK", "无法连接后端服务，请确认服务已启动", 0);
  }

  let payload: ApiEnvelope<T> | null = null;
  try {
    payload = (await res.json()) as ApiEnvelope<T>;
  } catch {
    payload = null;
  }

  if (!res.ok || !payload?.ok) {
    const code = payload?.error?.code ?? "INTERNAL";
    const text = payload?.error?.message ?? ERROR_MESSAGES[code] ?? "请求失败";
    if (fallback !== undefined) return fallback as T;
    if (!silent && code !== ERROR_CODES.AUTH_REQUIRED) {
      message.error(text);
    }
    throw new ApiError(code, text, res.status, payload?.data);
  }
  return payload.data as T;
}

export { ERROR_CODES };

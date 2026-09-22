// apiClient centralizes base URL (/api only, no hardcoded hosts), JWT
// injection and the error envelope used by every entity api module.
import { ERROR_MESSAGES } from "../constants/errorMessages";

export interface ApiError {
  code: string;
  message: string;
  details?: unknown;
  result?: unknown;
  status: number;
}

const TOKEN_KEY = "ground-turn-token";

export const tokenStore = {
  get: () => localStorage.getItem(TOKEN_KEY) ?? "",
  set: (token: string) => localStorage.setItem(TOKEN_KEY, token),
  clear: () => localStorage.removeItem(TOKEN_KEY),
};

export class RequestError extends Error implements ApiError {
  code: string;
  details?: unknown;
  result?: unknown;
  status: number;

  constructor(status: number, code: string, message: string, details?: unknown, result?: unknown) {
    super(message);
    this.status = status;
    this.code = code;
    this.details = details;
    this.result = result;
  }
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  query?: Record<string, string | number | undefined>;
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = "GET", body, query } = options;
  let url = `/api${path}`;
  if (query) {
    const params = new URLSearchParams();
    Object.entries(query).forEach(([key, value]) => {
      if (value !== undefined && value !== "") params.append(key, String(value));
    });
    const qs = params.toString();
    if (qs) url += `?${qs}`;
  }
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  const token = tokenStore.get();
  if (token) headers.Authorization = `Bearer ${token}`;

  const res = await fetch(url, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (res.status === 401) {
    tokenStore.clear();
  }
  const text = await res.text();
  const payload = text ? JSON.parse(text) : {};
  if (!res.ok) {
    throw new RequestError(
      res.status,
      payload.code ?? "INTERNAL_ERROR",
      payload.message ?? ERROR_MESSAGES.INTERNAL_ERROR,
      payload.details,
      payload.result,
    );
  }
  return payload as T;
}

export interface ListEnvelope<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

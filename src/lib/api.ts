import { getInitData } from "./telegram";

// База API. Можно переопределить через VITE_API_URL.
export const API_BASE: string =
  (import.meta as any).env?.VITE_API_URL || "/api";

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const initData = getInitData();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init.headers as Record<string, string> | undefined),
  };
  if (initData) headers["X-Telegram-Init-Data"] = initData;

  const res = await fetch(`${API_BASE}${path}`, { ...init, headers });
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(`API ${res.status}: ${text || res.statusText}`);
  }
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

export const api = {
  authTest: () => request<{ user: any; token: string }>("/auth/test", { method: "POST" }),
  getProfile: () => request<any>("/profile"),
  getTasks: () => request<{ tasks: any[] }>("/tasks"),
  createTask: (data: { title: string; description: string; task_type: string }) =>
    request<any>("/tasks", { method: "POST", body: JSON.stringify(data) }),
  startTask: (id: number) =>
    request<any>(`/tasks/${id}/start`, { method: "POST" }),
  completeTask: (id: number) =>
    request<any>(`/tasks/${id}/complete`, { method: "POST" }),
};
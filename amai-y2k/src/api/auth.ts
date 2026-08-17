import { api } from "./client";
import type { ApiMessage } from "./client";

export async function login(
  username: string,
  password: string
): Promise<ApiMessage> {
  const { data } = await api.post<ApiMessage>("/auth/login", null, {
    auth: { username, password },
  });
  return data;
}

export async function logout(): Promise<ApiMessage> {
  const { data } = await api.post<ApiMessage>("/auth/logout");
  return data;
}
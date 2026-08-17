import axios, { AxiosError } from "axios";
import type { AxiosInstance } from "axios";

interface ViteEnv {
  VITE_API_BASE_URL?: string;
}
interface ViteImportMeta {
  env?: ViteEnv;
}

const BASE_URL: string =
  (import.meta as unknown as ViteImportMeta).env?.VITE_API_BASE_URL ?? "/api";

export const api: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  withCredentials: true,
});

export interface ApiMessage {
  message: string;
}

interface ApiErrorShape {
  message: string;
}

export function getApiErrorMessage(
  error: unknown,
  fallback = "Something went wrong"
): string {
  if (axios.isAxiosError(error)) {
    const err = error as AxiosError<ApiErrorShape>;
    return err.response?.data?.message ?? fallback;
  }
  return fallback;
}

export { BASE_URL };
import { defineStore } from "pinia";
import { ref } from "vue";
import { login as apiLogin, logout as apiLogout } from "@/api/auth";
import { getApiErrorMessage, api } from "@/api/client";
import router from "@/router";

export const useAuthStore = defineStore("auth", () => {
  const isAuthenticated = ref(false);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  async function login(username: string, password: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      await apiLogin(username, password);
      isAuthenticated.value = true;
      return true;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to log in");
      isAuthenticated.value = false;
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  async function logout(): Promise<void> {
    isLoading.value = true;
    error.value = null;
    try {
      await apiLogout();
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to log out");
    } finally {
      isAuthenticated.value = false;
      isLoading.value = false;
      router.push({ name: "home" });
    }
  }

  async function checkAuth(): Promise<void> {
    try {
      await api.get("/admin/whoami");
      isAuthenticated.value = true;
    } catch {
      isAuthenticated.value = false;
    }
  }

  function forceLogout(): void {
    isAuthenticated.value = false;
    router.push({ name: "home" });
  }

  return { isAuthenticated, isLoading, error, login, logout, checkAuth, forceLogout };
});
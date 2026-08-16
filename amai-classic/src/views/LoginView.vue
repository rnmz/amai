<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const validationError = ref('')

const errorMessage = computed(() => validationError.value || authStore.error || '')

async function onSubmit() {
  validationError.value = ''

  if (!username.value || !password.value) {
    validationError.value = t('login.empty')
    return
  }

  const success = await authStore.login(username.value, password.value)
  if (success) {
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  }
}
</script>

<template>
  <div class="login-page">
    <form class="login-card" @submit.prevent="onSubmit">
      <h1 class="login-title">Amai</h1>

      <label class="field">
        <span class="field-label">{{ t('login.login') }}</span>
        <input
          v-model="username"
          type="text"
          autocomplete="username"
          placeholder="admin"
          :disabled="authStore.isLoading"
        />
      </label>

      <label class="field">
        <span class="field-label">{{ t('login.password') }}</span>
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="••••••••"
          :disabled="authStore.isLoading"
        />
      </label>

      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>

      <button type="submit" class="submit-btn" :disabled="authStore.isLoading">
        {{ authStore.isLoading ? t('common.loading') : t('login.enter') }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 90vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 380px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 32px;
  background-color: #fff;
  border: 1px solid #d0d7de;
  border-radius: 14px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.login-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1f2328;
  text-align: center;
  margin: 0 0 8px 0;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 0.85rem;
  font-weight: 500;
  color: #57606a;
}

.field input {
  padding: 10px 12px;
  border: 1px solid #d0d7de;
  border-radius: 8px;
  font-size: 0.95rem;
  color: #1f2328;
  background-color: #fff;
  transition:
    border-color 0.15s ease,
    box-shadow 0.15s ease;
}

.field input:focus {
  outline: none;
  border-color: #0b7a30;
  box-shadow: 0 0 0 3px rgba(11, 122, 48, 0.12);
}

.field input:disabled {
  background-color: #f6f8fa;
  cursor: not-allowed;
}

.error-message {
  margin: 0;
  padding: 8px 12px;
  border-radius: 8px;
  background-color: rgba(192, 57, 43, 0.08);
  color: #c0392b;
  font-size: 0.85rem;
}

.submit-btn {
  margin-top: 4px;
  padding: 11px;
  border: none;
  border-radius: 8px;
  background-color: #0b7a30;
  color: #fff;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    opacity 0.15s ease;
}

.submit-btn:hover:not(:disabled) {
  background-color: #095f26;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
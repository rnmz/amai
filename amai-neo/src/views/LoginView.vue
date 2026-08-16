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
  min-height: 80vh;
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
  gap: 18px;
  padding: 36px 32px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 18px;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.login-title {
  font-family: var(--main-font);
  font-size: 1.6rem;
  font-weight: 700;
  color: var(--color-text);
  text-align: center;
  margin: 0 0 4px 0;
  letter-spacing: -0.01em;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-family: var(--main-font);
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--md-text-secondary);
}

.field input {
  padding: 12px 14px;
  border: 1px solid var(--md-border);
  border-radius: 10px;
  font-family: var(--main-font);
  font-size: 0.95rem;
  color: var(--color-text);
  background-color: var(--color-surface);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.field input:focus {
  outline: none;
  border-color: var(--md-accent);
  box-shadow: 0 0 0 3px var(--md-selection-bg);
}

.field input:disabled {
  background-color: var(--md-bg-pre);
  color: var(--md-text-secondary);
  cursor: not-allowed;
  opacity: 0.7;
}

.error-message {
  margin: 0;
  padding: 10px 14px;
  border-radius: 8px;
  background-color: var(--md-mark-red-bg);
  color: var(--md-num-negative);
  border: 1px solid rgba(239, 68, 68, 0.3);
  font-size: 0.85rem;
  font-weight: 500;
  text-align: center;
}

.submit-btn {
  margin-top: 6px;
  padding: 12px;
  border: 1px solid var(--md-accent);
  border-radius: 10px;
  background-color: var(--md-accent);
  color: #ffffff;
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.submit-btn:hover:not(:disabled) {
  background-color: transparent;
  color: var(--md-accent);
  border-color: var(--md-accent);
  transform: translateY(-1px);
}

.submit-btn:active:not(:disabled) {
  transform: scale(0.97);
  background-color: var(--md-bg-blockquote);
}

.submit-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  border-color: var(--md-border);
  background-color: var(--md-bg-pre);
  color: var(--md-text-secondary);
}
</style>
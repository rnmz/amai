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
      <h1 class="login-title">AMAI</h1>

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
  min-height: 70vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 360px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 28px 24px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  border-radius: 8px;
  box-shadow: 1px 1px 3px rgba(0, 0, 0, 0.08);
}

.login-title {
  font-family: var(--main-font);
  font-size: 1.6rem;
  font-weight: 800;
  color: var(--md-text);
  text-align: center;
  margin: 0 0 4px 0;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--md-text);
}

.field input {
  padding: 10px 12px;
  border: 1px solid var(--md-border);
  border-radius: 4px;
  background-color: var(--color-surface);
  font-family: var(--main-font);
  font-size: 0.9rem;
  color: var(--md-text);
  transition: all 0.15s ease;
  box-sizing: border-box;
}

.field input:focus {
  outline: none;
  border-color: var(--md-accent);
  box-shadow: 0 0 0 2px var(--md-bg-table-row-hover);
}

.field input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-message {
  margin: 0;
  padding: 8px 12px;
  background-color: var(--md-bg-disclaimer);
  color: var(--md-num-negative);
  border: 1px solid var(--md-num-negative);
  border-radius: 4px;
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 600;
  text-align: center;
}

.submit-btn {
  margin-top: 4px;
  padding: 10px;
  border: 1px solid var(--md-accent);
  border-radius: 4px;
  background-color: var(--md-accent);
  color: #ffffff;
  font-family: var(--main-font);
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 1px 1px 2px rgba(0, 0, 0, 0.08);
}

.submit-btn:hover:not(:disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.submit-btn:active:not(:disabled) {
  transform: scale(0.97);
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  border-color: var(--md-border);
  background-color: var(--md-bg-code);
  color: var(--md-text-secondary);
  box-shadow: none;
}
</style>
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
      <h1 class="login-title" data-text="AMAI_SYS">AMAI_SYS</h1>

      <label class="field">
        <span class="field-label">// {{ t('login.login') }}</span>
        <input
          v-model="username"
          type="text"
          autocomplete="username"
          placeholder="admin"
          :disabled="authStore.isLoading"
        />
      </label>

      <label class="field">
        <span class="field-label">// {{ t('login.password') }}</span>
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="••••••••"
          :disabled="authStore.isLoading"
        />
      </label>

      <p v-if="errorMessage" class="error-message">[SYS_ERROR] :: {{ errorMessage }}</p>

      <button type="submit" class="submit-btn" :disabled="authStore.isLoading">
        [{{ authStore.isLoading ? t('common.loading') : t('login.enter') }}]
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
  gap: 20px;
  padding: 36px 32px;
  background-color: var(--color-surface);
  border: 1px solid var(--md-border);
  clip-path: polygon(
    0 0,
    calc(100% - 15px) 0,
    100% 15px,
    100% 100%,
    15px 100%,
    0 calc(100% - 15px)
  );
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.15);
}

.login-title {
  font-family: var(--main-font);
  font-size: 1.8rem;
  font-weight: 900;
  color: #ffffff;
  text-align: center;
  margin: 0 0 8px 0;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  text-shadow: 0 0 10px var(--md-accent);
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--md-accent);
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.field input {
  padding: 12px 14px;
  border: 1px solid var(--md-border-light);
  background-color: var(--md-bg-pre);
  font-family: 'Orbitron', monospace;
  font-size: 0.9rem;
  color: var(--md-text);
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%);
  transition: all 0.2s ease;
}

.field input:focus {
  outline: none;
  border-color: var(--md-accent);
  box-shadow: 0 0 10px rgba(0, 240, 255, 0.2);
}

.field input:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.error-message {
  margin: 0;
  padding: 10px 14px;
  background-color: var(--md-mark-red-bg);
  color: var(--md-num-negative);
  border: 1px solid var(--md-num-negative);
  font-family: var(--main-font);
  font-size: 0.85rem;
  font-weight: 700;
  text-align: center;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%);
}

.submit-btn {
  margin-top: 8px;
  padding: 12px;
  border: 1px solid var(--md-accent);
  background-color: var(--md-accent);
  color: #050811;
  font-family: var(--main-font);
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  cursor: pointer;
  clip-path: polygon(0 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%);
  box-shadow: 0 0 12px var(--md-accent);
  transition: all 0.2s ease;
}

.submit-btn:hover:not(:disabled) {
  background-color: #00f0ff;
  border-color: #00f0ff;
  box-shadow: 0 0 18px #00f0ff;
  transform: translateY(-1px);
}

.submit-btn:active:not(:disabled) {
  transform: scale(0.97);
}

.submit-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  border-color: var(--md-border-light);
  background-color: var(--color-surface);
  color: var(--md-text-secondary);
  box-shadow: none;
}
</style>
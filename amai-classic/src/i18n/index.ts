import { createI18n } from 'vue-i18n'
import { watch } from 'vue'
import ru from './locales/ru.json'
import en from './locales/en.json'
import jp from './locales/jp.json'

const savedLang = localStorage.getItem('user-lang') || 'en-US'

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'en-US',
  messages: {
    'en-US': en,
    'ru-RU': ru,
    'ja-JP': jp,
  },
  globalInjection: true,
})

watch(
  i18n.global.locale,
  (newLocale) => {
    document.documentElement.lang = newLocale
    localStorage.setItem('user-lang', newLocale)
  },
  { immediate: true },
)

export default i18n

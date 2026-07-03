import { createI18n, type I18n } from "vue-i18n";
import ru from './locales/ru.json'
import en from './locales/en.json'

const i18n = createI18n({
    legacy: false,
    locale: localStorage.getItem('user-lang') || 'en-US',
    fallbackLocale: 'en-US',
    messages: {
        'en-US': en,
        'ru-RU': ru
    },
    globalInjection: true
})

export default i18n
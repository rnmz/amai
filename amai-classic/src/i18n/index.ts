import { createI18n } from "vue-i18n";
import ru from './locales/ru.json'
import en from './locales/en.json'
import jp from './locales/jp.json'

const i18n = createI18n({
    legacy: false,
    locale: localStorage.getItem('user-lang') || 'en-US',
    fallbackLocale: 'en-US',
    messages: {
        'en-US': en,
        'ru-RU': ru,
        'ja-JP': jp
    },
    globalInjection: true
})

export default i18n

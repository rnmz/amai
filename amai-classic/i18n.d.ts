import { DefineLocaleMessage } from 'vue-i18n'
import en from './src/i18n/locales/en.json'

declare module 'vue-i18n' {
    export interface DefineLocaleMessage extends typeof en {}
}
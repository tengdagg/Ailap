import { createI18n } from 'vue-i18n'
import en from './en-US'
import zh from './zh-CN'

const i18n = createI18n({
    legacy: false, // use Composition API
    globalInjection: true,
    locale: localStorage.getItem('locale') || 'zh-CN',
    fallbackLocale: 'zh-CN',
    messages: {
        'zh-CN': zh,
        'en-US': en
    }
})

export default i18n

import { ref } from 'vue'

const THEME_STORAGE_KEY = 'theme'
const systemTheme = window.matchMedia('(prefers-color-scheme: dark)')

export const isDark = ref(false)

const applyTheme = (dark: boolean) => {
    isDark.value = dark
    document.documentElement.dataset.theme = dark ? 'dark' : 'light'
}

export const initializeTheme = () => {
    const storedTheme = localStorage.getItem(THEME_STORAGE_KEY)
    applyTheme(storedTheme ? storedTheme === 'dark' : systemTheme.matches)

    if (!storedTheme) {
        systemTheme.addEventListener('change', event => applyTheme(event.matches))
    }
}

export const toggleTheme = () => {
    const nextIsDark = !isDark.value
    applyTheme(nextIsDark)
    localStorage.setItem(THEME_STORAGE_KEY, nextIsDark ? 'dark' : 'light')
}

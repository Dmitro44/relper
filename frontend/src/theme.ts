import { ref } from 'vue'

export type Theme = 'light' | 'dark'

const STORAGE_KEY = 'relper-theme'
const media = window.matchMedia('(prefers-color-scheme: dark)')
const theme = ref<Theme>('light')

function storedTheme(): Theme | null {
  try {
    const value = localStorage.getItem(STORAGE_KEY)
    return value === 'light' || value === 'dark' ? value : null
  } catch {
    return null
  }
}

function systemTheme(): Theme {
  return media.matches ? 'dark' : 'light'
}

function apply(value: Theme): void {
  theme.value = value
  document.documentElement.dataset.theme = value
}

function onSystemChange(event: MediaQueryListEvent): void {
  if (!storedTheme()) {
    apply(event.matches ? 'dark' : 'light')
  }
}

function persistTheme(value: Theme): void {
  try {
    localStorage.setItem(STORAGE_KEY, value)
  } catch {
    return
  }
}

export function initTheme(): void {
  apply(storedTheme() ?? systemTheme())
  media.addEventListener('change', onSystemChange)
}

export function toggleTheme(): void {
  const next = theme.value === 'dark' ? 'light' : 'dark'
  apply(next)
  persistTheme(next)
}

export { theme }

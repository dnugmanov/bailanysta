import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { applyTheme, getSystemTheme } from '@/lib/utils'

type Theme = 'light' | 'dark' | 'system'

interface ThemeState {
  theme: Theme
  actualTheme: 'light' | 'dark'
  setTheme: (theme: Theme) => void
  toggleTheme: () => void
  initialize: () => void
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      theme: 'system',
      actualTheme: 'light',

      setTheme: (theme: Theme) => {
        const actualTheme = theme === 'system' ? getSystemTheme() : theme
        applyTheme(theme)
        set({ theme, actualTheme })
      },

      toggleTheme: () => {
        const { actualTheme } = get()
        const newTheme = actualTheme === 'light' ? 'dark' : 'light'
        const { setTheme } = get()
        setTheme(newTheme)
      },

      initialize: () => {
        const { theme, setTheme } = get()
        setTheme(theme)

        // Listen for system theme changes
        if (typeof window !== 'undefined') {
          const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
          const handleChange = () => {
            const { theme, setTheme } = get()
            if (theme === 'system') {
              setTheme('system') // This will recalculate the actual theme
            }
          }

          mediaQuery.addEventListener('change', handleChange)
          
          // Cleanup function (though this won't be called in Zustand)
          return () => mediaQuery.removeEventListener('change', handleChange)
        }
      },
    }),
    {
      name: 'theme-storage',
      partialize: (state) => ({ theme: state.theme }),
    }
  )
)
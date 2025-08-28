import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '@/types/api'
import { apiClient } from '@/api/client'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  setUser: (user: User | null) => void
  initialize: () => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,

      login: async (email: string, password: string) => {
        set({ isLoading: true })
        try {
          const response = await apiClient.login({ email, password })
          console.log('✅ Login successful:', response.user.username)
          console.log('✅ Token received:', response.tokens.access_token?.substring(0, 20) + '...')
          
          // Ensure token is set properly
          if (response.tokens.access_token) {
            apiClient.setAccessToken(response.tokens.access_token)
          }
          
          set({
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      register: async (email: string, username: string, password: string) => {
        set({ isLoading: true })
        try {
          const response = await apiClient.register({ email, username, password })
          console.log('✅ Registration successful:', response.user.username)
          console.log('✅ Token received:', response.tokens.access_token?.substring(0, 20) + '...')
          
          // Ensure token is set properly
          if (response.tokens.access_token) {
            apiClient.setAccessToken(response.tokens.access_token)
          }
          
          set({
            user: response.user,
            isAuthenticated: true,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      logout: async () => {
        set({ isLoading: true })
        try {
          await apiClient.logout()
        } catch (error) {
          console.error('Logout error:', error)
        } finally {
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
          })
        }
      },

      setUser: (user: User | null) => {
        set({
          user,
          isAuthenticated: !!user,
        })
      },

      initialize: async () => {
        const token = localStorage.getItem('access_token')
        console.log('🔄 Auth initialize - token found:', !!token)
        
        if (!token) {
          set({ isAuthenticated: false, user: null })
          return
        }

        // Set the token in API client
        apiClient.setAccessToken(token)
        
        set({ isLoading: true })
        try {
          const user = await apiClient.getCurrentUser()
          set({
            user,
            isAuthenticated: true,
            isLoading: false,
          })
        } catch (error) {
          // Token is invalid, clear it
          localStorage.removeItem('access_token')
          apiClient.setAccessToken(null)
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
          })
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
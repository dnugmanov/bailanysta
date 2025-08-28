import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { Toaster } from '@/components/ui/toaster'
import Layout from '@/components/layout/Layout'
import LoginPage from '@/pages/Login'
import RegisterPage from '@/pages/Register'
import FeedPage from '@/pages/Feed'
import ProfilePage from '@/pages/Profile'
import ComposePage from '@/pages/Compose'
import SearchPage from '@/pages/Search'
import PostPage from '@/pages/Post'
import NotificationsPage from '@/pages/Notifications'
import CoursesPage from '@/pages/Courses'
import CommunityPage from '@/pages/Community'
import TrendingPage from '@/pages/Trending'

import ProtectedRoute from '@/components/ProtectedRoute'
import { QueryErrorResetBoundary } from '@tanstack/react-query'
import { ErrorBoundary } from 'react-error-boundary'

function ErrorFallback({ error, resetErrorBoundary }: { error: Error; resetErrorBoundary: () => void }) {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center space-y-4">
        <h2 className="text-2xl font-bold text-destructive">Что-то пошло не так</h2>
        <p className="text-muted-foreground">{error.message}</p>
        <button
          onClick={resetErrorBoundary}
          className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
        >
          Попробовать снова
        </button>
      </div>
    </div>
  )
}

function App() {
  const { initialize, isLoading } = useAuthStore()
  const { initialize: initializeTheme } = useThemeStore()

  useEffect(() => {
    initialize()
    initializeTheme()
  }, [initialize, initializeTheme])

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    )
  }

  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary
          FallbackComponent={ErrorFallback}
          onReset={reset}
        >
          <div className="min-h-screen bg-background text-foreground">
            <Routes>
              {/* Public routes */}
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />
              
              {/* Protected routes */}
              <Route element={<ProtectedRoute />}>
                <Route element={<Layout />}>
                  <Route path="/" element={<Navigate to="/feed" replace />} />
                  <Route path="/feed" element={<FeedPage />} />
                  <Route path="/compose" element={<ComposePage />} />
                  <Route path="/search" element={<SearchPage />} />
                  <Route path="/notifications" element={<NotificationsPage />} />
                  <Route path="/profile/:username" element={<ProfilePage />} />
                  <Route path="/post/:id" element={<PostPage />} />
                  <Route path="/courses" element={<CoursesPage />} />
                  <Route path="/community" element={<CommunityPage />} />
                  <Route path="/trending" element={<TrendingPage />} />

                </Route>
              </Route>
              
              {/* Catch all */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
            <Toaster />
          </div>
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  )
}

export default App
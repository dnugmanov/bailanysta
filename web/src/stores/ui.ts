import { create } from 'zustand'

interface UIState {
  sidebarOpen: boolean
  composerOpen: boolean
  searchOpen: boolean
  notificationsOpen: boolean
  mobileMenuOpen: boolean
  
  // Actions
  setSidebarOpen: (open: boolean) => void
  setComposerOpen: (open: boolean) => void
  setSearchOpen: (open: boolean) => void
  setNotificationsOpen: (open: boolean) => void
  setMobileMenuOpen: (open: boolean) => void
  toggleSidebar: () => void
  toggleComposer: () => void
  toggleSearch: () => void
  toggleNotifications: () => void
  toggleMobileMenu: () => void
  closeAll: () => void
}

export const useUIStore = create<UIState>((set, get) => ({
  sidebarOpen: true,
  composerOpen: false,
  searchOpen: false,
  notificationsOpen: false,
  mobileMenuOpen: false,

  setSidebarOpen: (open: boolean) => set({ sidebarOpen: open }),
  setComposerOpen: (open: boolean) => set({ composerOpen: open }),
  setSearchOpen: (open: boolean) => set({ searchOpen: open }),
  setNotificationsOpen: (open: boolean) => set({ notificationsOpen: open }),
  setMobileMenuOpen: (open: boolean) => set({ mobileMenuOpen: open }),

  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  toggleComposer: () => set((state) => ({ composerOpen: !state.composerOpen })),
  toggleSearch: () => set((state) => ({ searchOpen: !state.searchOpen })),
  toggleNotifications: () => set((state) => ({ notificationsOpen: !state.notificationsOpen })),
  toggleMobileMenu: () => set((state) => ({ mobileMenuOpen: !state.mobileMenuOpen })),

  closeAll: () => set({
    composerOpen: false,
    searchOpen: false,
    notificationsOpen: false,
    mobileMenuOpen: false,
  }),
}))

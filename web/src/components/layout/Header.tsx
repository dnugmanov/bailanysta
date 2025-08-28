import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Bell,
  Menu,
  Search,
  Settings,
  LogOut,
  User,
  Sparkles,
  Moon,
  Sun,
  Monitor,
  Plus,
} from 'lucide-react'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { useUIStore } from '@/stores/ui'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import * as DropdownMenu from '@radix-ui/react-dropdown-menu'
import * as Popover from '@radix-ui/react-popover'
import { getInitials } from '@/lib/utils'
import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/api/client'

export default function Header() {
  const [searchOpen, setSearchOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const navigate = useNavigate()
  
  const { user, logout } = useAuthStore()
  const { theme, setTheme } = useThemeStore()
  const { toggleSidebar, toggleMobileMenu, setComposerOpen } = useUIStore()

  // Fetch unread notifications count
  const { data: unreadCount } = useQuery({
    queryKey: ['notifications', 'unread-count'],
    queryFn: () => apiClient.getUnreadCount(),
    refetchInterval: 30000, // Refetch every 30 seconds
  })

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery)}`)
      setSearchOpen(false)
      setSearchQuery('')
    }
  }

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const themeIcons = {
    light: Sun,
    dark: Moon,
    system: Monitor,
  }
  const themeIcon = themeIcons[theme as keyof typeof themeIcons] || Monitor

  const ThemeIcon = themeIcon

  return (
    <header className="fixed top-0 left-0 right-0 z-50 glass-effect border-b">
      <div className="flex items-center justify-between px-4 py-3">
        {/* Left section */}
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => {
              if (window.innerWidth >= 1024) {
                toggleSidebar()
              } else {
                toggleMobileMenu()
              }
            }}
            className="lg:inline-flex"
          >
            <Menu className="h-5 w-5" />
          </Button>

          <Link to="/feed" className="flex items-center gap-2">
            <motion.div
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              className="flex items-center justify-center w-8 h-8 bg-primary rounded-lg"
            >
              <Sparkles className="w-5 h-5 text-primary-foreground" />
            </motion.div>
            <span className="hidden sm:block text-xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              Bailanysta
            </span>
          </Link>
        </div>

        {/* Center section - Search */}
        <div className="flex-1 max-w-md mx-4">
          <Popover.Root open={searchOpen} onOpenChange={setSearchOpen}>
            <Popover.Trigger asChild>
              <Button
                variant="outline"
                className="w-full justify-start text-muted-foreground"
              >
                <Search className="h-4 w-4 mr-2" />
                <span className="hidden sm:inline">Поиск постов и пользователей...</span>
                <span className="sm:hidden">Поиск...</span>
              </Button>
            </Popover.Trigger>
            <Popover.Portal>
              <Popover.Content
                className="w-96 p-4 bg-popover border rounded-lg shadow-lg z-50"
                sideOffset={5}
              >
                <form onSubmit={handleSearch} className="space-y-3">
                  <Input
                    placeholder="Поиск постов, пользователей, хештегов..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    autoFocus
                  />
                  <div className="flex justify-end">
                    <Button type="submit" size="sm">
                      Найти
                    </Button>
                  </div>
                </form>
              </Popover.Content>
            </Popover.Portal>
          </Popover.Root>
        </div>

        {/* Right section */}
        <div className="flex items-center gap-2">
          {/* Compose button */}
          <Link to="/compose">
            <Button
              size="sm"
              className="hidden sm:flex"
            >
              <Plus className="h-4 w-4 mr-1" />
              Создать
            </Button>
          </Link>
          
          <Link to="/compose">
            <Button
              size="icon"
              className="sm:hidden"
            >
              <Plus className="h-4 w-4" />
            </Button>
          </Link>

          {/* Notifications */}
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate('/notifications')}
            className="relative"
          >
            <Bell className="h-5 w-5" />
            {unreadCount && unreadCount.count > 0 && (
              <motion.span
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                className="absolute -top-1 -right-1 bg-destructive text-destructive-foreground text-xs rounded-full h-5 w-5 flex items-center justify-center"
              >
                {unreadCount.count > 9 ? '9+' : unreadCount.count}
              </motion.span>
            )}
          </Button>

          {/* Theme toggle */}
          <DropdownMenu.Root>
            <DropdownMenu.Trigger asChild>
              <Button variant="ghost" size="icon">
                <ThemeIcon className="h-5 w-5" />
              </Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Portal>
              <DropdownMenu.Content
                className="min-w-[8rem] bg-popover border rounded-md p-1 shadow-md z-50"
                sideOffset={5}
              >
                <DropdownMenu.Item
                  className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm"
                  onClick={() => setTheme('light')}
                >
                  <Sun className="h-4 w-4" />
                  Светлая
                </DropdownMenu.Item>
                <DropdownMenu.Item
                  className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm"
                  onClick={() => setTheme('dark')}
                >
                  <Moon className="h-4 w-4" />
                  Темная
                </DropdownMenu.Item>
                <DropdownMenu.Item
                  className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm"
                  onClick={() => setTheme('system')}
                >
                  <Monitor className="h-4 w-4" />
                  Системная
                </DropdownMenu.Item>
              </DropdownMenu.Content>
            </DropdownMenu.Portal>
          </DropdownMenu.Root>

          {/* User menu */}
          <DropdownMenu.Root>
            <DropdownMenu.Trigger asChild>
              <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                <Avatar className="h-8 w-8">
                  <AvatarImage src={user?.avatar_url || undefined} alt={user?.username} />
                  <AvatarFallback className="text-xs">
                    {user ? getInitials(user.username) : 'U'}
                  </AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenu.Trigger>
            <DropdownMenu.Portal>
              <DropdownMenu.Content
                className="w-56 bg-popover border rounded-md p-1 shadow-md z-50"
                sideOffset={5}
                align="end"
              >
                <div className="flex items-center justify-start gap-2 p-2">
                  <Avatar className="h-8 w-8">
                    <AvatarImage src={user?.avatar_url || undefined} alt={user?.username} />
                    <AvatarFallback className="text-xs">
                      {user ? getInitials(user.username) : 'U'}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col space-y-1 leading-none">
                    <p className="font-medium text-sm">{user?.username}</p>
                    <p className="text-xs text-muted-foreground">{user?.email}</p>
                  </div>
                </div>
                <DropdownMenu.Separator className="my-1 h-px bg-border" />
                
                <DropdownMenu.Item
                  className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm"
                  onClick={() => navigate(`/profile/${user?.username}`)}
                >
                  <User className="h-4 w-4" />
                  Профиль
                </DropdownMenu.Item>
                
                <DropdownMenu.Item
                  className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm"
                  onClick={() => navigate('/settings')}
                >
                  <Settings className="h-4 w-4" />
                  Настройки
                </DropdownMenu.Item>
                
                <DropdownMenu.Separator className="my-1 h-px bg-border" />
                
                <DropdownMenu.Item
                  className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm text-destructive"
                  onClick={handleLogout}
                >
                  <LogOut className="h-4 w-4" />
                  Выйти
                </DropdownMenu.Item>
              </DropdownMenu.Content>
            </DropdownMenu.Portal>
          </DropdownMenu.Root>
        </div>
      </div>
    </header>
  )
}

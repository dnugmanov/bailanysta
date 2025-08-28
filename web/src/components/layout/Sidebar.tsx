import { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Home,
  Search,
  Bell,
  BookOpen,
  Users,
  TrendingUp,
  Sparkles,
} from 'lucide-react'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import AIAssistant from '@/components/features/AIAssistant'

interface NavItem {
  icon: React.ComponentType<{ className?: string }>
  label: string
  path: string
  badge?: number
}

export default function Sidebar() {
  const location = useLocation()
  const { user } = useAuthStore()
  const { sidebarOpen } = useUIStore()
  const [aiAssistantOpen, setAiAssistantOpen] = useState(false)

  const navItems: NavItem[] = [
    { icon: Home, label: 'Лента', path: '/feed' },
    { icon: Search, label: 'Поиск', path: '/search' },
    { icon: Bell, label: 'Уведомления', path: '/notifications' },
    { icon: BookOpen, label: 'Курсы', path: '/courses' },
    { icon: Users, label: 'Сообщество', path: '/community' },
    { icon: TrendingUp, label: 'Популярное', path: '/trending' },
  ]

  const secondaryItems: NavItem[] = []

  const NavLink = ({ item }: { item: NavItem }) => {
    const isActive = location.pathname === item.path
    const Icon = item.icon

    return (
      <Link to={item.path}>
        <motion.div
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          className={cn(
            "flex items-center gap-3 px-3 py-2 rounded-lg transition-colors relative",
            isActive
              ? "bg-primary text-primary-foreground"
              : "hover:bg-accent hover:text-accent-foreground"
          )}
        >
          <Icon className="h-5 w-5 flex-shrink-0" />
          {sidebarOpen && (
            <motion.span
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -10 }}
              className="font-medium truncate"
            >
              {item.label}
            </motion.span>
          )}
          {item.badge && item.badge > 0 && sidebarOpen && (
            <motion.span
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              className="ml-auto bg-destructive text-destructive-foreground text-xs rounded-full h-5 w-5 flex items-center justify-center"
            >
              {item.badge > 9 ? '9+' : item.badge}
            </motion.span>
          )}
          
          {isActive && (
            <motion.div
              layoutId="activeTab"
              className="absolute inset-0 bg-primary rounded-lg -z-10"
              transition={{ type: "spring", stiffness: 300, damping: 30 }}
            />
          )}
        </motion.div>
      </Link>
    )
  }

  return (
    <div className="h-full bg-card border-r border-border flex flex-col">
      <div className="flex-1 p-4 space-y-2">
        <nav className="space-y-1">
          {navItems.map((item) => (
            <NavLink key={item.path} item={item} />
          ))}
        </nav>

        {sidebarOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="py-4"
          >
            <div className="px-3 py-2">
              <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                Быстрые действия
              </h3>
            </div>
            <div className="space-y-2">
              <Link to="/compose" className="w-full">
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full justify-start"
                >
                  <motion.div
                    whileHover={{ rotate: 90 }}
                    transition={{ duration: 0.2 }}
                    className="mr-2"
                  >
                    ✨
                  </motion.div>
                  Создать пост
                </Button>
              </Link>

              <Button
                variant="outline"
                size="sm"
                className="w-full justify-start"
                onClick={() => setAiAssistantOpen(true)}
              >
                <motion.div
                  whileHover={{ rotate: 360 }}
                  transition={{ duration: 0.3 }}
                  className="mr-2"
                >
                  <Sparkles className="h-4 w-4" />
                </motion.div>
                AI Помощник
              </Button>

            </div>
          </motion.div>
        )}
      </div>

      {/* Bottom section */}
      <div className="border-t border-border p-4 space-y-1">
        {secondaryItems.map((item) => (
          <NavLink key={item.path} item={item} />
        ))}
        
        {sidebarOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="pt-4 text-xs text-muted-foreground"
          >
            <p>© 2025 Bailanysta</p>
            <p>Образовательная платформа</p>
          </motion.div>
        )}
      </div>

      {/* AI Assistant Modal */}
      <AIAssistant
        isOpen={aiAssistantOpen}
        onClose={() => setAiAssistantOpen(false)}
      />
    </div>
  )
}

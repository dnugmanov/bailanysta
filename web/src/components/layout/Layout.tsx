import { Outlet } from 'react-router-dom'
import { motion } from 'framer-motion'
import Header from './Header'
import Sidebar from './Sidebar'
import { useUIStore } from '@/stores/ui'
import { cn } from '@/lib/utils'

export default function Layout() {
  const { sidebarOpen, mobileMenuOpen } = useUIStore()

  return (
    <div className="min-h-screen bg-background">
      <Header />
      
      <div className="flex">
        {/* Desktop Sidebar */}
        <aside className={cn(
          "hidden lg:flex flex-col fixed left-0 top-16 h-[calc(100vh-4rem)] transition-all duration-300 ease-in-out z-30",
          sidebarOpen ? "w-64" : "w-16"
        )}>
          <Sidebar />
        </aside>

        {/* Mobile Sidebar Overlay */}
        {mobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 z-40 lg:hidden"
            onClick={() => useUIStore.getState().setMobileMenuOpen(false)}
          />
        )}

        {/* Mobile Sidebar */}
        <motion.aside
          initial={{ x: "-100%" }}
          animate={{ x: mobileMenuOpen ? 0 : "-100%" }}
          transition={{ type: "spring", stiffness: 300, damping: 30 }}
          className="fixed left-0 top-16 h-[calc(100vh-4rem)] w-64 z-50 lg:hidden"
        >
          <Sidebar />
        </motion.aside>

        {/* Main Content */}
        <main className={cn(
          "flex-1 transition-all duration-300 ease-in-out",
          "lg:ml-16", // Base margin for collapsed sidebar
          sidebarOpen && "lg:ml-64" // Expanded margin for open sidebar
        )}>
          <div className="pt-16"> {/* Account for header height */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3 }}
              className="container mx-auto px-4 py-6 max-w-4xl"
            >
              <Outlet />
            </motion.div>
          </div>
        </main>
      </div>
    </div>
  )
}

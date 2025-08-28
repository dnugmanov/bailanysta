import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Users, UserPlus, Crown, TrendingUp, Search, Filter } from 'lucide-react'
import { apiClient } from '@/api/client'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import UserCard from '@/components/features/UserCard'

export default function CommunityPage() {
  const [activeTab, setActiveTab] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')

  // Fetch real users from API
  const { data: usersData, isLoading, error } = useQuery({
    queryKey: ['community-users'],
    queryFn: () => apiClient.getAllUsers(50, 0), // Get more users for community
  })

  const communityMembers = usersData?.users || []

  // Filter users based on search and active tab
  const filteredMembers = communityMembers.filter(member => {
    const matchesSearch = member.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         (member.bio && member.bio.toLowerCase().includes(searchQuery.toLowerCase()))
    
    switch (activeTab) {
      case 'following':
        return matchesSearch && member.is_following
      case 'popular':
        return matchesSearch && (member.followers_count || 0) > 2 // Adjusted for real data
      default:
        return matchesSearch
    }
  })

  // Stats calculation
  const totalUsers = communityMembers.length
  const followingCount = communityMembers.filter(u => u.is_following).length
  const popularCount = communityMembers.filter(u => (u.followers_count || 0) > 2).length

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Сообщество</h1>
        </div>
        
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i} className="animate-pulse">
              <CardContent className="p-4">
                <div className="flex items-center space-x-3">
                  <div className="h-12 w-12 bg-muted rounded-full"></div>
                  <div className="space-y-2 flex-1">
                    <div className="h-4 bg-muted rounded"></div>
                    <div className="h-3 bg-muted rounded w-2/3"></div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Сообщество</h1>
        </div>
        
        <Card>
          <CardContent className="py-12 text-center">
            <Users className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">Не удалось загрузить пользователей</h3>
            <p className="text-muted-foreground">
              Попробуйте обновить страницу
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Сообщество</h1>
          <p className="text-muted-foreground mt-1">
            Знакомьтесь с участниками образовательной платформы
          </p>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <Users className="h-8 w-8 text-primary" />
                <div className="ml-4">
                  <p className="text-sm font-medium text-muted-foreground">
                    Всего участников
                  </p>
                  <p className="text-2xl font-bold">{totalUsers}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <UserPlus className="h-8 w-8 text-green-500" />
                <div className="ml-4">
                  <p className="text-sm font-medium text-muted-foreground">
                    Ваши подписки
                  </p>
                  <p className="text-2xl font-bold">{followingCount}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center">
                <TrendingUp className="h-8 w-8 text-orange-500" />
                <div className="ml-4">
                  <p className="text-sm font-medium text-muted-foreground">
                    Популярные
                  </p>
                  <p className="text-2xl font-bold">{popularCount}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </div>

      {/* Search and Filter */}
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Поиск участников..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10"
              />
            </div>
            <div className="flex gap-2">
              <Button
                variant={activeTab === 'all' ? 'default' : 'outline'}
                onClick={() => setActiveTab('all')}
                size="sm"
              >
                Все ({totalUsers})
              </Button>
              <Button
                variant={activeTab === 'following' ? 'default' : 'outline'}
                onClick={() => setActiveTab('following')}
                size="sm"
              >
                Подписки ({followingCount})
              </Button>
              <Button
                variant={activeTab === 'popular' ? 'default' : 'outline'}
                onClick={() => setActiveTab('popular')}
                size="sm"
              >
                Популярные ({popularCount})
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Users Grid */}
      {filteredMembers.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <Users className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">
              {searchQuery ? 'Пользователи не найдены' : 'Нет участников'}
            </h3>
            <p className="text-muted-foreground">
              {searchQuery 
                ? 'Попробуйте изменить поисковый запрос'
                : 'Скоро здесь появятся участники сообщества'
              }
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {filteredMembers.map((member, index) => (
            <motion.div
              key={member.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
            >
              <UserCard user={member} compact />
            </motion.div>
          ))}
        </div>
      )}
    </div>
  )
}
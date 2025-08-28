import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import {
  Settings,
  UserPlus,
  UserMinus,
  Calendar,
  MapPin,
  Link as LinkIcon,
  Mail,
} from 'lucide-react'
import { apiClient } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import PostCard from '@/components/features/PostCard'
import { formatDate, formatCount, getInitials } from '@/lib/utils'
import { useToast } from '@/hooks/use-toast'

export default function ProfilePage() {
  const { username } = useParams<{ username: string }>()
  const { user: currentUser } = useAuthStore()
  const { toast } = useToast()
  const queryClient = useQueryClient()

  const { data: user, isLoading: userLoading, error } = useQuery({
    queryKey: ['user', username],
    queryFn: () => {
      if (!username) throw new Error('Username is required')
      // For now, we'll need to get user by username - this might need backend adjustment
      return apiClient.getCurrentUser() // Placeholder
    },
    enabled: !!username,
  })

  const { data: userPosts, isLoading: postsLoading } = useQuery({
    queryKey: ['user-posts', username],
    queryFn: () => {
      // This would be a getUserPosts API call
      return { posts: [] } // Placeholder
    },
    enabled: !!username,
  })

  const followMutation = useMutation({
    mutationFn: (userId: string) => 
      user?.is_following ? apiClient.unfollowUser(userId) : apiClient.followUser(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', username] })
      toast({
        title: user?.is_following ? 'Отписка' : 'Подписка',
        description: user?.is_following 
          ? 'Вы отписались от пользователя' 
          : 'Вы подписались на пользователя',
      })
    },
    onError: () => {
      toast({
        title: 'Ошибка',
        description: 'Не удалось изменить статус подписки',
        variant: 'destructive',
      })
    },
  })

  if (userLoading) {
    return (
      <div className="space-y-6">
        {/* Profile skeleton */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-start gap-6">
              <div className="w-24 h-24 bg-muted rounded-full animate-pulse" />
              <div className="flex-1 space-y-3">
                <div className="h-6 bg-muted rounded animate-pulse w-1/3" />
                <div className="h-4 bg-muted rounded animate-pulse w-1/2" />
                <div className="h-4 bg-muted rounded animate-pulse w-2/3" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (error || !user) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-bold text-destructive">Пользователь не найден</h2>
        <p className="text-muted-foreground mt-2">
          Пользователь с именем @{username} не существует
        </p>
      </div>
    )
  }

  const isOwnProfile = currentUser?.id === user.id

  return (
    <div className="space-y-6">
      {/* Profile Header */}
      <Card>
        <CardContent className="p-6">
          <div className="flex flex-col sm:flex-row items-start gap-6">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ type: "spring", stiffness: 200 }}
            >
              <Avatar className="w-24 h-24 border-4 border-primary/20">
                <AvatarImage src={user.avatar_url || undefined} />
                <AvatarFallback className="text-2xl">
                  {getInitials(user.username)}
                </AvatarFallback>
              </Avatar>
            </motion.div>

            <div className="flex-1 min-w-0">
              <div className="flex flex-col sm:flex-row sm:items-center gap-4 mb-4">
                <div>
                  <h1 className="text-2xl font-bold">{user.username}</h1>
                  <p className="text-muted-foreground">{user.email}</p>
                </div>

                <div className="flex gap-2 sm:ml-auto">
                  {isOwnProfile ? (
                    <Button variant="outline">
                      <Settings className="w-4 h-4 mr-2" />
                      Редактировать профиль
                    </Button>
                  ) : (
                    <Button
                      onClick={() => followMutation.mutate(user.id)}
                      disabled={followMutation.isPending}
                      variant={user.is_following ? "outline" : "default"}
                    >
                      {user.is_following ? (
                        <>
                          <UserMinus className="w-4 h-4 mr-2" />
                          Отписаться
                        </>
                      ) : (
                        <>
                          <UserPlus className="w-4 h-4 mr-2" />
                          Подписаться
                        </>
                      )}
                    </Button>
                  )}
                </div>
              </div>

              {user.bio && (
                <p className="text-sm mb-4 leading-relaxed">{user.bio}</p>
              )}

              <div className="flex flex-wrap gap-4 text-sm text-muted-foreground mb-4">
                <div className="flex items-center gap-1">
                  <Calendar className="w-4 h-4" />
                  Регистрация {formatDate(user.created_at)}
                </div>
              </div>

              <div className="flex gap-6 text-sm">
                <div>
                  <span className="font-semibold text-foreground">
                    {formatCount(user.following_count || 0)}
                  </span>
                  <span className="text-muted-foreground ml-1">подписок</span>
                </div>
                <div>
                  <span className="font-semibold text-foreground">
                    {formatCount(user.followers_count || 0)}
                  </span>
                  <span className="text-muted-foreground ml-1">подписчиков</span>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Profile Tabs */}
      <Tabs defaultValue="posts" className="space-y-4">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="posts">Посты</TabsTrigger>
          <TabsTrigger value="liked">Понравилось</TabsTrigger>
          <TabsTrigger value="media">Медиа</TabsTrigger>
        </TabsList>

        <TabsContent value="posts" className="space-y-4">
          {postsLoading ? (
            Array.from({ length: 3 }).map((_, i) => (
              <Card key={i} className="post-skeleton">
                <CardContent className="p-6">
                  <div className="space-y-3">
                    <div className="h-4 bg-muted rounded animate-pulse" />
                    <div className="h-4 bg-muted rounded animate-pulse w-3/4" />
                    <div className="h-4 bg-muted rounded animate-pulse w-1/2" />
                  </div>
                </CardContent>
              </Card>
            ))
          ) : userPosts?.posts.length === 0 ? (
            <Card>
              <CardContent className="py-12 text-center">
                <p className="text-muted-foreground">
                  {isOwnProfile ? 'Вы еще не создали ни одного поста' : 'Пользователь еще не создал постов'}
                </p>
              </CardContent>
            </Card>
          ) : (
            userPosts?.posts.map((post, index) => (
              <motion.div
                key={post.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                <PostCard post={post} />
              </motion.div>
            ))
          )}
        </TabsContent>

        <TabsContent value="liked">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Понравившиеся посты пока не загружены</p>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="media">
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">Медиа файлы пока не поддерживаются</p>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { UserPlus, UserMinus } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/types/api'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { getInitials, formatCount } from '@/lib/utils'
import { useToast } from '@/hooks/use-toast'

interface UserCardProps {
  user: User
  compact?: boolean
}

export default function UserCard({ user, compact = false }: UserCardProps) {
  const { user: currentUser } = useAuthStore()
  const { toast } = useToast()
  const queryClient = useQueryClient()

  const followMutation = useMutation({
    mutationFn: () => 
      user.is_following ? apiClient.unfollowUser(user.id) : apiClient.followUser(user.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['search'] })
      queryClient.invalidateQueries({ queryKey: ['user', user.username] })
      toast({
        title: user.is_following ? 'Отписка' : 'Подписка',
        description: user.is_following 
          ? `Вы отписались от @${user.username}` 
          : `Вы подписались на @${user.username}`,
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

  const isOwnProfile = currentUser?.id === user.id

  if (compact) {
    return (
      <motion.div
        whileHover={{ scale: 1.02 }}
        transition={{ duration: 0.2 }}
      >
        <Card className="hover:shadow-md transition-shadow">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <Link to={`/profile/${user.username}`} className="flex items-center gap-3 flex-1 min-w-0">
                <Avatar className="h-10 w-10">
                  <AvatarImage src={user.avatar_url || undefined} />
                  <AvatarFallback>{getInitials(user.username)}</AvatarFallback>
                </Avatar>
                <div className="min-w-0 flex-1">
                  <p className="font-semibold truncate">{user.username}</p>
                  <p className="text-sm text-muted-foreground truncate">{user.email}</p>
                </div>
              </Link>
              
              {!isOwnProfile && (
                <Button
                  size="sm"
                  variant={user.is_following ? "outline" : "default"}
                  onClick={() => followMutation.mutate()}
                  disabled={followMutation.isPending}
                >
                  {user.is_following ? (
                    <UserMinus className="h-4 w-4" />
                  ) : (
                    <UserPlus className="h-4 w-4" />
                  )}
                </Button>
              )}
            </div>
          </CardContent>
        </Card>
      </motion.div>
    )
  }

  return (
    <motion.div
      whileHover={{ scale: 1.02 }}
      transition={{ duration: 0.2 }}
    >
      <Card className="hover:shadow-md transition-shadow">
        <CardContent className="p-6">
          <div className="flex items-start gap-4">
            <Link to={`/profile/${user.username}`}>
              <Avatar className="h-16 w-16 hover:ring-2 hover:ring-primary/20 transition-all">
                <AvatarImage src={user.avatar_url || undefined} />
                <AvatarFallback className="text-lg">
                  {getInitials(user.username)}
                </AvatarFallback>
              </Avatar>
            </Link>

            <div className="flex-1 min-w-0">
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0 flex-1">
                  <Link
                    to={`/profile/${user.username}`}
                    className="font-semibold text-lg hover:text-primary transition-colors truncate block"
                  >
                    {user.username}
                  </Link>
                  <p className="text-muted-foreground truncate">{user.email}</p>
                  
                  {user.bio && (
                    <p className="text-sm mt-2 leading-relaxed line-clamp-2">
                      {user.bio}
                    </p>
                  )}

                  <div className="flex gap-4 mt-3 text-sm">
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

                {!isOwnProfile && (
                  <Button
                    variant={user.is_following ? "outline" : "default"}
                    onClick={() => followMutation.mutate()}
                    disabled={followMutation.isPending}
                    className="shrink-0"
                  >
                    {user.is_following ? (
                      <>
                        <UserMinus className="h-4 w-4 mr-2" />
                        Отписаться
                      </>
                    ) : (
                      <>
                        <UserPlus className="h-4 w-4 mr-2" />
                        Подписаться
                      </>
                    )}
                  </Button>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  )
}

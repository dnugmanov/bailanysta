import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { useNavigate } from 'react-router-dom'
import { Check, Bell, Heart, MessageCircle, UserPlus, Trash2, FileText } from 'lucide-react'
import { apiClient } from '@/api/client'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { formatDate } from '@/lib/utils'
import { useToast } from '@/hooks/use-toast'
import type { Notification } from '@/types/api'

export default function NotificationsPage() {
  const [activeTab, setActiveTab] = useState('all')
  const { toast } = useToast()
  const queryClient = useQueryClient()

  const { data: allNotifications, isLoading, error: allError } = useQuery({
    queryKey: ['notifications', 'all'],
    queryFn: () => apiClient.getNotifications(false),
    retry: 1,
  })

  const { data: unreadNotifications, error: unreadError } = useQuery({
    queryKey: ['notifications', 'unread'],
    queryFn: () => apiClient.getNotifications(true),
    retry: 1,
  })

  // Debug logging
  console.log('All notifications:', allNotifications)
  console.log('Unread notifications:', unreadNotifications)
  console.log('Is loading:', isLoading)
  console.log('All error:', allError)
  console.log('Unread error:', unreadError)

  const markAllAsReadMutation = useMutation({
    mutationFn: () => apiClient.markAllNotificationsAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
      toast({
        title: 'Успешно',
        description: 'Все уведомления отмечены как прочитанные',
      })
    },
  })

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) => apiClient.markNotificationAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
    },
  })

  const deleteNotificationMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteNotification(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
      toast({
        title: 'Удалено',
        description: 'Уведомление удалено',
      })
    },
  })

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'like':
        return <Heart className="h-4 w-4 text-red-500" />
      case 'comment':
        return <MessageCircle className="h-4 w-4 text-blue-500" />
      case 'follow':
        return <UserPlus className="h-4 w-4 text-green-500" />
      case 'new_post':
        return <FileText className="h-4 w-4 text-purple-500" />
      default:
        return <Bell className="h-4 w-4" />
    }
  }

  const getNotificationText = (notification: Notification) => {
    const { type, actor } = notification
    const username = actor?.username || 'Пользователь'

    switch (type) {
      case 'like':
        return `${username} поставил лайк на ваш пост`
      case 'comment':
        return `${username} прокомментировал ваш пост`
      case 'follow':
        return `${username} подписался на вас`
      case 'new_post':
        return `${username} опубликовал новый пост`
      default:
        return 'Новое уведомление'
    }
  }

  const NotificationCard = ({ notification }: { notification: Notification }) => {
    const navigate = useNavigate()
    const isUnread = !notification.read_at

    const handleClick = () => {
      // Mark as read first
      if (isUnread) {
        markAsReadMutation.mutate(notification.id)
      }
      
      // Navigate based on notification type
      if (notification.type === 'like' || notification.type === 'comment' || notification.type === 'new_post') {
        if (notification.entity_id) {
          navigate(`/post/${notification.entity_id}`)
        }
      } else if (notification.type === 'follow') {
        if (notification.actor?.username) {
          navigate(`/profile/${notification.actor.username}`)
        }
      }
    }

    return (
      <Card 
        className={`transition-all cursor-pointer hover:shadow-md ${isUnread ? 'ring-2 ring-primary/20 bg-primary/5' : ''}`}
        onClick={handleClick}
      >
        <CardContent className="p-4">
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 mt-1">
              {getNotificationIcon(notification.type)}
            </div>
            
            <div className="flex-1 min-w-0">
              <div className="flex items-start justify-between gap-2">
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium">
                    {getNotificationText(notification)}
                  </p>
                  {notification.post && (
                    <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                      {notification.post.text.substring(0, 100)}...
                    </p>
                  )}
                  {notification.payload.comment_text && (
                    <p className="text-xs text-muted-foreground mt-1 italic">
                      "{notification.payload.comment_text}"
                    </p>
                  )}
                  <p className="text-xs text-muted-foreground mt-1">
                    {formatDate(notification.created_at)}
                  </p>
                </div>
                
                <div className="flex items-center gap-1">
                  {isUnread && (
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={(e) => {
                        e.stopPropagation()
                        markAsReadMutation.mutate(notification.id)
                      }}
                    >
                      <Check className="h-3 w-3" />
                    </Button>
                  )}
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={(e) => {
                      e.stopPropagation()
                      deleteNotificationMutation.mutate(notification.id)
                    }}
                  >
                    <Trash2 className="h-3 w-3" />
                  </Button>
                </div>
              </div>
            </div>
            
            {isUnread && (
              <div className="flex-shrink-0 w-2 h-2 bg-primary rounded-full mt-2" />
            )}
          </div>
        </CardContent>
      </Card>
    )
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <div className="h-8 bg-muted rounded animate-pulse w-48" />
            <div className="h-4 bg-muted rounded animate-pulse w-64 mt-2" />
          </div>
        </div>

        <div className="space-y-4">
          {Array.from({ length: 5 }).map((_, i) => (
            <Card key={i} className="post-skeleton">
              <CardContent className="p-4">
                <div className="flex items-center space-x-3">
                  <div className="w-4 h-4 bg-muted rounded animate-pulse" />
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-muted rounded animate-pulse w-3/4" />
                    <div className="h-3 bg-muted rounded animate-pulse w-1/4" />
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Уведомления</h1>
          <p className="text-muted-foreground">
            Следите за активностью в вашем профиле
          </p>
        </div>
        
        {(unreadNotifications?.length || 0) > 0 && (
          <Button
            variant="outline"
            onClick={() => markAllAsReadMutation.mutate()}
            disabled={markAllAsReadMutation.isPending}
          >
            <Check className="h-4 w-4 mr-2" />
            Отметить все как прочитанные
          </Button>
        )}
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="all">
            Все ({allNotifications?.length || 0})
          </TabsTrigger>
          <TabsTrigger value="unread">
            Непрочитанные ({unreadNotifications?.length || 0})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="all" className="space-y-4">
          {allNotifications && allNotifications.length > 0 ? (
            allNotifications.map((notification, index) => (
              <motion.div
                key={notification.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
              >
                <NotificationCard notification={notification} />
              </motion.div>
            ))
          ) : (
            <Card>
              <CardContent className="py-12 text-center">
                <Bell className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">Уведомлений пока нет</h3>
                <p className="text-muted-foreground">
                  Когда кто-то поставит лайк, оставит комментарий или подпишется на вас, здесь появятся уведомления
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="unread" className="space-y-4">
          {unreadNotifications && unreadNotifications.length > 0 ? (
            unreadNotifications.map((notification, index) => (
              <motion.div
                key={notification.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
              >
                <NotificationCard notification={notification} />
              </motion.div>
            ))
          ) : (
            <Card>
              <CardContent className="py-12 text-center">
                <Check className="h-12 w-12 text-green-500 mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">Все уведомления прочитаны</h3>
                <p className="text-muted-foreground">
                  Отлично! Вы в курсе всех событий
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}

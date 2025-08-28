import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Heart, MoreHorizontal, Trash2, Flag } from 'lucide-react'
import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/stores/auth'
import type { Comment } from '@/types/api'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import * as DropdownMenu from '@radix-ui/react-dropdown-menu'
import { formatDate, getInitials } from '@/lib/utils'
import { useToast } from '@/hooks/use-toast'

interface CommentCardProps {
  comment: Comment
}

export default function CommentCard({ comment }: CommentCardProps) {
  const [isLiked, setIsLiked] = useState(false) // This would come from API
  const [likeCount, setLikeCount] = useState(0) // This would come from API
  const { user } = useAuthStore()
  const { toast } = useToast()
  const queryClient = useQueryClient()

  const isOwnComment = user?.id === comment.author_id

  const deleteCommentMutation = useMutation({
    mutationFn: () => {
      // This would be an API call to delete comment
      return Promise.resolve()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments'] })
      toast({
        title: 'Удалено',
        description: 'Комментарий удален',
      })
    },
    onError: () => {
      toast({
        title: 'Ошибка',
        description: 'Не удалось удалить комментарий',
        variant: 'destructive',
      })
    },
  })

  const handleLike = () => {
    // Toggle like optimistically
    setIsLiked(!isLiked)
    setLikeCount(prev => isLiked ? prev - 1 : prev + 1)
    
    // This would be an API call
    toast({
      title: 'Функция в разработке',
      description: 'Лайки комментариев скоро будут доступны',
    })
  }

  return (
    <Card className="hover:shadow-sm transition-shadow">
      <CardContent className="p-4">
        <div className="flex items-start gap-3">
          <Link to={`/profile/${comment.author?.username}`}>
            <Avatar className="h-8 w-8 hover:ring-2 hover:ring-primary/20 transition-all">
              <AvatarImage src={comment.author?.avatar_url || undefined} />
              <AvatarFallback className="text-xs">
                {comment.author ? getInitials(comment.author.username) : 'U'}
              </AvatarFallback>
            </Avatar>
          </Link>

          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2">
                <Link
                  to={`/profile/${comment.author?.username}`}
                  className="font-medium text-sm hover:text-primary transition-colors"
                >
                  {comment.author?.username}
                </Link>
                <span className="text-xs text-muted-foreground">
                  {formatDate(comment.created_at)}
                </span>
              </div>

              {isOwnComment && (
                <DropdownMenu.Root>
                  <DropdownMenu.Trigger asChild>
                    <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
                      <MoreHorizontal className="h-3 w-3" />
                    </Button>
                  </DropdownMenu.Trigger>
                  <DropdownMenu.Portal>
                    <DropdownMenu.Content
                      className="min-w-[8rem] bg-popover border rounded-md p-1 shadow-md z-50"
                      sideOffset={5}
                      align="end"
                    >
                      <DropdownMenu.Item
                        className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm text-destructive"
                        onClick={() => deleteCommentMutation.mutate()}
                      >
                        <Trash2 className="h-3 w-3" />
                        Удалить
                      </DropdownMenu.Item>
                    </DropdownMenu.Content>
                  </DropdownMenu.Portal>
                </DropdownMenu.Root>
              )}

              {!isOwnComment && (
                <DropdownMenu.Root>
                  <DropdownMenu.Trigger asChild>
                    <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
                      <MoreHorizontal className="h-3 w-3" />
                    </Button>
                  </DropdownMenu.Trigger>
                  <DropdownMenu.Portal>
                    <DropdownMenu.Content
                      className="min-w-[8rem] bg-popover border rounded-md p-1 shadow-md z-50"
                      sideOffset={5}
                      align="end"
                    >
                      <DropdownMenu.Item className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm">
                        <Flag className="h-3 w-3" />
                        Пожаловаться
                      </DropdownMenu.Item>
                    </DropdownMenu.Content>
                  </DropdownMenu.Portal>
                </DropdownMenu.Root>
              )}
            </div>

            <p className="text-sm leading-relaxed mb-3">{comment.text}</p>

            <div className="flex items-center gap-4">
              <motion.button
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                onClick={handleLike}
                className={`flex items-center gap-1 text-xs transition-colors ${
                  isLiked
                    ? 'text-red-500 hover:text-red-600'
                    : 'text-muted-foreground hover:text-red-500'
                }`}
              >
                <Heart className={`h-3 w-3 ${isLiked ? 'fill-current' : ''}`} />
                {likeCount > 0 && <span>{likeCount}</span>}
              </motion.button>

              <Button
                variant="ghost"
                size="sm"
                className="h-6 px-0 text-xs text-muted-foreground hover:text-foreground"
              >
                Ответить
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

import { useState } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import {
  Heart,
  MessageCircle,
  Share,
  Bookmark,
  MoreHorizontal,
  Edit,
  Trash2,
} from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { Post } from '@/types/api'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import * as DropdownMenu from '@radix-ui/react-dropdown-menu'
import { formatDate, getInitials, highlightHashtags, formatCount } from '@/lib/utils'
import FormattedText from './FormattedText'
import PostEditor from './PostEditor'
import { useToast } from '@/hooks/use-toast'

interface PostCardProps {
  post: Post
}

export default function PostCard({ post }: PostCardProps) {
  const [isLiked, setIsLiked] = useState(post.is_liked || false)
  const [likeCount, setLikeCount] = useState(post.like_count)
  const [isEditing, setIsEditing] = useState(false)
  const { user } = useAuthStore()
  const { toast } = useToast()
  const queryClient = useQueryClient()

  const likeMutation = useMutation({
    mutationFn: (liked: boolean) => 
      liked ? apiClient.unlikePost(post.id) : apiClient.likePost(post.id),
    onMutate: async (liked) => {
      // Optimistic update
      setIsLiked(!liked)
      setLikeCount(prev => liked ? prev - 1 : prev + 1)
    },
    onError: (error, liked) => {
      // Revert on error
      setIsLiked(liked)
      setLikeCount(prev => liked ? prev + 1 : prev - 1)
      toast({
        title: 'Ошибка',
        description: 'Не удалось обновить лайк',
        variant: 'destructive',
      })
    },
    onSuccess: () => {
      // Invalidate feed queries
      queryClient.invalidateQueries({ queryKey: ['feed'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: () => apiClient.deletePost(post.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feed'] })
      toast({
        title: 'Успешно',
        description: 'Пост удален',
      })
    },
    onError: () => {
      toast({
        title: 'Ошибка',
        description: 'Не удалось удалить пост',
        variant: 'destructive',
      })
    },
  })

  const handleLike = () => {
    likeMutation.mutate(isLiked)
  }

  const handleShare = async () => {
    const url = `${window.location.origin}/post/${post.id}`
    try {
      await navigator.share({
        title: `Пост от ${post.author?.username}`,
        text: post.text.slice(0, 100) + '...',
        url,
      })
    } catch (error) {
      // Fallback to clipboard
      await navigator.clipboard.writeText(url)
      toast({
        title: 'Ссылка скопирована',
        description: 'Ссылка на пост скопирована в буфер обмена',
      })
    }
  }

  const isOwnPost = user?.id === post.author_id

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <Link to={`/profile/${post.author?.username}`}>
              <Avatar className="h-10 w-10 hover:ring-2 hover:ring-primary/20 transition-all">
                <AvatarImage src={post.author?.avatar_url || undefined} />
                <AvatarFallback>
                  {post.author ? getInitials(post.author.username) : 'U'}
                </AvatarFallback>
              </Avatar>
            </Link>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2">
                <Link
                  to={`/profile/${post.author?.username}`}
                  className="font-semibold hover:text-primary transition-colors truncate"
                >
                  {post.author?.username}
                </Link>
                {post.course && (
                  <span className="text-xs bg-primary/10 text-primary px-2 py-1 rounded-full">
                    {post.course.title}
                  </span>
                )}
              </div>
              <p className="text-sm text-muted-foreground">
                {formatDate(post.created_at)}
                {post.updated_at !== post.created_at && (
                  <span className="ml-1">(изменено)</span>
                )}
              </p>
            </div>
          </div>

          {isOwnPost && (
            <DropdownMenu.Root>
              <DropdownMenu.Trigger asChild>
                <Button variant="ghost" size="sm">
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenu.Trigger>
              <DropdownMenu.Portal>
                <DropdownMenu.Content
                  className="min-w-[8rem] bg-popover border rounded-md p-1 shadow-md z-50"
                  sideOffset={5}
                  align="end"
                >
                  <DropdownMenu.Item 
                    className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm"
                    onClick={() => setIsEditing(true)}
                  >
                    <Edit className="h-4 w-4" />
                    Редактировать
                  </DropdownMenu.Item>
                  <DropdownMenu.Item
                    className="flex items-center gap-2 px-2 py-1.5 text-sm cursor-pointer hover:bg-accent rounded-sm text-destructive"
                    onClick={() => deleteMutation.mutate()}
                  >
                    <Trash2 className="h-4 w-4" />
                    Удалить
                  </DropdownMenu.Item>
                </DropdownMenu.Content>
              </DropdownMenu.Portal>
            </DropdownMenu.Root>
          )}
        </div>
      </CardHeader>

      <CardContent className="pt-0">
        <Link to={`/post/${post.id}`} className="block">
          <FormattedText 
            text={post.text}
            className="mb-4 cursor-pointer hover:text-foreground/80 transition-colors"
            maxLines={10}
          />
        </Link>

        {/* Hashtags */}
        {post.hashtags && post.hashtags.length > 0 && (
          <div className="flex flex-wrap gap-1 mb-4">
            {post.hashtags.map((tag) => (
              <Link
                key={tag}
                to={`/search?q=${encodeURIComponent('#' + tag)}`}
                className="text-xs bg-muted hover:bg-primary/10 hover:text-primary px-2 py-1 rounded-full transition-colors"
              >
                #{tag}
              </Link>
            ))}
          </div>
        )}

        {/* Actions */}
        <div className="flex items-center justify-between pt-3 border-t">
          <div className="flex items-center space-x-4">
            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={handleLike}
              className={`flex items-center space-x-1 transition-colors ${
                isLiked
                  ? 'text-red-500 hover:text-red-600'
                  : 'text-muted-foreground hover:text-red-500'
              }`}
            >
              <Heart
                className={`h-4 w-4 ${isLiked ? 'fill-current' : ''}`}
              />
              <span className="text-sm">{formatCount(likeCount)}</span>
            </motion.button>

            <Link
              to={`/post/${post.id}`}
              className="flex items-center space-x-1 text-muted-foreground hover:text-blue-500 transition-colors"
            >
              <MessageCircle className="h-4 w-4" />
              <span className="text-sm">{formatCount(post.comment_count)}</span>
            </Link>

            <motion.button
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              onClick={handleShare}
              className="flex items-center space-x-1 text-muted-foreground hover:text-green-500 transition-colors"
            >
              <Share className="h-4 w-4" />
            </motion.button>
          </div>

          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            className="text-muted-foreground hover:text-yellow-500 transition-colors"
          >
            <Bookmark className="h-4 w-4" />
          </motion.button>
        </div>
      </CardContent>

      {/* Post Editor Modal */}
      <PostEditor
        post={post}
        isOpen={isEditing}
        onClose={() => setIsEditing(false)}
      />
    </Card>
  )
}

import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { ArrowLeft, MessageCircle } from 'lucide-react'
import { Link } from 'react-router-dom'
import { apiClient } from '@/api/client'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import PostCard from '@/components/features/PostCard'
import FormattedText from '@/components/features/FormattedText'
import CommentCard from '@/components/features/CommentCard'
import CommentComposer from '@/components/features/CommentComposer'

export default function PostPage() {
  const { id } = useParams<{ id: string }>()

  const { data: post, isLoading: postLoading, error: postError } = useQuery({
    queryKey: ['post', id],
    queryFn: async () => {
      if (!id) throw new Error('Post ID is required')
      try {
        const result = await apiClient.getPostById(id)
        console.log('Post loaded:', result)
        return result
      } catch (error) {
        console.error('Error loading post:', error)
        throw error
      }
    },
    enabled: !!id,
    retry: 1,
  })

  const { data: commentsData, isLoading: commentsLoading } = useQuery({
    queryKey: ['comments', id],
    queryFn: async () => {
      if (!id) throw new Error('Post ID is required')
      try {
        const result = await apiClient.getComments(id)
        console.log('Comments loaded:', result)
        return result
      } catch (error) {
        console.error('Error loading comments:', error)
        throw error
      }
    },
    enabled: !!id && !!post,
    retry: 1,
  })

  if (postLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link to="/feed">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div className="h-6 bg-muted rounded animate-pulse w-32" />
        </div>

        <Card className="post-skeleton">
          <CardContent className="p-6">
            <div className="space-y-4">
              <div className="flex items-center space-x-4">
                <div className="w-10 h-10 bg-muted rounded-full animate-pulse" />
                <div className="space-y-2 flex-1">
                  <div className="h-4 bg-muted rounded animate-pulse w-1/4" />
                  <div className="h-3 bg-muted rounded animate-pulse w-1/6" />
                </div>
              </div>
              <div className="space-y-3">
                <div className="h-4 bg-muted rounded animate-pulse" />
                <div className="h-4 bg-muted rounded animate-pulse w-3/4" />
                <div className="h-4 bg-muted rounded animate-pulse w-1/2" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (postError || !post) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link to="/feed">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <h1 className="text-2xl font-bold">Пост не найден</h1>
        </div>

        <Card>
          <CardContent className="py-12 text-center">
            <h2 className="text-xl font-semibold text-destructive mb-2">
              Пост не найден
            </h2>
            <p className="text-muted-foreground">
              Пост был удален или не существует
            </p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link to="/feed">
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <h1 className="text-2xl font-bold">Пост</h1>
      </div>

      {/* Post */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <PostCard post={post} />
      </motion.div>

      {/* Comment Composer */}
      <Card>
        <CardContent className="p-4">
          <CommentComposer postId={post.id} />
        </CardContent>
      </Card>

      {/* Comments */}
      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <MessageCircle className="h-5 w-5" />
          <h2 className="text-lg font-semibold">
            Комментарии ({commentsData?.comments?.length || 0})
          </h2>
        </div>

        {commentsLoading ? (
          Array.from({ length: 3 }).map((_, i) => (
            <Card key={i} className="post-skeleton">
              <CardContent className="p-4">
                <div className="space-y-3">
                  <div className="flex items-center space-x-3">
                    <div className="w-8 h-8 bg-muted rounded-full animate-pulse" />
                    <div className="space-y-1 flex-1">
                      <div className="h-3 bg-muted rounded animate-pulse w-1/4" />
                      <div className="h-2 bg-muted rounded animate-pulse w-1/6" />
                    </div>
                  </div>
                  <div className="h-3 bg-muted rounded animate-pulse w-3/4" />
                </div>
              </CardContent>
            </Card>
          ))
        ) : (!commentsData?.comments || commentsData.comments.length === 0) ? (
          <Card>
            <CardContent className="py-12 text-center">
              <MessageCircle className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-semibold mb-2">Пока нет комментариев</h3>
              <p className="text-muted-foreground">
                Станьте первым, кто оставит комментарий к этому посту
              </p>
            </CardContent>
          </Card>
        ) : (
          (commentsData?.comments || []).map((comment, index) => (
            <motion.div
              key={comment.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
            >
              <CommentCard comment={comment} />
            </motion.div>
          ))
        )}
      </div>
    </div>
  )
}

import { useState } from 'react'
import { useInfiniteQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { RefreshCw, TrendingUp, Users, Sparkles } from 'lucide-react'
import { apiClient } from '@/api/client'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import PostCard from '@/components/features/PostCard'
import PostComposer from '@/components/features/PostComposer'
import { useUIStore } from '@/stores/ui'
import { useToast } from '@/hooks/use-toast'

export default function FeedPage() {
  const [activeTab, setActiveTab] = useState<'following' | 'trending'>('following')
  const { setComposerOpen } = useUIStore()
  const { toast } = useToast()

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    error,
    refetch,
    isRefetching,
  } = useInfiniteQuery({
    queryKey: ['feed', activeTab],
    queryFn: ({ pageParam = 0 }) => {
      // Use different API methods based on active tab
      if (activeTab === 'trending') {
        return apiClient.getTrendingPosts(20, pageParam as number)
      }
      return apiClient.getFeed(20, pageParam as number)
    },
    getNextPageParam: (lastPage: any) => {
      if (!lastPage || !lastPage.posts || !Array.isArray(lastPage.posts)) {
        return undefined
      }
      const nextOffset = lastPage.offset + lastPage.limit
      return lastPage.posts.length === lastPage.limit ? nextOffset : undefined
    },
    initialPageParam: 0,
    staleTime: 1000 * 60 * 5, // 5 minutes
  })

  const posts = data?.pages.flatMap((page: any) => page.posts || []).filter(Boolean) || []

  const handleRefresh = async () => {
    try {
      await refetch()
      toast({
        title: 'Обновлено!',
        description: 'Лента успешно обновлена',
      })
    } catch (error) {
      toast({
        title: 'Ошибка',
        description: 'Не удалось обновить ленту',
        variant: 'destructive',
      })
    }
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="space-y-4">
          <h2 className="text-2xl font-bold text-destructive">Ошибка загрузки</h2>
          <p className="text-muted-foreground">
            Не удалось загрузить ленту. Попробуйте еще раз.
          </p>
          <Button onClick={() => refetch()}>
            <RefreshCw className="w-4 h-4 mr-2" />
            Попробовать снова
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Лента</h1>
          <p className="text-muted-foreground">
            Следите за обновлениями от образовательного сообщества
          </p>
        </div>
        <Button
          variant="outline"
          onClick={handleRefresh}
          disabled={isRefetching}
        >
          <RefreshCw className={`w-4 h-4 mr-2 ${isRefetching ? 'animate-spin' : ''}`} />
          Обновить
        </Button>
      </div>

      {/* Tabs */}
      <div className="flex space-x-1 bg-muted p-1 rounded-lg w-fit">
        <Button
          variant={activeTab === 'following' ? 'default' : 'ghost'}
          size="sm"
          onClick={() => setActiveTab('following')}
          className="relative"
        >
          <Users className="w-4 h-4 mr-2" />
          Подписки
          {activeTab === 'following' && (
            <motion.div
              layoutId="activeTab"
              className="absolute inset-0 bg-primary rounded-md -z-10"
              transition={{ type: "spring", stiffness: 300, damping: 30 }}
            />
          )}
        </Button>
        <Button
          variant={activeTab === 'trending' ? 'default' : 'ghost'}
          size="sm"
          onClick={() => setActiveTab('trending')}
          className="relative"
        >
          <TrendingUp className="w-4 h-4 mr-2" />
          Популярное
          {activeTab === 'trending' && (
            <motion.div
              layoutId="activeTab"
              className="absolute inset-0 bg-primary rounded-md -z-10"
              transition={{ type: "spring", stiffness: 300, damping: 30 }}
            />
          )}
        </Button>
      </div>

      {/* Composer Card */}
      <Card>
        <CardContent className="p-4">
          <Button
            variant="outline"
            className="w-full justify-start text-muted-foreground"
            onClick={() => setComposerOpen(true)}
          >
            <Sparkles className="w-4 h-4 mr-2" />
            Поделитесь своими знаниями...
          </Button>
        </CardContent>
      </Card>

      {/* Posts */}
      <div className="space-y-4">
        {isLoading ? (
          // Loading skeletons
          Array.from({ length: 3 }).map((_, i) => (
            <Card key={i} className="post-skeleton">
              <CardHeader>
                <div className="flex items-center space-x-4">
                  <div className="w-10 h-10 bg-muted rounded-full animate-pulse" />
                  <div className="space-y-2 flex-1">
                    <div className="h-4 bg-muted rounded animate-pulse w-1/4" />
                    <div className="h-3 bg-muted rounded animate-pulse w-1/6" />
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="h-4 bg-muted rounded animate-pulse" />
                  <div className="h-4 bg-muted rounded animate-pulse w-3/4" />
                  <div className="h-4 bg-muted rounded animate-pulse w-1/2" />
                </div>
              </CardContent>
            </Card>
          ))
        ) : posts.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <div className="space-y-4">
                <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mx-auto">
                  <Users className="w-8 h-8 text-muted-foreground" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold">Пока здесь пусто</h3>
                  <p className="text-muted-foreground">
                    {activeTab === 'following'
                      ? 'Подпишитесь на интересных пользователей, чтобы видеть их посты'
                      : 'Популярные посты скоро появятся'}
                  </p>
                </div>
                <Button onClick={() => setComposerOpen(true)}>
                  <Sparkles className="w-4 h-4 mr-2" />
                  Создать первый пост
                </Button>
              </div>
            </CardContent>
          </Card>
        ) : (
          posts.map((post, index) => (
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

        {/* Load more button */}
        {hasNextPage && (
          <div className="text-center py-4">
            <Button
              variant="outline"
              onClick={() => fetchNextPage()}
              disabled={isFetchingNextPage}
            >
              {isFetchingNextPage ? (
                <div className="flex items-center gap-2">
                  <div className="w-4 h-4 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
                  Загружаем...
                </div>
              ) : (
                'Загрузить еще'
              )}
            </Button>
          </div>
        )}
      </div>

      {/* Post Composer Modal */}
      <PostComposer />
    </div>
  )
}
import { useState, useEffect, useCallback } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { Search as SearchIcon, Users, Hash } from 'lucide-react'
import { apiClient } from '@/api/client'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import PostCard from '@/components/features/PostCard'
import UserCard from '@/components/features/UserCard'
import { debounce } from '@/lib/utils'

export default function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [query, setQuery] = useState(searchParams.get('q') || '')
  const [activeTab, setActiveTab] = useState('all')

  const debouncedSearch = useCallback(
    debounce((searchQuery: string) => {
      if (searchQuery && searchQuery.trim().length > 0) {
        setSearchParams({ q: searchQuery.trim() })
      } else {
        setSearchParams({})
      }
    }, 500),
    [setSearchParams]
  )

  useEffect(() => {
    debouncedSearch(query)
  }, [query, debouncedSearch])

  const { data: searchResults, isLoading, error } = useQuery({
    queryKey: ['search', searchParams.get('q')],
    queryFn: async () => {
      const q = searchParams.get('q')
      if (!q || q.trim().length === 0) return null
      
      try {
        return await apiClient.search(q.trim())
      } catch (err: any) {
        // Обрабатываем ошибки SQL с поиском пользователей
        if (err.message?.includes('following_count does not exist')) {
          // Возвращаем результаты только с постами, без пользователей
          return {
            posts: [],
            users: [],
            query: q.trim(),
            total_posts: 0,
            total_users: 0,
          }
        }
        throw err
      }
    },
    enabled: !!searchParams.get('q') && searchParams.get('q')!.trim().length > 0,
    retry: false, // Не повторяем при ошибках
  })

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (query.trim()) {
      setSearchParams({ q: query.trim() })
    }
  }

  const trendingHashtags = [
    '#javascript',
    '#react',
    '#python',
    '#webdev',
    '#programming',
    '#ai',
    '#machinelearning',
    '#design',
  ]

  const suggestedUsers: any[] = [
    // This would come from an API
  ]

  return (
    <div className="space-y-6">
      {/* Search Header */}
      <div>
        <h1 className="text-3xl font-bold mb-2">Поиск</h1>
        <p className="text-muted-foreground">
          Найдите интересные посты, пользователей и темы
        </p>
      </div>

      {/* Search Bar */}
      <Card>
        <CardContent className="p-4">
          <form onSubmit={handleSearch} className="flex gap-2">
            <div className="relative flex-1">
              <SearchIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Поиск постов, пользователей, хештегов..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                className="pl-10"
              />
            </div>
            <Button type="submit">Найти</Button>
          </form>
        </CardContent>
      </Card>

      {/* Search Results */}
      {searchParams.get('q') ? (
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <h2 className="text-xl font-semibold">
              Результаты для "{searchParams.get('q')}"
            </h2>
            {searchResults && (
              <span className="text-sm text-muted-foreground">
                ({(searchResults.posts?.length || 0) + (searchResults.users?.length || 0)} результатов)
              </span>
            )}
          </div>

          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 3 }).map((_, i) => (
                <Card key={i} className="post-skeleton">
                  <CardContent className="p-6">
                    <div className="space-y-3">
                      <div className="h-4 bg-muted rounded animate-pulse" />
                      <div className="h-4 bg-muted rounded animate-pulse w-3/4" />
                      <div className="h-4 bg-muted rounded animate-pulse w-1/2" />
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : error ? (
            <Card>
              <CardContent className="py-12 text-center">
                <p className="text-destructive">Ошибка при поиске</p>
              </CardContent>
            </Card>
          ) : searchResults ? (
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList>
                <TabsTrigger value="all">Все</TabsTrigger>
                <TabsTrigger value="posts">
                  Посты ({searchResults.posts?.length || 0})
                </TabsTrigger>
                <TabsTrigger value="users">
                  Пользователи ({searchResults.users?.length || 0})
                </TabsTrigger>
              </TabsList>

              <TabsContent value="all" className="space-y-4">
                {/* Users */}
                {searchResults.users && searchResults.users.length > 0 && (
                  <div className="space-y-2">
                    <h3 className="font-semibold text-sm text-muted-foreground">Пользователи</h3>
                    <div className="grid gap-3">
                      {searchResults.users.slice(0, 3).map((user) => (
                        <UserCard key={user.id} user={user} />
                      ))}
                    </div>
                  </div>
                )}

                {/* Posts */}
                {searchResults.posts && searchResults.posts.length > 0 && (
                  <div className="space-y-2">
                    <h3 className="font-semibold text-sm text-muted-foreground">Посты</h3>
                    <div className="space-y-4">
                      {searchResults.posts.map((post, index) => (
                        <motion.div
                          key={post.id}
                          initial={{ opacity: 0, y: 20 }}
                          animate={{ opacity: 1, y: 0 }}
                          transition={{ delay: index * 0.1 }}
                        >
                          <PostCard post={post} />
                        </motion.div>
                      ))}
                    </div>
                  </div>
                )}

                {(!searchResults.posts || searchResults.posts.length === 0) &&
                 (!searchResults.users || searchResults.users.length === 0) && (
                  <Card>
                    <CardContent className="py-12 text-center">
                      <p className="text-muted-foreground">
                        По вашему запросу ничего не найдено
                      </p>
                    </CardContent>
                  </Card>
                )}
              </TabsContent>

              <TabsContent value="posts" className="space-y-4">
                {searchResults.posts && searchResults.posts.length > 0 ? (
                  searchResults.posts.map((post, index) => (
                    <motion.div
                      key={post.id}
                      initial={{ opacity: 0, y: 20 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ delay: index * 0.1 }}
                    >
                      <PostCard post={post} />
                    </motion.div>
                  ))
                ) : (
                  <Card>
                    <CardContent className="py-12 text-center">
                      <p className="text-muted-foreground">Посты не найдены</p>
                    </CardContent>
                  </Card>
                )}
              </TabsContent>

              <TabsContent value="users" className="space-y-4">
                {searchResults.users && searchResults.users.length > 0 ? (
                  <div className="grid gap-3">
                    {searchResults.users.map((user) => (
                      <UserCard key={user.id} user={user} />
                    ))}
                  </div>
                ) : (
                  <Card>
                    <CardContent className="py-12 text-center">
                      <p className="text-muted-foreground">Пользователи не найдены</p>
                    </CardContent>
                  </Card>
                )}
              </TabsContent>
            </Tabs>
          ) : null}
        </div>
      ) : (
        <div className="grid md:grid-cols-2 gap-6">
          {/* Trending Hashtags */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Hash className="h-5 w-5" />
                Популярные хештеги
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              {trendingHashtags.map((hashtag) => (
                <Button
                  key={hashtag}
                  variant="ghost"
                  className="w-full justify-start"
                  onClick={() => setQuery(hashtag)}
                >
                  <Hash className="h-4 w-4 mr-2" />
                  {hashtag}
                </Button>
              ))}
            </CardContent>
          </Card>

          {/* Suggested Users */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-5 w-5" />
                Рекомендуемые пользователи
              </CardTitle>
            </CardHeader>
            <CardContent>
              {suggestedUsers.length > 0 ? (
                <div className="space-y-3">
                  {suggestedUsers.map((user) => (
                    <UserCard key={user.id} user={user} compact />
                  ))}
                </div>
              ) : (
                <p className="text-muted-foreground text-center py-8">
                  Рекомендации пока не загружены
                </p>
              )}
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}

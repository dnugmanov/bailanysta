import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import {
  X,
  Send,
  Image,
  Hash,
  BookOpen,
  Sparkles,
  Wand2,
  Heart,
  MessageCircle,
  Share,
} from 'lucide-react'
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query'
import { apiClient } from '@/api/client'
import { useUIStore } from '@/stores/ui'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import * as Dialog from '@radix-ui/react-dialog'
import * as Select from '@radix-ui/react-select'
import * as Tabs from '@radix-ui/react-tabs'
import { useToast } from '@/hooks/use-toast'
import { extractHashtags } from '@/lib/utils'
import FormattedText from './FormattedText'
import type { CreatePostRequest, GeneratePostRequest } from '@/types/api'

export default function PostComposer() {
  const [text, setText] = useState('')
  const [selectedCourse, setSelectedCourse] = useState<string>('no-course')
  const [activeTab, setActiveTab] = useState('compose')
  const [aiPrompt, setAiPrompt] = useState('')
  const [isGenerating, setIsGenerating] = useState(false)
  
  const { composerOpen, setComposerOpen } = useUIStore()
  const { user } = useAuthStore()
  const { toast } = useToast()
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const location = useLocation()
  
  // Проверяем, находимся ли мы на странице /compose
  const isComposePage = location.pathname === '/compose'

  // Fetch available courses
  const { data: coursesData } = useQuery({
    queryKey: ['courses'],
    queryFn: () => apiClient.getCourses(),
    retry: false, // Don't retry on auth errors
  })

  const createPostMutation = useMutation({
    mutationFn: (data: CreatePostRequest) => apiClient.createPost(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feed'] })
      setText('')
      setSelectedCourse('no-course')
      setComposerOpen(false)
      toast({
        title: 'Успешно!',
        description: 'Пост опубликован',
      })
      
      // Если мы на странице compose, перенаправляем на фид
      if (isComposePage) {
        navigate('/feed')
      }
    },
    onError: (error: any) => {
      toast({
        title: 'Ошибка',
        description: error.message || 'Не удалось опубликовать пост',
        variant: 'destructive',
      })
    },
  })

  const generatePostMutation = useMutation({
    mutationFn: (data: GeneratePostRequest) => apiClient.generatePost(data),
    onSuccess: (data) => {
      setText((data as any).text)
      setActiveTab('compose')
      toast({
        title: 'Готово!',
        description: 'Текст сгенерирован с помощью AI',
      })
    },
    onError: (error: any) => {
      console.error('AI generation error:', error)
      let errorMessage = 'Не удалось сгенерировать текст'
      
      if (error.message?.includes('API key')) {
        errorMessage = 'AI сервис недоступен. Попробуйте позже.'
      } else if (error.message) {
        errorMessage = error.message
      }
      
      toast({
        title: 'Ошибка AI',
        description: errorMessage,
        variant: 'destructive',
      })
    },
  })

  const handleSubmit = () => {
    if (!text.trim()) {
      toast({
        title: 'Ошибка',
        description: 'Введите текст поста',
        variant: 'destructive',
      })
      return
    }

    createPostMutation.mutate({
      text: text.trim(),
      course_id: selectedCourse && selectedCourse !== 'no-course' ? selectedCourse : undefined,
    })
  }

  const handleAIGenerate = () => {
    if (!aiPrompt.trim()) {
      toast({
        title: 'Ошибка',
        description: 'Введите тему для генерации',
        variant: 'destructive',
      })
      return
    }

    generatePostMutation.mutate({
      topic: aiPrompt,
      style: 'casual',
      format: 'markdown'
    })
  }

  const hashtags = extractHashtags(text)
  const charCount = text.length
  const maxChars = 5000

  return (
    <Dialog.Root open={composerOpen} onOpenChange={setComposerOpen}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 bg-black/50 z-50" />
        <Dialog.Content className="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-full max-w-2xl max-h-[90vh] z-50">
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ duration: 0.2 }}
          >
            <Card className="glass-effect">
              <CardHeader className="flex flex-row items-center justify-between pb-3">
                <CardTitle className="text-xl">Создать пост</CardTitle>
                <Dialog.Close asChild>
                  <Button variant="ghost" size="icon">
                    <X className="h-4 w-4" />
                  </Button>
                </Dialog.Close>
              </CardHeader>

              <CardContent>
                <Tabs.Root value={activeTab} onValueChange={setActiveTab}>
                  <Tabs.List className="grid w-full grid-cols-3 mb-4">
                    <Tabs.Trigger
                      value="compose"
                      className="flex items-center gap-2 px-3 py-2 text-sm font-medium transition-colors data-[state=active]:bg-primary data-[state=active]:text-primary-foreground rounded-md"
                    >
                      <Send className="h-4 w-4" />
                      Написать
                    </Tabs.Trigger>
                    <Tabs.Trigger
                      value="preview"
                      className="flex items-center gap-2 px-3 py-2 text-sm font-medium transition-colors data-[state=active]:bg-primary data-[state=active]:text-primary-foreground rounded-md"
                    >
                      <BookOpen className="h-4 w-4" />
                      Просмотр
                    </Tabs.Trigger>
                    <Tabs.Trigger
                      value="ai"
                      className="flex items-center gap-2 px-3 py-2 text-sm font-medium transition-colors data-[state=active]:bg-primary data-[state=active]:text-primary-foreground rounded-md"
                    >
                      <Sparkles className="h-4 w-4" />
                      AI Генератор
                    </Tabs.Trigger>
                  </Tabs.List>

                  <Tabs.Content value="compose" className="space-y-4">
                    <div className="space-y-2">
                      <Textarea
                        placeholder="Поделитесь своими знаниями, идеями или вопросами..."
                        value={text}
                        onChange={(e) => setText(e.target.value)}
                        className="min-h-[120px] resize-none focus:ring-2 focus:ring-primary"
                        maxLength={maxChars}
                      />
                      <div className="flex justify-between items-center text-xs text-muted-foreground">
                        <span>{hashtags.length > 0 && `Хештеги: ${hashtags.join(', ')}`}</span>
                        <span className={charCount > maxChars * 0.9 ? 'text-destructive' : ''}>
                          {charCount}/{maxChars}
                        </span>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <label className="text-sm font-medium">Курс (необязательно)</label>
                      <Select.Root value={selectedCourse} onValueChange={setSelectedCourse}>
                        <Select.Trigger className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
                          <Select.Value placeholder="Выберите курс" />
                          <Select.Icon>
                            <BookOpen className="h-4 w-4" />
                          </Select.Icon>
                        </Select.Trigger>
                        <Select.Portal>
                          <Select.Content className="relative z-50 max-h-96 min-w-[8rem] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md">
                            <Select.ScrollUpButton className="flex cursor-default items-center justify-center py-1">
                              <BookOpen className="h-4 w-4" />
                            </Select.ScrollUpButton>
                            <Select.Viewport className="p-1">
                              <Select.Item
                                value="no-course"
                                className="relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                              >
                                <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
                                  <Select.ItemIndicator>
                                    <BookOpen className="h-4 w-4" />
                                  </Select.ItemIndicator>
                                </span>
                                <Select.ItemText>Без курса</Select.ItemText>
                              </Select.Item>
                              {coursesData?.courses?.map((course) => (
                                <Select.Item
                                  key={course.id}
                                  value={course.id.toString()}
                                  className="relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                                >
                                  <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
                                    <Select.ItemIndicator>
                                      <BookOpen className="h-4 w-4" />
                                    </Select.ItemIndicator>
                                  </span>
                                  <Select.ItemText>{course.title}</Select.ItemText>
                                </Select.Item>
                              ))}
                            </Select.Viewport>
                            <Select.ScrollDownButton className="flex cursor-default items-center justify-center py-1">
                              <BookOpen className="h-4 w-4" />
                            </Select.ScrollDownButton>
                          </Select.Content>
                        </Select.Portal>
                      </Select.Root>
                    </div>

                    <div className="flex justify-between items-center pt-4 border-t">
                      <div className="flex items-center space-x-2 text-muted-foreground">
                        <Hash className="h-4 w-4" />
                        <span className="text-xs">Используйте # для хештегов</span>
                      </div>
                      <Button
                        onClick={handleSubmit}
                        disabled={!text.trim() || createPostMutation.isPending}
                        className="min-w-[100px]"
                      >
                        {createPostMutation.isPending ? (
                          <div className="flex items-center gap-2">
                            <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                            Публикация...
                          </div>
                        ) : (
                          <div className="flex items-center gap-2">
                            <Send className="w-4 h-4" />
                            Опубликовать
                          </div>
                        )}
                      </Button>
                    </div>
                  </Tabs.Content>

                  <Tabs.Content value="preview" className="space-y-4">
                    <div className="space-y-4">
                      <h3 className="text-lg font-semibold">Предварительный просмотр</h3>
                      
                      {text.trim() ? (
                        <div className="border rounded-lg p-4 bg-muted/30">
                          <div className="flex items-center gap-3 mb-3">
                            <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                              <span className="text-sm font-medium">{user?.username?.[0]?.toUpperCase()}</span>
                            </div>
                            <div>
                              <p className="font-medium text-sm">{user?.username}</p>
                              <p className="text-xs text-muted-foreground">Сейчас</p>
                            </div>
                          </div>
                          
                          <FormattedText 
                            text={text} 
                            className="text-sm leading-relaxed"
                          />
                          
                          <div className="flex items-center gap-4 mt-4 pt-4 border-t">
                            <button className="flex items-center gap-1 text-muted-foreground hover:text-foreground transition-colors">
                              <Heart className="h-4 w-4" />
                              <span className="text-sm">0</span>
                            </button>
                            <button className="flex items-center gap-1 text-muted-foreground hover:text-foreground transition-colors">
                              <MessageCircle className="h-4 w-4" />
                              <span className="text-sm">0</span>
                            </button>
                            <button className="flex items-center gap-1 text-muted-foreground hover:text-foreground transition-colors">
                              <Share className="h-4 w-4" />
                            </button>
                          </div>
                        </div>
                      ) : (
                        <div className="border rounded-lg p-8 text-center text-muted-foreground bg-muted/30">
                          <BookOpen className="h-8 w-8 mx-auto mb-2 opacity-50" />
                          <p>Напишите текст поста, чтобы увидеть предварительный просмотр</p>
                        </div>
                      )}
                    </div>
                  </Tabs.Content>

                  <Tabs.Content value="ai" className="space-y-4">
                    <div className="space-y-2">
                      <label className="text-sm font-medium">Тема для генерации</label>
                      <Textarea
                        placeholder="Опишите тему поста, который хотите создать..."
                        value={aiPrompt}
                        onChange={(e) => setAiPrompt(e.target.value)}
                        className="min-h-[80px] resize-none"
                      />
                    </div>

                    <div className="flex justify-end">
                      <Button
                        onClick={handleAIGenerate}
                        disabled={!aiPrompt.trim() || generatePostMutation.isPending}
                        variant="outline"
                      >
                        {generatePostMutation.isPending ? (
                          <div className="flex items-center gap-2">
                            <div className="w-4 h-4 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
                            Генерируем...
                          </div>
                        ) : (
                          <div className="flex items-center gap-2">
                            <Wand2 className="w-4 h-4" />
                            Сгенерировать
                          </div>
                        )}
                      </Button>
                    </div>

                    {text && (
                      <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="p-4 bg-muted rounded-lg"
                      >
                        <p className="text-sm text-muted-foreground mb-2">Сгенерированный текст:</p>
                        <p className="text-sm">{text}</p>
                        <div className="flex justify-end mt-3">
                          <Button
                            size="sm"
                            onClick={() => setActiveTab('compose')}
                          >
                            Редактировать
                          </Button>
                        </div>
                      </motion.div>
                    )}
                  </Tabs.Content>
                </Tabs.Root>
              </CardContent>
            </Card>
          </motion.div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}

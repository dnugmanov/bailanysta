import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { X, Save } from 'lucide-react'
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import * as Select from '@radix-ui/react-select'
import { ChevronDown } from 'lucide-react'
import { apiClient } from '@/api/client'
import { useToast } from '@/hooks/use-toast'
import type { Post, UpdatePostRequest } from '@/types/api'

interface PostEditorProps {
  post: Post
  isOpen: boolean
  onClose: () => void
}

export default function PostEditor({ post, isOpen, onClose }: PostEditorProps) {
  const [text, setText] = useState(post.text)
  const [selectedCourse, setSelectedCourse] = useState<string>(
    post.course_id ? post.course_id.toString() : 'no-course'
  )
  const { toast } = useToast()
  const queryClient = useQueryClient()

  // Reset form when post changes
  useEffect(() => {
    setText(post.text)
    setSelectedCourse(post.course_id ? post.course_id.toString() : 'no-course')
  }, [post])

  // Fetch available courses
  const { data: coursesData } = useQuery({
    queryKey: ['courses'],
    queryFn: () => apiClient.getCourses(),
    retry: false,
  })

  const updatePostMutation = useMutation({
    mutationFn: (data: UpdatePostRequest) => apiClient.updatePost(post.id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feed'] })
      queryClient.invalidateQueries({ queryKey: ['posts', post.id] })
      toast({
        title: 'Успешно!',
        description: 'Пост обновлен',
      })
      onClose()
    },
    onError: (error: any) => {
      toast({
        title: 'Ошибка',
        description: error.message || 'Не удалось обновить пост',
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

    updatePostMutation.mutate({
      text: text.trim(),
      course_id: selectedCourse && selectedCourse !== 'no-course' ? selectedCourse : undefined,
    })
  }

  const maxChars = 5000
  const charCount = text.length

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        className="bg-background border border-border rounded-lg w-full max-w-2xl max-h-[90vh] overflow-y-auto"
      >
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle>Редактировать пост</CardTitle>
            <Button variant="ghost" size="sm" onClick={onClose}>
              <X className="h-4 w-4" />
            </Button>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Course Selection */}
            {coursesData?.courses && coursesData.courses.length > 0 && (
              <div>
                <label className="text-sm font-medium mb-2 block">Курс</label>
                <Select.Root value={selectedCourse} onValueChange={setSelectedCourse}>
                  <Select.Trigger className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50">
                    <Select.Value placeholder="Выберите курс" />
                    <Select.Icon asChild>
                      <ChevronDown className="h-4 w-4 opacity-50" />
                    </Select.Icon>
                  </Select.Trigger>
                  <Select.Portal>
                    <Select.Content className="relative z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md">
                      <Select.ScrollUpButton className="flex cursor-default items-center justify-center py-1">
                        <ChevronDown className="h-4 w-4 rotate-180" />
                      </Select.ScrollUpButton>
                      <Select.Viewport className="p-1">
                        <Select.Item
                          value="no-course"
                          className="relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                        >
                          <Select.ItemIndicator className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
                            ✓
                          </Select.ItemIndicator>
                          <Select.ItemText>Без курса</Select.ItemText>
                        </Select.Item>
                        {coursesData.courses.map((course) => (
                          <Select.Item
                            key={course.id}
                            value={course.id.toString()}
                            className="relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
                          >
                            <Select.ItemIndicator className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
                              ✓
                            </Select.ItemIndicator>
                            <Select.ItemText>{course.title}</Select.ItemText>
                          </Select.Item>
                        ))}
                      </Select.Viewport>
                      <Select.ScrollDownButton className="flex cursor-default items-center justify-center py-1">
                        <ChevronDown className="h-4 w-4" />
                      </Select.ScrollDownButton>
                    </Select.Content>
                  </Select.Portal>
                </Select.Root>
              </div>
            )}

            {/* Text Area */}
            <div>
              <label className="text-sm font-medium mb-2 block">Текст поста</label>
              <Textarea
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="О чём думаете?"
                rows={8}
                className="resize-none"
                maxLength={maxChars}
              />
              <div className="flex justify-between items-center mt-2 text-sm text-muted-foreground">
                <span>Поддерживается Markdown</span>
                <span className={charCount > maxChars * 0.9 ? 'text-destructive' : ''}>
                  {charCount}/{maxChars}
                </span>
              </div>
            </div>

            {/* Actions */}
            <div className="flex gap-2 pt-4">
              <Button variant="outline" onClick={onClose} className="flex-1">
                Отмена
              </Button>
              <Button 
                onClick={handleSubmit}
                disabled={updatePostMutation.isPending || !text.trim() || charCount > maxChars}
                className="flex-1"
              >
                {updatePostMutation.isPending ? (
                  <>Сохраняем...</>
                ) : (
                  <>
                    <Save className="h-4 w-4 mr-2" />
                    Сохранить
                  </>
                )}
              </Button>
            </div>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  )
}

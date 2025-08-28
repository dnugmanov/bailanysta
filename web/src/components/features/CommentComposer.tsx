import { useState } from 'react'
import { motion } from 'framer-motion'
import { Send, Smile } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { useToast } from '@/hooks/use-toast'
import { getInitials } from '@/lib/utils'

interface CommentComposerProps {
  postId: string
}

export default function CommentComposer({ postId }: CommentComposerProps) {
  const [text, setText] = useState('')
  const { user } = useAuthStore()
  const { toast } = useToast()
  const queryClient = useQueryClient()

  const createCommentMutation = useMutation({
    mutationFn: (text: string) => apiClient.createComment(postId, { text }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', postId] })
      queryClient.invalidateQueries({ queryKey: ['post', postId] }) // Update comment count
      setText('')
      toast({
        title: 'Успешно!',
        description: 'Комментарий добавлен',
      })
    },
    onError: (error: any) => {
      toast({
        title: 'Ошибка',
        description: error.message || 'Не удалось добавить комментарий',
        variant: 'destructive',
      })
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!text.trim()) {
      toast({
        title: 'Ошибка',
        description: 'Введите текст комментария',
        variant: 'destructive',
      })
      return
    }

    createCommentMutation.mutate(text.trim())
  }

  const maxChars = 500
  const charCount = text.length

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="flex gap-3">
        <Avatar className="h-8 w-8 flex-shrink-0">
          <AvatarImage src={user?.avatar_url || undefined} />
          <AvatarFallback className="text-xs">
            {user ? getInitials(user.username) : 'U'}
          </AvatarFallback>
        </Avatar>

        <div className="flex-1 space-y-2">
          <Textarea
            placeholder="Напишите комментарий..."
            value={text}
            onChange={(e) => setText(e.target.value)}
            className="min-h-[80px] resize-none focus:ring-2 focus:ring-primary"
            maxLength={maxChars}
          />
          
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 text-muted-foreground">
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="h-6 px-2"
                disabled
              >
                <Smile className="h-3 w-3" />
              </Button>
              <span className="text-xs">
                Эмодзи скоро будут доступны
              </span>
            </div>

            <div className="flex items-center gap-3">
              <span 
                className={`text-xs ${
                  charCount > maxChars * 0.9 ? 'text-destructive' : 'text-muted-foreground'
                }`}
              >
                {charCount}/{maxChars}
              </span>
              
              <Button
                type="submit"
                size="sm"
                disabled={!text.trim() || createCommentMutation.isPending}
              >
                {createCommentMutation.isPending ? (
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Отправка...
                  </div>
                ) : (
                  <div className="flex items-center gap-2">
                    <Send className="w-3 h-3" />
                    Отправить
                  </div>
                )}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </form>
  )
}

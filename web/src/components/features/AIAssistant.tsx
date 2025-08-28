import { useState } from 'react'
import { motion } from 'framer-motion'
import {
  Sparkles,
  FileText,
  MessageCircle,
  HelpCircle,
  Wand2,
  BookOpen,
  Brain,
  X,
} from 'lucide-react'
import { useMutation } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import * as Tabs from '@radix-ui/react-tabs'
import { apiClient } from '@/api/client'
import { useToast } from '@/hooks/use-toast'
import FormattedText from './FormattedText'

interface AIAssistantProps {
  isOpen: boolean
  onClose: () => void
}

export default function AIAssistant({ isOpen, onClose }: AIAssistantProps) {
  const [activeTab, setActiveTab] = useState('study-notes')
  const [prompt, setPrompt] = useState('')
  const [context, setContext] = useState('')
  const [result, setResult] = useState('')
  const { toast } = useToast()

  const generateStudyNotesMutation = useMutation({
    mutationFn: (data: { topic: string; course_context?: string }) =>
      apiClient.generateStudyNotes(data),
    onSuccess: (data) => {
      setResult((data as any).text)
      toast({
        title: 'Готово!',
        description: 'Учебные заметки созданы',
      })
    },
    onError: (error: any) => {
      toast({
        title: 'Ошибка AI',
        description: error.message || 'Не удалось создать заметки',
        variant: 'destructive',
      })
    },
  })



  const explainConceptMutation = useMutation({
    mutationFn: (data: { concept: string; level?: string; context?: string }) =>
      apiClient.explainConcept(data),
    onSuccess: (data) => {
      setResult((data as any).text)
      toast({
        title: 'Готово!',
        description: 'Концепт объяснен',
      })
    },
    onError: (error: any) => {
      toast({
        title: 'Ошибка AI',
        description: error.message || 'Не удалось объяснить концепт',
        variant: 'destructive',
      })
    },
  })

  const handleGenerate = () => {
    if (!prompt.trim()) {
      toast({
        title: 'Ошибка',
        description: 'Введите запрос',
        variant: 'destructive',
      })
      return
    }

    switch (activeTab) {
      case 'study-notes':
        generateStudyNotesMutation.mutate({
          topic: prompt,
          course_context: context || undefined,
        })
        break
      case 'explain':
        explainConceptMutation.mutate({
          concept: prompt,
          level: 'intermediate',
          context: context || undefined,
        })
        break
    }
  }

  const isLoading = generateStudyNotesMutation.isPending || 
                   explainConceptMutation.isPending

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        className="bg-background border border-border rounded-lg w-full max-w-4xl max-h-[90vh] overflow-y-auto"
      >
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-primary" />
              AI Помощник
            </CardTitle>
            <Button variant="ghost" size="sm" onClick={onClose}>
              <X className="h-4 w-4" />
            </Button>
          </CardHeader>
          <CardContent>
            <Tabs.Root value={activeTab} onValueChange={setActiveTab}>
              <Tabs.List className="grid grid-cols-2 mb-6">
                <Tabs.Trigger
                  value="study-notes"
                  className="flex items-center gap-2 data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
                >
                  <BookOpen className="h-4 w-4" />
                  Конспект
                </Tabs.Trigger>
                <Tabs.Trigger
                  value="explain"
                  className="flex items-center gap-2 data-[state=active]:bg-primary data-[state=active]:text-primary-foreground"
                >
                  <Brain className="h-4 w-4" />
                  Объяснение
                </Tabs.Trigger>
              </Tabs.List>

              <Tabs.Content value="study-notes" className="space-y-4">
                <div>
                  <label className="text-sm font-medium mb-2 block">
                    Тема для изучения
                  </label>
                  <Input
                    value={prompt}
                    onChange={(e) => setPrompt(e.target.value)}
                    placeholder="Например: Основы машинного обучения"
                  />
                </div>
                <div>
                  <label className="text-sm font-medium mb-2 block">
                    Контекст курса (опционально)
                  </label>
                  <Input
                    value={context}
                    onChange={(e) => setContext(e.target.value)}
                    placeholder="Например: Python для анализа данных"
                  />
                </div>
              </Tabs.Content>

              <Tabs.Content value="explain" className="space-y-4">
                <div>
                  <label className="text-sm font-medium mb-2 block">
                    Концепт для объяснения
                  </label>
                  <Textarea
                    value={prompt}
                    onChange={(e) => setPrompt(e.target.value)}
                    placeholder="Например: Рекурсия в программировании, Асинхронность в JavaScript"
                    rows={3}
                  />
                </div>
                <div>
                  <label className="text-sm font-medium mb-2 block">
                    Дополнительный контекст (опционально)
                  </label>
                  <Textarea
                    value={context}
                    onChange={(e) => setContext(e.target.value)}
                    placeholder="Например: Для начинающих программистов, углубленное объяснение для продвинутых"
                    rows={2}
                  />
                </div>
              </Tabs.Content>

              <Button
                onClick={handleGenerate}
                disabled={isLoading || !prompt.trim()}
                className="w-full mb-6"
              >
                {isLoading ? (
                  <>
                    <Wand2 className="h-4 w-4 mr-2 animate-spin" />
                    Генерируем...
                  </>
                ) : (
                  <>
                    <Sparkles className="h-4 w-4 mr-2" />
                    Сгенерировать
                  </>
                )}
              </Button>

              {result && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">Результат:</label>
                  <Card>
                    <CardContent className="p-6 max-h-96 overflow-y-auto">
                      <FormattedText text={result} />
                    </CardContent>
                  </Card>
                </div>
              )}
            </Tabs.Root>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  )
}

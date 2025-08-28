import { useEffect } from 'react'
import { useUIStore } from '@/stores/ui'
import PostComposer from '@/components/features/PostComposer'

export default function ComposePage() {
  const { setComposerOpen } = useUIStore()

  useEffect(() => {
    setComposerOpen(true)
  }, [setComposerOpen])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Создать пост</h1>
        <p className="text-muted-foreground">
          Поделитесь своими знаниями с сообществом
        </p>
      </div>
      <PostComposer />
    </div>
  )
}

import React from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'
import { cn } from '@/lib/utils'

interface FormattedTextProps {
  text: string
  className?: string
  maxLines?: number
}

export default function FormattedText({ text, className, maxLines }: FormattedTextProps) {
  return (
    <div
      className={cn(
        'text-sm leading-relaxed break-words prose prose-sm max-w-none',
        'prose-headings:text-foreground prose-p:text-foreground prose-p:my-2',
        'prose-strong:text-foreground prose-code:text-foreground',
        'prose-blockquote:text-muted-foreground prose-blockquote:border-l-border',
        'prose-code:bg-muted prose-code:px-1 prose-code:py-0.5 prose-code:rounded prose-code:before:content-none prose-code:after:content-none',
        'prose-pre:bg-muted prose-pre:text-foreground prose-pre:border prose-pre:border-border',
        'prose-a:text-primary prose-a:no-underline hover:prose-a:underline',
        'prose-ul:my-2 prose-ol:my-2 prose-li:my-0',
        maxLines && `line-clamp-${maxLines}`,
        className
      )}
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight]}
      >
        {text}
      </ReactMarkdown>
    </div>
  )
}

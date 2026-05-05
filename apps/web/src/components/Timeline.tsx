import { ReactNode } from 'react'

interface TimelineItem {
  id: string
  time?: string
  title: string
  description?: string
  icon?: ReactNode
}

interface TimelineProps {
  items: TimelineItem[]
  className?: string
}

export function Timeline({ items, className = '' }: TimelineProps) {
  return (
    <div className={`relative ${className}`}>
      <div className="absolute left-3 top-0 h-full w-px bg-gray-200 dark:bg-gray-700" />
      <ul className="space-y-4">
        {items.map((item) => (
          <li key={item.id} className="relative pl-8">
            <div className="absolute left-1.5 top-1.5 h-3 w-3 rounded-full border-2 border-blue-500 bg-white dark:bg-gray-900" />
            <div className="flex flex-col">
              {item.time && <span className="text-xs text-gray-400">{item.time}</span>}
              <span className="text-sm font-medium text-gray-900 dark:text-white">
                {item.title}
              </span>
              {item.description && (
                <span className="text-xs text-gray-500 dark:text-gray-400">{item.description}</span>
              )}
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}

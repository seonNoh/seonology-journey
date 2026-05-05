import { X } from 'lucide-react'

interface ChipProps {
  label: string
  onRemove?: () => void
  variant?: 'filled' | 'outlined'
}

export function Chip({ label, onRemove, variant = 'filled' }: ChipProps) {
  const base =
    variant === 'filled'
      ? 'bg-sakura-100 text-sakura-800'
      : 'border border-sakura-300 text-sakura-700'

  return (
    <span className={`inline-flex items-center gap-1 rounded-full px-3 py-1 text-sm ${base}`}>
      {label}
      {onRemove && (
        <button
          onClick={onRemove}
          className="ml-0.5 rounded-full p-0.5 hover:bg-sakura-200"
          aria-label={`Remove ${label}`}
        >
          <X className="h-3 w-3" />
        </button>
      )}
    </span>
  )
}

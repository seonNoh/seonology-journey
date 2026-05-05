import { useState } from 'react'
import { Heart } from 'lucide-react'

interface RatingStarsProps {
  value?: number
  max?: number
  onChange?: (rating: number) => void
  readOnly?: boolean
  className?: string
}

export function RatingStars({
  value = 0,
  max = 5,
  onChange,
  readOnly = false,
  className = '',
}: RatingStarsProps) {
  const [hovered, setHovered] = useState(0)

  const handleClick = (rating: number) => {
    if (!readOnly && onChange) {
      onChange(rating === value ? 0 : rating)
    }
  }

  return (
    <div className={`flex items-center gap-0.5 ${className}`}>
      {Array.from({ length: max }, (_, i) => {
        const rating = i + 1
        const filled = rating <= (hovered || value)
        return (
          <button
            key={rating}
            type="button"
            disabled={readOnly}
            onClick={() => handleClick(rating)}
            onMouseEnter={() => !readOnly && setHovered(rating)}
            onMouseLeave={() => !readOnly && setHovered(0)}
            className="p-0.5 disabled:cursor-default"
            aria-label={`${rating} of ${max}`}
          >
            <Heart
              className={`h-5 w-5 transition-colors ${
                filled
                  ? 'fill-rose-500 text-rose-500'
                  : 'fill-none text-gray-300 dark:text-gray-600'
              }`}
            />
          </button>
        )
      })}
    </div>
  )
}

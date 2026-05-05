import { useState } from 'react'
import { Clock } from 'lucide-react'

interface TimePickerProps {
  value?: string
  onChange: (time: string) => void
  step?: number
  className?: string
}

export function TimePicker({ value = '', onChange, step = 15, className = '' }: TimePickerProps) {
  const [time, setTime] = useState(value)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value
    setTime(val)
    onChange(val)
  }

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <Clock className="h-4 w-4 text-gray-400 shrink-0" />
      <input
        type="time"
        value={time}
        onChange={handleChange}
        step={step * 60}
        className="rounded-md border border-gray-300 px-2 py-1 text-sm focus:border-blue-500 focus:outline-none dark:border-gray-600 dark:bg-gray-800 dark:text-white"
      />
    </div>
  )
}

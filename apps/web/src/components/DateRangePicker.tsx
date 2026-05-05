import { useState } from 'react'
import { Calendar } from 'lucide-react'

interface DateRangePickerProps {
  startDate?: string
  endDate?: string
  onChange: (range: { start: string; end: string }) => void
  placeholder?: string
  className?: string
}

export function DateRangePicker({
  startDate = '',
  endDate = '',
  onChange,
  placeholder = 'Select dates',
  className = '',
}: DateRangePickerProps) {
  const [start, setStart] = useState(startDate)
  const [end, setEnd] = useState(endDate)

  const handleStartChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value
    setStart(val)
    if (val && end && val <= end) {
      onChange({ start: val, end })
    }
  }

  const handleEndChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value
    setEnd(val)
    if (start && val && start <= val) {
      onChange({ start, end: val })
    }
  }

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <Calendar className="h-4 w-4 text-gray-400 shrink-0" />
      <input
        type="date"
        value={start}
        onChange={handleStartChange}
        placeholder={placeholder}
        className="rounded-md border border-gray-300 px-2 py-1 text-sm focus:border-blue-500 focus:outline-none dark:border-gray-600 dark:bg-gray-800 dark:text-white"
      />
      <span className="text-gray-400">-</span>
      <input
        type="date"
        value={end}
        onChange={handleEndChange}
        min={start}
        className="rounded-md border border-gray-300 px-2 py-1 text-sm focus:border-blue-500 focus:outline-none dark:border-gray-600 dark:bg-gray-800 dark:text-white"
      />
    </div>
  )
}

import { useState } from 'react'

interface CurrencyInputProps {
  value?: number
  currency?: string
  onChange: (amount: number) => void
  currencies?: string[]
  onCurrencyChange?: (currency: string) => void
  className?: string
}

export function CurrencyInput({
  value = 0,
  currency = 'JPY',
  onChange,
  currencies = ['JPY', 'KRW', 'USD', 'EUR', 'THB', 'VND'],
  onCurrencyChange,
  className = '',
}: CurrencyInputProps) {
  const [amount, setAmount] = useState(value.toString())

  const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value.replace(/[^0-9.]/g, '')
    setAmount(raw)
    const num = parseFloat(raw)
    if (!isNaN(num)) {
      onChange(num)
    }
  }

  return (
    <div className={`flex items-center gap-1 ${className}`}>
      <select
        value={currency}
        onChange={(e) => onCurrencyChange?.(e.target.value)}
        className="rounded-l-md border border-gray-300 bg-gray-50 px-2 py-1 text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-white"
      >
        {currencies.map((c) => (
          <option key={c} value={c}>
            {c}
          </option>
        ))}
      </select>
      <input
        type="text"
        inputMode="decimal"
        value={amount}
        onChange={handleAmountChange}
        placeholder="0"
        className="w-full rounded-r-md border border-l-0 border-gray-300 px-2 py-1 text-sm text-right focus:border-blue-500 focus:outline-none dark:border-gray-600 dark:bg-gray-800 dark:text-white"
      />
    </div>
  )
}

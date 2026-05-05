import { useState } from 'react'

interface Tab {
  id: string
  label: string
}

interface TabsProps {
  tabs: Tab[]
  defaultTab?: string
  onChange?: (id: string) => void
  children: (activeId: string) => React.ReactNode
}

export function Tabs({ tabs, defaultTab, onChange, children }: TabsProps) {
  const [active, setActive] = useState(defaultTab ?? tabs[0]?.id ?? '')

  const handleClick = (id: string) => {
    setActive(id)
    onChange?.(id)
  }

  return (
    <div>
      <div className="flex border-b border-slate-200" role="tablist">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            role="tab"
            aria-selected={active === tab.id}
            onClick={() => handleClick(tab.id)}
            className={`px-4 py-2 text-sm font-medium transition-colors ${
              active === tab.id
                ? 'border-b-2 border-sakura-500 text-sakura-700'
                : 'text-slate-500 hover:text-slate-700'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>
      <div className="pt-4" role="tabpanel">
        {children(active)}
      </div>
    </div>
  )
}

interface EmptyStateProps {
  title: string
  description?: string
  action?: React.ReactNode
  illustration?: React.ReactNode
}

export function EmptyState({ title, description, action, illustration }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <div className="mb-4">
        {illustration || (
          <div className="rounded-full bg-sakura-50 p-6">
            <svg
              width="48"
              height="48"
              viewBox="0 0 48 48"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
            >
              <path
                d="M24 4L6 14v20l18 10 18-10V14L24 4z"
                stroke="currentColor"
                strokeWidth="2"
                className="text-sakura-300"
                fill="none"
              />
              <path
                d="M24 24v20M6 14l18 10 18-10"
                stroke="currentColor"
                strokeWidth="2"
                className="text-sakura-300"
              />
            </svg>
          </div>
        )}
      </div>
      <h3 className="text-lg font-medium text-slate-900">{title}</h3>
      {description && <p className="mt-1 text-sm text-slate-500 max-w-sm">{description}</p>}
      {action && <div className="mt-4">{action}</div>}
    </div>
  )
}

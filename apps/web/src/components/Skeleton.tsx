interface SkeletonProps {
  className?: string
  variant?: 'text' | 'circle' | 'rect'
}

export function Skeleton({ className = '', variant = 'text' }: SkeletonProps) {
  const baseClasses: Record<string, string> = {
    text: 'h-4 w-full rounded',
    circle: 'h-10 w-10 rounded-full',
    rect: 'h-24 w-full rounded-lg',
  }

  return (
    <div
      className={`animate-pulse bg-slate-200 ${baseClasses[variant]} ${className}`}
      aria-hidden="true"
    />
  )
}

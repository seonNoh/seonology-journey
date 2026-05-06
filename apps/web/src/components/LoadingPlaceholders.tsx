import { Skeleton } from './Skeleton'
import { CuteLoader } from './CuteLoader'

/**
 * Reusable loading placeholders that match each list/detail shape in the
 * app. Each component renders the same layout as the filled state so the
 * page doesn't reflow when real data arrives, while still showing a small
 * CuteLoader for character.
 */

export function TripListSkeleton({ count = 4 }: { count?: number }) {
  return (
    <div className="space-y-3">
      <div className="flex items-center justify-start">
        <CuteLoader message="여행을 꺼내오고 있어요" />
      </div>
      <ul className="grid grid-cols-1 gap-3 md:grid-cols-2" aria-hidden="true">
        {Array.from({ length: count }).map((_, i) => (
          <li key={i} className="rounded-xl bg-white p-4 shadow-sm space-y-2">
            <Skeleton className="h-5 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
            <Skeleton className="h-4 w-2/3" />
          </li>
        ))}
      </ul>
    </div>
  )
}

export function TripDetailSkeleton() {
  return (
    <div className="space-y-3" aria-hidden="true">
      <CuteLoader message="여행을 펼치고 있어요" />
      <div className="rounded-2xl bg-white p-6 shadow-sm space-y-2">
        <Skeleton className="h-7 w-2/3" />
        <Skeleton className="h-4 w-1/3" />
        <Skeleton className="h-4 w-1/2" />
      </div>
    </div>
  )
}

export function ListRowsSkeleton({ count = 3, message }: { count?: number; message?: string }) {
  return (
    <div className="space-y-3">
      {message && <CuteLoader message={message} />}
      <ul className="space-y-2" aria-hidden="true">
        {Array.from({ length: count }).map((_, i) => (
          <li key={i} className="rounded-xl bg-white p-4 shadow-sm space-y-2">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-3 w-1/2" />
          </li>
        ))}
      </ul>
    </div>
  )
}

export function GridSkeleton({ count = 6, message }: { count?: number; message?: string }) {
  return (
    <div className="space-y-3">
      {message && <CuteLoader message={message} />}
      <ul className="grid grid-cols-2 gap-2 sm:grid-cols-3" aria-hidden="true">
        {Array.from({ length: count }).map((_, i) => (
          <li key={i} className="aspect-square rounded-xl bg-white shadow-sm overflow-hidden">
            <Skeleton className="h-full w-full rounded-none" />
          </li>
        ))}
      </ul>
    </div>
  )
}

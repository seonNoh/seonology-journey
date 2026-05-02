import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ArrowLeft, CalendarDays, MapPin, Wallet, StickyNote, ListChecks } from 'lucide-react'
import { api } from '../lib/api'
import type { Day, ListDaysResponse, Trip } from '../lib/types'

export function TripDetailPage() {
  const { tripId = '' } = useParams()
  const trip = useQuery({
    queryKey: ['trip', tripId],
    queryFn: () => api.get<{ trip: Trip }>(`/trips/${tripId}`),
    enabled: !!tripId,
  })
  const days = useQuery({
    queryKey: ['days', tripId],
    queryFn: () => api.get<ListDaysResponse>(`/trips/${tripId}/days`),
    enabled: !!tripId,
  })
  return (
    <section className="space-y-4">
      <Link to="/trips" className="inline-flex items-center gap-1 text-sm text-sakura-700">
        <ArrowLeft className="h-4 w-4" /> 목록으로
      </Link>
      {trip.isLoading && <p>불러오는 중…</p>}
      {trip.error && <p className="text-red-500">{(trip.error as Error).message}</p>}
      {trip.data && (
        <div className="rounded-2xl bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-bold text-slate-800">{trip.data.trip.title}</h1>
          {trip.data.trip.destination && (
            <p className="mt-1 flex items-center gap-1 text-slate-600">
              <MapPin className="h-4 w-4" /> {trip.data.trip.destination}
            </p>
          )}
          {trip.data.trip.startDate && (
            <p className="mt-1 flex items-center gap-1 text-slate-500">
              <CalendarDays className="h-4 w-4" /> {trip.data.trip.startDate} →{' '}
              {trip.data.trip.endDate}
            </p>
          )}
          {trip.data.trip.description && (
            <p className="mt-3 text-slate-600 whitespace-pre-wrap">{trip.data.trip.description}</p>
          )}
        </div>
      )}
      <div className="grid grid-cols-3 gap-2">
        <Link
          to={`/trips/${tripId}/expenses`}
          className="flex items-center justify-center gap-2 rounded-xl bg-white p-3 text-sm font-bold text-sakura-700 shadow-sm hover:shadow"
        >
          <Wallet className="h-4 w-4" /> 지출
        </Link>
        <Link
          to={`/trips/${tripId}/notes`}
          className="flex items-center justify-center gap-2 rounded-xl bg-white p-3 text-sm font-bold text-sakura-700 shadow-sm hover:shadow"
        >
          <StickyNote className="h-4 w-4" /> 메모
        </Link>
        <Link
          to={`/trips/${tripId}/checklist`}
          className="flex items-center justify-center gap-2 rounded-xl bg-white p-3 text-sm font-bold text-sakura-700 shadow-sm hover:shadow"
        >
          <ListChecks className="h-4 w-4" /> 체크
        </Link>
      </div>
      <h2 className="text-lg font-bold text-slate-800">일정</h2>
      <ul className="space-y-2">
        {days.data?.days?.map((d: Day) => (
          <li key={d.id}>
            <Link
              to={`/days/${d.id}`}
              className="flex items-center justify-between rounded-xl bg-white p-4 shadow-sm hover:shadow"
            >
              <div>
                <span className="text-xs font-bold text-sakura-600">DAY {d.dayNumber}</span>
                <p className="text-slate-800">
                  {d.date} {d.dayOfWeek && <span className="text-slate-400">({d.dayOfWeek})</span>}
                </p>
                {d.region && <p className="text-sm text-slate-500">{d.region}</p>}
              </div>
              <span className="text-slate-400">→</span>
            </Link>
          </li>
        ))}
        {days.data?.days?.length === 0 && (
          <li className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
            Day 가 없습니다. 시작일/종료일을 설정한 새 여행을 만들면 자동으로 생성됩니다.
          </li>
        )}
      </ul>
    </section>
  )
}

import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ArrowLeft, Clock } from 'lucide-react'
import { api } from '../lib/api'
import type { ListSchedulesResponse } from '../lib/types'

export function DayDetailPage() {
  const { dayId = '' } = useParams()
  const schedules = useQuery({
    queryKey: ['schedules', dayId],
    queryFn: () => api.get<ListSchedulesResponse>(`/days/${dayId}/schedules`),
    enabled: !!dayId,
  })
  return (
    <section className="space-y-4">
      <Link to="/trips" className="inline-flex items-center gap-1 text-sm text-sakura-700">
        <ArrowLeft className="h-4 w-4" /> 목록으로
      </Link>
      <h1 className="text-xl font-bold text-slate-800">일정 타임라인</h1>
      {schedules.isLoading && <p>불러오는 중…</p>}
      <ol className="relative border-l border-sakura-200 pl-4 space-y-3">
        {schedules.data?.schedules
          ?.slice()
          .sort((a, b) => (a.startTime ?? '').localeCompare(b.startTime ?? ''))
          .map((s) => (
            <li key={s.scheduleId} className="rounded-xl bg-white p-4 shadow-sm">
              <span className="absolute -left-1.5 mt-1 h-3 w-3 rounded-full bg-sakura-500" />
              <p className="flex items-center gap-1 text-xs text-slate-500">
                <Clock className="h-3 w-3" /> {s.startTime ?? '-'} → {s.endTime ?? ''}
              </p>
              <p className="font-bold text-slate-800">{s.title}</p>
              {s.description && <p className="mt-1 text-sm text-slate-600">{s.description}</p>}
              {s.cost ? (
                <p className="mt-1 text-sm text-sakura-600">
                  {s.cost} {s.costCurrency ?? ''}
                </p>
              ) : null}
            </li>
          ))}
        {schedules.data?.schedules?.length === 0 && (
          <li className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">일정이 없습니다.</li>
        )}
      </ol>
    </section>
  )
}

import { useParams, Link } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import type { ReactNode } from 'react'
import {
  ArrowLeft,
  CalendarDays,
  MapPin,
  Wallet,
  StickyNote,
  ListChecks,
  Ticket,
  Tag as TagIcon,
  Users,
  Image as ImageIcon,
  Plus,
} from 'lucide-react'
import { api } from '../lib/api'
import { Modal } from '../components/Modal'
import type { Day, ListDaysResponse, Trip } from '../lib/types'

export function TripDetailPage() {
  const { tripId = '' } = useParams()
  const queryClient = useQueryClient()
  const [showDateModal, setShowDateModal] = useState(false)
  const [newStartDate, setNewStartDate] = useState('')
  const [newEndDate, setNewEndDate] = useState('')

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

  const updateDates = useMutation({
    mutationFn: (data: { startDate: string; endDate: string }) =>
      api.patch<{ trip: Trip }>(`/trips/${tripId}`, {
        startDate: data.startDate,
        endDate: data.endDate,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['trip', tripId] })
      queryClient.invalidateQueries({ queryKey: ['days', tripId] })
      setShowDateModal(false)
    },
  })

  const handleOpenDateModal = () => {
    setNewStartDate(trip.data?.trip.startDate ?? '')
    setNewEndDate(trip.data?.trip.endDate ?? '')
    setShowDateModal(true)
  }

  const handleSaveDates = () => {
    if (newStartDate && newEndDate) {
      updateDates.mutate({ startDate: newStartDate, endDate: newEndDate })
    }
  }
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
      <div className="grid grid-cols-4 gap-2 sm:grid-cols-7">
        <NavBtn
          to={`/trips/${tripId}/expenses`}
          icon={<Wallet className="h-4 w-4" />}
          label="지출"
        />
        <NavBtn
          to={`/trips/${tripId}/notes`}
          icon={<StickyNote className="h-4 w-4" />}
          label="메모"
        />
        <NavBtn
          to={`/trips/${tripId}/checklist`}
          icon={<ListChecks className="h-4 w-4" />}
          label="체크"
        />
        <NavBtn
          to={`/trips/${tripId}/reservations`}
          icon={<Ticket className="h-4 w-4" />}
          label="예약"
        />
        <NavBtn to={`/trips/${tripId}/tags`} icon={<TagIcon className="h-4 w-4" />} label="태그" />
        <NavBtn
          to={`/trips/${tripId}/companions`}
          icon={<Users className="h-4 w-4" />}
          label="동행"
        />
        <NavBtn
          to={`/trips/${tripId}/media`}
          icon={<ImageIcon className="h-4 w-4" />}
          label="사진"
        />
      </div>
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-bold text-slate-800">일정</h2>
        <button
          onClick={handleOpenDateModal}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-1.5 text-sm text-white hover:bg-sakura-600"
        >
          <Plus className="h-4 w-4" /> 날짜 변경
        </button>
      </div>
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

      <Modal open={showDateModal} onClose={() => setShowDateModal(false)} title="여행 날짜 변경">
        <p className="text-sm text-slate-600 mb-4">
          날짜를 변경하면 Day가 새로 생성됩니다. 기존 Day에 연결된 일정은 유지됩니다.
        </p>
        <div className="grid grid-cols-2 gap-3">
          <label className="block">
            <span className="text-sm text-slate-700">시작일</span>
            <input
              type="date"
              value={newStartDate}
              onChange={(e) => setNewStartDate(e.target.value)}
              className="mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 text-sm"
            />
          </label>
          <label className="block">
            <span className="text-sm text-slate-700">종료일</span>
            <input
              type="date"
              value={newEndDate}
              onChange={(e) => setNewEndDate(e.target.value)}
              className="mt-1 block w-full rounded-md border border-slate-300 px-3 py-2 text-sm"
            />
          </label>
        </div>
        <div className="mt-4 flex justify-end gap-2">
          <button
            onClick={() => setShowDateModal(false)}
            className="rounded-md border border-slate-200 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50"
          >
            취소
          </button>
          <button
            onClick={handleSaveDates}
            disabled={updateDates.isPending}
            className="rounded-md bg-sakura-500 px-4 py-2 text-sm text-white hover:bg-sakura-600 disabled:opacity-50"
          >
            {updateDates.isPending ? '저장 중…' : '저장'}
          </button>
        </div>
      </Modal>
    </section>
  )
}

function NavBtn({ to, icon, label }: { to: string; icon: ReactNode; label: string }) {
  return (
    <Link
      to={to}
      className="flex flex-col items-center justify-center gap-1 rounded-xl bg-white p-3 text-xs font-bold text-sakura-700 shadow-sm hover:shadow"
    >
      {icon}
      {label}
    </Link>
  )
}

import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Plus, Trash2, Ticket } from 'lucide-react'
import { api } from '../lib/api'
import type {
  CreateReservationInput,
  ListReservationsResponse,
  Money,
  Reservation,
  ReservationType,
} from '../lib/types'

const TYPE_LABEL: Record<ReservationType, string> = {
  RESERVATION_TYPE_UNSPECIFIED: '기타',
  RESERVATION_TYPE_FLIGHT: '항공',
  RESERVATION_TYPE_HOTEL: '호텔',
  RESERVATION_TYPE_ACTIVITY: '액티비티',
  RESERVATION_TYPE_RESTAURANT: '식당',
  RESERVATION_TYPE_TRANSPORT: '교통',
}

const TYPES: ReservationType[] = [
  'RESERVATION_TYPE_FLIGHT',
  'RESERVATION_TYPE_HOTEL',
  'RESERVATION_TYPE_ACTIVITY',
  'RESERVATION_TYPE_RESTAURANT',
  'RESERVATION_TYPE_TRANSPORT',
]

export function ReservationsPage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [open, setOpen] = useState(false)

  const list = useQuery({
    queryKey: ['reservations', tripId],
    queryFn: () =>
      api.get<ListReservationsResponse>(`/trips/${tripId}/reservations`),
    enabled: !!tripId,
  })

  const createMut = useMutation({
    mutationFn: (input: CreateReservationInput) =>
      api.post(`/trips/${tripId}/reservations`, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['reservations', tripId] })
      setOpen(false)
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/reservations/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['reservations', tripId] }),
  })

  const items = list.data?.reservations ?? []

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <div className="flex items-center justify-between">
        <h1 className="flex items-center gap-2 text-xl font-bold text-slate-800">
          <Ticket className="h-5 w-5 text-sakura-500" /> 예약
        </h1>
        <button
          onClick={() => setOpen(true)}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-1.5 text-white hover:bg-sakura-600"
        >
          <Plus className="h-4 w-4" /> 추가
        </button>
      </div>

      {list.isLoading && <p>불러오는 중…</p>}
      {list.error && <p className="text-red-500">{(list.error as Error).message}</p>}

      <ul className="space-y-2">
        {items.map((r) => (
          <ReservationCard
            key={r.id}
            r={r}
            onDelete={() => deleteMut.mutate(r.id)}
          />
        ))}
        {items.length === 0 && !list.isLoading && (
          <li className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
            아직 예약이 없습니다.
          </li>
        )}
      </ul>

      {open && (
        <CreateReservationModal
          onClose={() => setOpen(false)}
          onSubmit={(v) => createMut.mutate(v)}
          pending={createMut.isPending}
        />
      )}
    </section>
  )
}

function ReservationCard({ r, onDelete }: { r: Reservation; onDelete: () => void }) {
  return (
    <li className="flex items-start justify-between rounded-xl bg-white p-4 shadow-sm">
      <div className="flex-1">
        <p className="text-xs text-slate-500">{TYPE_LABEL[r.type]}</p>
        <p className="font-bold text-slate-800">{r.vendor ?? '(공급자 없음)'}</p>
        {r.confirmNumber && (
          <p className="text-xs text-slate-500">예약번호 {r.confirmNumber}</p>
        )}
        {r.reservedAt && (
          <p className="text-xs text-slate-500">
            {new Date(r.reservedAt).toLocaleString()}
          </p>
        )}
        {r.cost?.amount ? (
          <p className="mt-1 text-sm text-sakura-600">
            {r.cost.amount} {r.cost.currency}
          </p>
        ) : null}
        {r.notes && <p className="mt-1 text-sm text-slate-600">{r.notes}</p>}
      </div>
      <button
        onClick={onDelete}
        className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
        aria-label="삭제"
      >
        <Trash2 className="h-4 w-4" />
      </button>
    </li>
  )
}

function CreateReservationModal({
  onClose,
  onSubmit,
  pending,
}: {
  onClose: () => void
  onSubmit: (v: CreateReservationInput) => void
  pending: boolean
}) {
  const [form, setForm] = useState({
    type: 'RESERVATION_TYPE_FLIGHT' as ReservationType,
    vendor: '',
    confirmNumber: '',
    reservedAt: '',
    notes: '',
    costAmount: '',
    costCurrency: 'JPY',
  })
  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/30 p-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-slate-800">예약 추가</h2>
        <form
          className="mt-4 space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            const cost: Money | undefined = form.costAmount
              ? { currency: form.costCurrency, amount: Number(form.costAmount) }
              : undefined
            onSubmit({
              type: form.type,
              vendor: form.vendor || undefined,
              confirmNumber: form.confirmNumber || undefined,
              reservedAt: form.reservedAt
                ? new Date(form.reservedAt).toISOString()
                : undefined,
              notes: form.notes || undefined,
              cost,
            })
          }}
        >
          <label className="block">
            <span className="text-sm text-slate-700">유형</span>
            <select
              value={form.type}
              onChange={(e) =>
                setForm({ ...form, type: e.target.value as ReservationType })
              }
              className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            >
              {TYPES.map((t) => (
                <option key={t} value={t}>
                  {TYPE_LABEL[t]}
                </option>
              ))}
            </select>
          </label>
          <Field
            label="공급자"
            value={form.vendor}
            onChange={(v) => setForm({ ...form, vendor: v })}
          />
          <Field
            label="예약번호"
            value={form.confirmNumber}
            onChange={(v) => setForm({ ...form, confirmNumber: v })}
          />
          <Field
            label="예약 시각"
            type="datetime-local"
            value={form.reservedAt}
            onChange={(v) => setForm({ ...form, reservedAt: v })}
          />
          <div className="grid grid-cols-2 gap-2">
            <Field
              label="비용"
              type="number"
              value={form.costAmount}
              onChange={(v) => setForm({ ...form, costAmount: v })}
            />
            <Field
              label="통화"
              value={form.costCurrency}
              onChange={(v) => setForm({ ...form, costCurrency: v })}
            />
          </div>
          <Field
            label="메모"
            value={form.notes}
            onChange={(v) => setForm({ ...form, notes: v })}
          />
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded-md border px-3 py-1.5">
              취소
            </button>
            <button
              type="submit"
              disabled={pending}
              className="rounded-md bg-sakura-500 px-3 py-1.5 text-white disabled:opacity-50"
            >
              {pending ? '추가 중…' : '추가'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

function Field({
  label,
  value,
  onChange,
  type = 'text',
}: {
  label: string
  value: string
  onChange: (v: string) => void
  type?: string
}) {
  return (
    <label className="block">
      <span className="text-sm text-slate-700">{label}</span>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 focus:border-sakura-400 focus:outline-none"
      />
    </label>
  )
}

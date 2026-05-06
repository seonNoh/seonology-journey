import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { Plus, MapPin, CalendarDays } from 'lucide-react'
import { api } from '../lib/api'
import type { CreateTripInput, ListTripsResponse, Trip } from '../lib/types'
import { useAuth } from '../hooks/useAuth'
import { CuteLoaderBlock } from '../components/CuteLoader'
import { TripListSkeleton } from '../components/LoadingPlaceholders'

export function TripListPage() {
  const auth = useAuth()
  const qc = useQueryClient()
  const [open, setOpen] = useState(false)

  const { data, isLoading, error } = useQuery({
    queryKey: ['trips'],
    queryFn: () => api.get<ListTripsResponse>('/trips'),
    enabled: auth.authenticated,
  })

  const createMut = useMutation({
    mutationFn: (input: CreateTripInput) => api.post<{ trip: Trip }>('/trips', input),
    // Optimistic UI: stub 하나를 캐시에 먼저 꽂아두고 서버 응답으로 교체한다.
    // 서버가 id 를 만들어 주므로 임시 id 를 쓴 뒤 onSuccess 에서 invalidate.
    onMutate: async (input) => {
      await qc.cancelQueries({ queryKey: ['trips'] })
      const prev = qc.getQueryData<ListTripsResponse>(['trips'])
      const optimistic: Trip = {
        id: `tmp-${Date.now()}`,
        ownerId: 'me',
        title: input.title,
        description: input.description,
        destination: input.destination,
        countryCode: input.countryCode,
        startDate: input.startDate,
        endDate: input.endDate,
        totalBudget: input.totalBudget,
        status: 'TRIP_STATUS_PLANNING',
      }
      qc.setQueryData<ListTripsResponse>(['trips'], (old) => ({
        trips: [optimistic, ...(old?.trips ?? [])],
      }))
      return { prev }
    },
    onError: (_err, _input, ctx) => {
      if (ctx?.prev) qc.setQueryData(['trips'], ctx.prev)
    },
    onSettled: () => {
      qc.invalidateQueries({ queryKey: ['trips'] })
      setOpen(false)
    },
  })

  if (!auth.ready) return <CuteLoaderBlock message="로그인 정보를 확인하고 있어요" />
  if (!auth.authenticated)
    return (
      <div className="rounded-2xl bg-white p-8 text-center shadow-sm">
        <p className="text-slate-700">로그인이 필요합니다.</p>
        <button
          onClick={auth.login}
          className="mt-4 rounded-md bg-sakura-500 px-4 py-2 text-white hover:bg-sakura-600"
        >
          로그인
        </button>
      </div>
    )

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-slate-800">내 여행</h1>
        <button
          onClick={() => setOpen(true)}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-1.5 text-white hover:bg-sakura-600"
        >
          <Plus className="h-4 w-4" /> 새 여행
        </button>
      </div>

      {isLoading && !data && <TripListSkeleton />}
      {error && <p className="text-red-500">{(error as Error).message}</p>}

      <ul className="grid grid-cols-1 gap-3 md:grid-cols-2">
        {data?.trips?.map((t) => (
          <li key={t.id}>
            <Link
              to={`/trips/${t.id}`}
              className="block rounded-xl bg-white p-4 shadow-sm hover:shadow transition"
            >
              <h2 className="font-bold text-slate-800">{t.title}</h2>
              {t.destination && (
                <p className="mt-1 flex items-center gap-1 text-sm text-slate-600">
                  <MapPin className="h-4 w-4" /> {t.destination}
                </p>
              )}
              {t.startDate && (
                <p className="mt-1 flex items-center gap-1 text-sm text-slate-500">
                  <CalendarDays className="h-4 w-4" /> {t.startDate} → {t.endDate}
                </p>
              )}
            </Link>
          </li>
        ))}
        {data?.trips && data.trips.length === 0 && (
          <li className="col-span-full rounded-xl bg-white p-8 text-center text-slate-500 shadow-sm">
            아직 여행이 없습니다. 새 여행을 만들어 보세요.
          </li>
        )}
      </ul>

      {open && (
        <CreateTripModal
          onClose={() => setOpen(false)}
          onSubmit={(v) => createMut.mutate(v)}
          pending={createMut.isPending}
        />
      )}
    </section>
  )
}

function CreateTripModal({
  onClose,
  onSubmit,
  pending,
}: {
  onClose: () => void
  onSubmit: (v: CreateTripInput) => void
  pending: boolean
}) {
  const [form, setForm] = useState<CreateTripInput & { currency: string }>({
    title: '',
    description: '',
    destination: '',
    startDate: '',
    endDate: '',
    countryCode: '',
    currency: 'JPY',
  })
  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/30 p-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-slate-800">새 여행</h2>
        <form
          className="mt-4 space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            const { currency, ...rest } = form
            onSubmit({
              ...rest,
              totalBudget: currency ? { currency, amount: 0 } : undefined,
            })
          }}
        >
          <Field
            label="제목"
            required
            value={form.title}
            onChange={(v) => setForm({ ...form, title: v })}
          />
          <Field
            label="목적지"
            value={form.destination ?? ''}
            onChange={(v) => setForm({ ...form, destination: v })}
          />
          <Field
            label="국가 코드 (JP/KR/...)"
            value={form.countryCode ?? ''}
            onChange={(v) => setForm({ ...form, countryCode: v.toUpperCase() })}
          />
          <div className="grid grid-cols-2 gap-2">
            <Field
              label="시작일"
              type="date"
              value={form.startDate ?? ''}
              onChange={(v) => setForm({ ...form, startDate: v })}
            />
            <Field
              label="종료일"
              type="date"
              value={form.endDate ?? ''}
              onChange={(v) => setForm({ ...form, endDate: v })}
            />
          </div>
          <Field
            label="설명"
            value={form.description ?? ''}
            onChange={(v) => setForm({ ...form, description: v })}
          />
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded-md border px-3 py-1.5">
              취소
            </button>
            <button
              type="submit"
              disabled={pending || !form.title}
              className="rounded-md bg-sakura-500 px-3 py-1.5 text-white disabled:opacity-50"
            >
              {pending ? '생성 중…' : '만들기'}
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
  required,
}: {
  label: string
  value: string
  onChange: (v: string) => void
  type?: string
  required?: boolean
}) {
  return (
    <label className="block">
      <span className="text-sm text-slate-700">
        {label}
        {required && <span className="text-sakura-500"> *</span>}
      </span>
      <input
        type={type}
        value={value}
        required={required}
        onChange={(e) => onChange(e.target.value)}
        className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2 focus:border-sakura-400 focus:outline-none"
      />
    </label>
  )
}

import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  ArrowLeft,
  Clock,
  Plus,
  Trash2,
  GripVertical,
  Utensils,
  BedDouble,
  Pencil,
} from 'lucide-react'
import {
  DndContext,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import {
  SortableContext,
  arrayMove,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { api } from '../lib/api'
import { MapPicker } from '../components/MapPicker'
import { ListRowsSkeleton } from '../components/LoadingPlaceholders'
import type {
  Accommodation,
  CreateScheduleInput,
  GeoPoint,
  ListMealsResponse,
  ListSchedulesResponse,
  Meal,
  MealSource,
  MealType,
  Money,
  Schedule,
} from '../lib/types'

const MEAL_TYPES: { value: Exclude<MealType, 'MEAL_TYPE_UNSPECIFIED'>; label: string }[] = [
  { value: 'MEAL_TYPE_BREAKFAST', label: '아침' },
  { value: 'MEAL_TYPE_LUNCH', label: '점심' },
  { value: 'MEAL_TYPE_DINNER', label: '저녁' },
]

const MEAL_SOURCES: { value: MealSource; label: string }[] = [
  { value: 'MEAL_SOURCE_HOTEL', label: '호텔식' },
  { value: 'MEAL_SOURCE_LOCAL', label: '현지식' },
  { value: 'MEAL_SOURCE_CONVENIENCE', label: '편의점/도시락' },
  { value: 'MEAL_SOURCE_SKIP', label: '거름' },
]

export function DayDetailPage() {
  const { dayId = '' } = useParams()
  const qc = useQueryClient()
  const [open, setOpen] = useState(false)

  const list = useQuery({
    queryKey: ['schedules', dayId],
    queryFn: () => api.get<ListSchedulesResponse>(`/days/${dayId}/schedules`),
    enabled: !!dayId,
  })

  const sorted = (list.data?.schedules ?? [])
    .slice()
    .sort((a, b) => (a.order ?? 0) - (b.order ?? 0))

  const createMut = useMutation({
    mutationFn: (input: CreateScheduleInput) => api.post(`/days/${dayId}/schedules`, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['schedules', dayId] })
      setOpen(false)
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/schedules/${id}`),
    // 리스트에서 즉시 제거하여 UX 를 개선. 실패 시 롤백.
    onMutate: async (id) => {
      await qc.cancelQueries({ queryKey: ['schedules', dayId] })
      const prev = qc.getQueryData<ListSchedulesResponse>(['schedules', dayId])
      qc.setQueryData<ListSchedulesResponse>(['schedules', dayId], (old) => ({
        schedules: (old?.schedules ?? []).filter((s) => s.id !== id),
      }))
      return { prev }
    },
    onError: (_err, _id, ctx) => {
      if (ctx?.prev) qc.setQueryData(['schedules', dayId], ctx.prev)
    },
    onSettled: () => qc.invalidateQueries({ queryKey: ['schedules', dayId] }),
  })

  const reorderMut = useMutation({
    mutationFn: (ids: string[]) =>
      api.post(`/days/${dayId}/schedules:reorder`, { scheduleIdsInOrder: ids }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['schedules', dayId] }),
  })

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 4 } }))

  const onDragEnd = (e: DragEndEvent) => {
    const { active, over } = e
    if (!over || active.id === over.id) return
    const ids = sorted.map((s) => s.id)
    const from = ids.indexOf(String(active.id))
    const to = ids.indexOf(String(over.id))
    if (from < 0 || to < 0) return
    const next = arrayMove(ids, from, to)
    qc.setQueryData<ListSchedulesResponse>(['schedules', dayId], (old) => {
      if (!old?.schedules) return old
      const map = new Map(old.schedules.map((s) => [s.id, s]))
      return {
        schedules: next.map((id, i) => ({ ...(map.get(id) as Schedule), order: i })),
      }
    })
    reorderMut.mutate(next)
  }

  return (
    <section className="space-y-4">
      <Link to="/trips" className="inline-flex items-center gap-1 text-sm text-sakura-700">
        <ArrowLeft className="h-4 w-4" /> 목록으로
      </Link>
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-slate-800">일정 타임라인</h1>
        <button
          onClick={() => setOpen(true)}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-1.5 text-white hover:bg-sakura-600"
        >
          <Plus className="h-4 w-4" /> 일정 추가
        </button>
      </div>

      {list.isLoading && !list.data && <ListRowsSkeleton message="일정을 펼치고 있어요" />}
      {list.error && <p className="text-red-500">{(list.error as Error).message}</p>}

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={onDragEnd}>
        <SortableContext items={sorted.map((s) => s.id)} strategy={verticalListSortingStrategy}>
          <ul className="space-y-2">
            {sorted.map((s) => (
              <SortableItem key={s.id} schedule={s} onDelete={() => deleteMut.mutate(s.id)} />
            ))}
            {sorted.length === 0 && !list.isLoading && (
              <li className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
                아직 일정이 없습니다.
              </li>
            )}
          </ul>
        </SortableContext>
      </DndContext>

      {open && (
        <CreateScheduleModal
          onClose={() => setOpen(false)}
          onSubmit={(v) => createMut.mutate(v)}
          pending={createMut.isPending}
        />
      )}

      <MealSection dayId={dayId} />
      <AccommodationSection dayId={dayId} />
    </section>
  )
}

function SortableItem({ schedule, onDelete }: { schedule: Schedule; onDelete: () => void }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: schedule.id,
  })
  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.6 : 1,
  }
  return (
    <li
      ref={setNodeRef}
      style={style}
      className="flex items-start gap-3 rounded-xl bg-white p-4 shadow-sm"
    >
      <button
        {...attributes}
        {...listeners}
        className="mt-1 cursor-grab text-slate-400 hover:text-slate-600"
        aria-label="순서 이동"
      >
        <GripVertical className="h-4 w-4" />
      </button>
      <div className="flex-1">
        <p className="flex items-center gap-1 text-xs text-slate-500">
          <Clock className="h-3 w-3" /> {schedule.startTime ?? '-'} → {schedule.endTime ?? ''}
        </p>
        <p className="font-bold text-slate-800">{schedule.title}</p>
        {schedule.placeName && (
          <p className="mt-0.5 text-xs text-slate-500">{schedule.placeName}</p>
        )}
        {schedule.notes && <p className="mt-1 text-sm text-slate-600">{schedule.notes}</p>}
        {schedule.cost ? (
          <p className="mt-1 text-sm text-sakura-600">
            {schedule.cost.amount} {schedule.cost.currency}
          </p>
        ) : null}
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

function CreateScheduleModal({
  onClose,
  onSubmit,
  pending,
}: {
  onClose: () => void
  onSubmit: (v: CreateScheduleInput) => void
  pending: boolean
}) {
  const [form, setForm] = useState({
    title: '',
    startTime: '',
    endTime: '',
    notes: '',
    placeName: '',
    costAmount: '',
    costCurrency: 'JPY',
  })
  const [location, setLocation] = useState<GeoPoint | undefined>(undefined)
  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/30 p-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl max-h-[90vh] overflow-y-auto">
        <h2 className="text-lg font-bold text-slate-800">일정 추가</h2>
        <form
          className="mt-4 space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            const amount = form.costAmount ? Number(form.costAmount) : 0
            onSubmit({
              title: form.title,
              startTime: form.startTime || undefined,
              endTime: form.endTime || undefined,
              notes: form.notes || undefined,
              placeName: form.placeName || undefined,
              location,
              cost: amount ? { currency: form.costCurrency, amount } : undefined,
            })
          }}
        >
          <Field
            label="제목"
            required
            value={form.title}
            onChange={(v) => setForm({ ...form, title: v })}
          />
          <div className="grid grid-cols-2 gap-2">
            <Field
              label="시작 시간"
              type="time"
              value={form.startTime}
              onChange={(v) => setForm({ ...form, startTime: v })}
            />
            <Field
              label="종료 시간"
              type="time"
              value={form.endTime}
              onChange={(v) => setForm({ ...form, endTime: v })}
            />
          </div>
          <Field
            label="장소"
            value={form.placeName}
            onChange={(v) => setForm({ ...form, placeName: v })}
          />
          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">위치 선택</label>
            <MapPicker
              latitude={location?.latitude}
              longitude={location?.longitude}
              onChange={(coords) => setLocation({ latitude: coords.lat, longitude: coords.lng })}
              accessToken={import.meta.env.VITE_MAPBOX_TOKEN}
            />
            {location && (
              <p className="mt-1 text-xs text-slate-500">
                {location.latitude.toFixed(5)}, {location.longitude.toFixed(5)}
              </p>
            )}
          </div>
          <Field label="메모" value={form.notes} onChange={(v) => setForm({ ...form, notes: v })} />
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
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded-md border px-3 py-1.5">
              취소
            </button>
            <button
              type="submit"
              disabled={pending || !form.title}
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

// =====================
// Meal section
// =====================

function MealSection({ dayId }: { dayId: string }) {
  const qc = useQueryClient()
  const [editing, setEditing] = useState<MealType | null>(null)

  const list = useQuery({
    queryKey: ['meals', dayId],
    queryFn: () => api.get<ListMealsResponse>(`/days/${dayId}/meals`),
    enabled: !!dayId,
  })

  const upsertMut = useMutation({
    mutationFn: (m: Partial<Meal> & { mealType: MealType }) => api.put(`/days/${dayId}/meals`, m),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['meals', dayId] })
      setEditing(null)
    },
  })

  const deleteMut = useMutation({
    mutationFn: (mealType: MealType) => api.del(`/days/${dayId}/meals/${mealType}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['meals', dayId] }),
  })

  const byType = new Map((list.data?.meals ?? []).map((m) => [m.mealType, m]))

  return (
    <div className="rounded-2xl bg-white p-4 shadow-sm">
      <h2 className="mb-3 flex items-center gap-2 text-base font-bold text-slate-800">
        <Utensils className="h-4 w-4 text-sakura-500" /> 식사
      </h2>
      <ul className="space-y-2">
        {MEAL_TYPES.map(({ value, label }) => {
          const meal = byType.get(value)
          return (
            <li
              key={value}
              className="flex items-center justify-between rounded-lg border border-slate-100 px-3 py-2"
            >
              <div className="flex-1">
                <p className="text-xs text-slate-500">{label}</p>
                {meal ? (
                  <p className="text-sm text-slate-700">
                    {meal.restaurantName || meal.menu || '(메모 없음)'}
                    {meal.cost?.amount ? (
                      <span className="ml-2 text-sakura-600">
                        {meal.cost.amount} {meal.cost.currency}
                      </span>
                    ) : null}
                  </p>
                ) : (
                  <p className="text-sm text-slate-400">미설정</p>
                )}
              </div>
              <div className="flex items-center gap-1">
                <button
                  onClick={() => setEditing(value)}
                  className="rounded-md p-1 text-slate-400 hover:bg-sakura-50 hover:text-sakura-600"
                  aria-label="편집"
                >
                  <Pencil className="h-4 w-4" />
                </button>
                {meal && (
                  <button
                    onClick={() => deleteMut.mutate(value)}
                    className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
                    aria-label="삭제"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                )}
              </div>
            </li>
          )
        })}
      </ul>
      {editing && (
        <MealEditModal
          dayId={dayId}
          mealType={editing}
          initial={byType.get(editing)}
          onClose={() => setEditing(null)}
          onSubmit={(m) => upsertMut.mutate(m)}
          pending={upsertMut.isPending}
        />
      )}
    </div>
  )
}

function MealEditModal({
  mealType,
  initial,
  onClose,
  onSubmit,
  pending,
}: {
  dayId: string
  mealType: MealType
  initial?: Meal
  onClose: () => void
  onSubmit: (m: Partial<Meal> & { mealType: MealType }) => void
  pending: boolean
}) {
  const [form, setForm] = useState({
    source: (initial?.source ?? 'MEAL_SOURCE_LOCAL') as MealSource,
    restaurantName: initial?.restaurantName ?? '',
    menu: initial?.menu ?? '',
    rating: String(initial?.rating ?? ''),
    review: initial?.review ?? '',
    costAmount: initial?.cost?.amount ? String(initial.cost.amount) : '',
    costCurrency: initial?.cost?.currency ?? 'JPY',
  })
  const label = MEAL_TYPES.find((t) => t.value === mealType)?.label ?? ''
  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/30 p-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-slate-800">{label} 입력</h2>
        <form
          className="mt-4 space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            const cost: Money | undefined = form.costAmount
              ? { currency: form.costCurrency, amount: Number(form.costAmount) }
              : undefined
            onSubmit({
              mealType,
              source: form.source,
              restaurantName: form.restaurantName || undefined,
              menu: form.menu || undefined,
              rating: form.rating ? Number(form.rating) : undefined,
              review: form.review || undefined,
              cost,
            })
          }}
        >
          <label className="block">
            <span className="text-sm text-slate-700">유형</span>
            <select
              value={form.source}
              onChange={(e) => setForm({ ...form, source: e.target.value as MealSource })}
              className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            >
              {MEAL_SOURCES.map((s) => (
                <option key={s.value} value={s.value}>
                  {s.label}
                </option>
              ))}
            </select>
          </label>
          <Field
            label="식당"
            value={form.restaurantName}
            onChange={(v) => setForm({ ...form, restaurantName: v })}
          />
          <Field label="메뉴" value={form.menu} onChange={(v) => setForm({ ...form, menu: v })} />
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
            label="평점 (1-5)"
            type="number"
            value={form.rating}
            onChange={(v) => setForm({ ...form, rating: v })}
          />
          <Field
            label="후기"
            value={form.review}
            onChange={(v) => setForm({ ...form, review: v })}
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
              {pending ? '저장 중…' : '저장'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// =====================
// Accommodation section
// =====================

function AccommodationSection({ dayId }: { dayId: string }) {
  const qc = useQueryClient()
  const [editing, setEditing] = useState(false)

  const get = useQuery({
    queryKey: ['accommodation', dayId],
    queryFn: async () => {
      try {
        return await api.get<{ accommodation?: Accommodation }>(`/days/${dayId}/accommodation`)
      } catch (e) {
        // 미설정이면 404 가능 — 빈 응답으로 처리
        if ((e as Error).message.includes('404')) return { accommodation: undefined }
        throw e
      }
    },
    enabled: !!dayId,
  })

  const upsertMut = useMutation({
    mutationFn: (a: Partial<Accommodation>) => api.put(`/days/${dayId}/accommodation`, a),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['accommodation', dayId] })
      setEditing(false)
    },
  })

  const deleteMut = useMutation({
    mutationFn: () => api.del(`/days/${dayId}/accommodation`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['accommodation', dayId] }),
  })

  const a = get.data?.accommodation

  return (
    <div className="rounded-2xl bg-white p-4 shadow-sm">
      <div className="mb-3 flex items-center justify-between">
        <h2 className="flex items-center gap-2 text-base font-bold text-slate-800">
          <BedDouble className="h-4 w-4 text-sakura-500" /> 숙박
        </h2>
        <div className="flex items-center gap-1">
          <button
            onClick={() => setEditing(true)}
            className="rounded-md p-1 text-slate-400 hover:bg-sakura-50 hover:text-sakura-600"
            aria-label="편집"
          >
            <Pencil className="h-4 w-4" />
          </button>
          {a && (
            <button
              onClick={() => deleteMut.mutate()}
              className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
              aria-label="삭제"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          )}
        </div>
      </div>
      {a ? (
        <div className="space-y-1 text-sm">
          <p className="font-bold text-slate-800">{a.name}</p>
          {a.address && <p className="text-slate-500">{a.address}</p>}
          {(a.checkInTime || a.checkOutTime) && (
            <p className="text-slate-500">
              체크인 {a.checkInTime ?? '-'} / 체크아웃 {a.checkOutTime ?? '-'}
            </p>
          )}
          {a.cost?.amount ? (
            <p className="text-sakura-600">
              {a.cost.amount} {a.cost.currency}
            </p>
          ) : null}
          {a.amenities && <p className="text-slate-500">시설: {a.amenities}</p>}
        </div>
      ) : (
        <p className="text-sm text-slate-400">아직 등록된 숙소가 없습니다.</p>
      )}
      {editing && (
        <AccommodationEditModal
          initial={a}
          onClose={() => setEditing(false)}
          onSubmit={(v) => upsertMut.mutate(v)}
          pending={upsertMut.isPending}
        />
      )}
    </div>
  )
}

function AccommodationEditModal({
  initial,
  onClose,
  onSubmit,
  pending,
}: {
  initial?: Accommodation
  onClose: () => void
  onSubmit: (a: Partial<Accommodation>) => void
  pending: boolean
}) {
  const [form, setForm] = useState({
    name: initial?.name ?? '',
    address: initial?.address ?? '',
    checkInTime: initial?.checkInTime ?? '',
    checkOutTime: initial?.checkOutTime ?? '',
    amenities: initial?.amenities ?? '',
    costAmount: initial?.cost?.amount ? String(initial.cost.amount) : '',
    costCurrency: initial?.cost?.currency ?? 'JPY',
  })
  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/30 p-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-slate-800">숙소 입력</h2>
        <form
          className="mt-4 space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            const cost: Money | undefined = form.costAmount
              ? { currency: form.costCurrency, amount: Number(form.costAmount) }
              : undefined
            onSubmit({
              name: form.name,
              address: form.address || undefined,
              checkInTime: form.checkInTime || undefined,
              checkOutTime: form.checkOutTime || undefined,
              amenities: form.amenities || undefined,
              cost,
            })
          }}
        >
          <Field
            label="숙소명"
            required
            value={form.name}
            onChange={(v) => setForm({ ...form, name: v })}
          />
          <Field
            label="주소"
            value={form.address}
            onChange={(v) => setForm({ ...form, address: v })}
          />
          <div className="grid grid-cols-2 gap-2">
            <Field
              label="체크인"
              type="time"
              value={form.checkInTime}
              onChange={(v) => setForm({ ...form, checkInTime: v })}
            />
            <Field
              label="체크아웃"
              type="time"
              value={form.checkOutTime}
              onChange={(v) => setForm({ ...form, checkOutTime: v })}
            />
          </div>
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
            label="시설/조식 등"
            value={form.amenities}
            onChange={(v) => setForm({ ...form, amenities: v })}
          />
          <div className="flex justify-end gap-2 pt-2">
            <button type="button" onClick={onClose} className="rounded-md border px-3 py-1.5">
              취소
            </button>
            <button
              type="submit"
              disabled={pending || !form.name}
              className="rounded-md bg-sakura-500 px-3 py-1.5 text-white disabled:opacity-50"
            >
              {pending ? '저장 중…' : '저장'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

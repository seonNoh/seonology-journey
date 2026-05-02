import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Clock, Plus, Trash2, GripVertical } from 'lucide-react'
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
import type { CreateScheduleInput, ListSchedulesResponse, Schedule } from '../lib/types'

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
    mutationFn: (input: CreateScheduleInput) =>
      api.post(`/days/${dayId}/schedules`, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['schedules', dayId] })
      setOpen(false)
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/schedules/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['schedules', dayId] }),
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

      {list.isLoading && <p>불러오는 중…</p>}
      {list.error && <p className="text-red-500">{(list.error as Error).message}</p>}

      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={onDragEnd}>
        <SortableContext items={sorted.map((s) => s.id)} strategy={verticalListSortingStrategy}>
          <ul className="space-y-2">
            {sorted.map((s) => (
              <SortableItem
                key={s.id}
                schedule={s}
                onDelete={() => deleteMut.mutate(s.id)}
              />
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
        {schedule.notes && (
          <p className="mt-1 text-sm text-slate-600">{schedule.notes}</p>
        )}
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
  return (
    <div className="fixed inset-0 z-20 flex items-center justify-center bg-black/30 p-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
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
              cost: amount
                ? { currency: form.costCurrency, amount }
                : undefined,
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
          <Field
            label="메모"
            value={form.notes}
            onChange={(v) => setForm({ ...form, notes: v })}
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

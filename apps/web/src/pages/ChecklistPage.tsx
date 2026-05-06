import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Plus, Trash2, Check } from 'lucide-react'
import { api } from '../lib/api'
import { ListRowsSkeleton } from '../components/LoadingPlaceholders'
import type { ChecklistCategory, ChecklistItem, ListChecklistResponse } from '../lib/types'

const CAT_LABEL: Record<ChecklistCategory, string> = {
  CHECKLIST_CATEGORY_UNSPECIFIED: '기타',
  CHECKLIST_CATEGORY_PACKING: '짐',
  CHECKLIST_CATEGORY_TODO: '할일',
  CHECKLIST_CATEGORY_BOOKING: '예약',
}

const CATS: ChecklistCategory[] = [
  'CHECKLIST_CATEGORY_PACKING',
  'CHECKLIST_CATEGORY_TODO',
  'CHECKLIST_CATEGORY_BOOKING',
]

export function ChecklistPage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [item, setItem] = useState('')
  const [category, setCategory] = useState<ChecklistCategory>('CHECKLIST_CATEGORY_PACKING')

  const list = useQuery({
    queryKey: ['checklist', tripId],
    queryFn: () => api.get<ListChecklistResponse>(`/trips/${tripId}/checklist`),
    enabled: !!tripId,
  })

  const createMut = useMutation({
    mutationFn: (input: { item: string; category: ChecklistCategory }) =>
      api.post<{ item: ChecklistItem }>(`/trips/${tripId}/checklist`, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['checklist', tripId] })
      setItem('')
    },
  })

  const toggleMut = useMutation({
    mutationFn: (it: ChecklistItem) =>
      api.patch(`/checklist/${it.id}`, { isChecked: !it.isChecked }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['checklist', tripId] }),
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/checklist/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['checklist', tripId] }),
  })

  const grouped = (list.data?.items ?? []).reduce<Record<string, ChecklistItem[]>>((acc, it) => {
    const key = it.category ?? 'CHECKLIST_CATEGORY_UNSPECIFIED'
    acc[key] = acc[key] ?? []
    acc[key].push(it)
    return acc
  }, {})

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행 상세
      </Link>
      <h1 className="text-xl font-bold text-slate-800">체크리스트</h1>

      <form
        className="rounded-2xl bg-white p-4 shadow-sm flex items-center gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          if (!item.trim()) return
          createMut.mutate({ item, category })
        }}
      >
        <select
          value={category}
          onChange={(e) => setCategory(e.target.value as ChecklistCategory)}
          className="rounded-md border border-slate-200 px-2 py-2 focus:border-sakura-400 focus:outline-none"
        >
          {CATS.map((c) => (
            <option key={c} value={c}>
              {CAT_LABEL[c]}
            </option>
          ))}
        </select>
        <input
          value={item}
          onChange={(e) => setItem(e.target.value)}
          placeholder="추가할 항목..."
          className="flex-1 rounded-md border border-slate-200 px-3 py-2 focus:border-sakura-400 focus:outline-none"
        />
        <button
          type="submit"
          disabled={createMut.isPending || !item.trim()}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-2 text-white disabled:opacity-50"
        >
          <Plus className="h-4 w-4" />
        </button>
      </form>

      {list.isLoading && !list.data && (
        <ListRowsSkeleton message="체크리스트를 가져오는 중이에요" />
      )}
      {Object.entries(grouped).map(([cat, items]) => (
        <div key={cat} className="space-y-2">
          <h2 className="text-sm font-bold text-slate-600">
            {CAT_LABEL[cat as ChecklistCategory] ?? cat}
          </h2>
          <ul className="space-y-1">
            {items.map((it) => (
              <li key={it.id} className="flex items-center gap-2 rounded-xl bg-white p-3 shadow-sm">
                <button
                  onClick={() => toggleMut.mutate(it)}
                  className={
                    'flex h-5 w-5 items-center justify-center rounded border ' +
                    (it.isChecked
                      ? 'border-sakura-500 bg-sakura-500 text-white'
                      : 'border-slate-300')
                  }
                  aria-label="토글"
                >
                  {it.isChecked && <Check className="h-3 w-3" />}
                </button>
                <span
                  className={
                    'flex-1 ' + (it.isChecked ? 'line-through text-slate-400' : 'text-slate-800')
                  }
                >
                  {it.item}
                </span>
                <button
                  onClick={() => deleteMut.mutate(it.id)}
                  className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
                  aria-label="삭제"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </li>
            ))}
          </ul>
        </div>
      ))}
      {list.data?.items?.length === 0 && (
        <p className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
          항목이 없습니다.
        </p>
      )}
    </section>
  )
}

import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Plus, Trash2, Tag as TagIcon, X } from 'lucide-react'
import { api } from '../lib/api'
import type { ListTagsResponse, Tag } from '../lib/types'

export function TagsPage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [name, setName] = useState('')
  const [color, setColor] = useState('#f9a8d4')

  const all = useQuery({
    queryKey: ['tags'],
    queryFn: () => api.get<ListTagsResponse>('/tags'),
  })

  const linked = useQuery({
    queryKey: ['trip-tags', tripId],
    queryFn: () => api.get<ListTagsResponse>(`/trips/${tripId}/tags`),
    enabled: !!tripId,
  })

  const createMut = useMutation({
    mutationFn: (input: { name: string; color: string }) =>
      api.post('/tags', input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['tags'] })
      setName('')
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/tags/${id}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['tags'] })
      qc.invalidateQueries({ queryKey: ['trip-tags', tripId] })
    },
  })

  const attachMut = useMutation({
    mutationFn: (tagId: string) => api.put(`/trips/${tripId}/tags/${tagId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['trip-tags', tripId] }),
  })

  const detachMut = useMutation({
    mutationFn: (tagId: string) => api.del(`/trips/${tripId}/tags/${tagId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['trip-tags', tripId] }),
  })

  const linkedIds = new Set((linked.data?.tags ?? []).map((t) => t.id))
  const allTags = all.data?.tags ?? []

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <h1 className="flex items-center gap-2 text-xl font-bold text-slate-800">
        <TagIcon className="h-5 w-5 text-sakura-500" /> 태그
      </h1>

      <form
        className="flex items-end gap-2 rounded-2xl bg-white p-4 shadow-sm"
        onSubmit={(e) => {
          e.preventDefault()
          if (!name.trim()) return
          createMut.mutate({ name: name.trim(), color })
        }}
      >
        <label className="flex-1 block">
          <span className="text-sm text-slate-700">이름</span>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
        <label className="block">
          <span className="text-sm text-slate-700">색상</span>
          <input
            type="color"
            value={color}
            onChange={(e) => setColor(e.target.value)}
            className="mt-1 h-10 w-14 rounded-md border border-slate-200"
          />
        </label>
        <button
          type="submit"
          disabled={createMut.isPending}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-2 text-white disabled:opacity-50"
        >
          <Plus className="h-4 w-4" /> 생성
        </button>
      </form>

      <div className="rounded-2xl bg-white p-4 shadow-sm">
        <h2 className="mb-2 text-sm font-bold text-slate-700">이 여행에 붙은 태그</h2>
        <div className="flex flex-wrap gap-2">
          {(linked.data?.tags ?? []).map((t) => (
            <TagPill
              key={t.id}
              tag={t}
              action="detach"
              onAction={() => detachMut.mutate(t.id)}
            />
          ))}
          {(linked.data?.tags ?? []).length === 0 && (
            <p className="text-sm text-slate-400">아직 붙은 태그가 없습니다.</p>
          )}
        </div>
      </div>

      <div className="rounded-2xl bg-white p-4 shadow-sm">
        <h2 className="mb-2 text-sm font-bold text-slate-700">전체 태그</h2>
        <div className="flex flex-wrap gap-2">
          {allTags.map((t) => (
            <div key={t.id} className="flex items-center gap-1">
              <button
                onClick={() =>
                  linkedIds.has(t.id)
                    ? detachMut.mutate(t.id)
                    : attachMut.mutate(t.id)
                }
                className="rounded-full px-3 py-1 text-xs font-bold"
                style={{
                  background: linkedIds.has(t.id) ? t.color ?? '#f472b6' : '#f1f5f9',
                  color: linkedIds.has(t.id) ? 'white' : '#475569',
                }}
              >
                {t.name}
              </button>
              <button
                onClick={() => deleteMut.mutate(t.id)}
                className="rounded-full p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
                aria-label="삭제"
              >
                <Trash2 className="h-3 w-3" />
              </button>
            </div>
          ))}
          {allTags.length === 0 && (
            <p className="text-sm text-slate-400">등록된 태그가 없습니다.</p>
          )}
        </div>
      </div>
    </section>
  )
}

function TagPill({
  tag,
  onAction,
}: {
  tag: Tag
  action: 'detach'
  onAction: () => void
}) {
  return (
    <span
      className="inline-flex items-center gap-1 rounded-full px-3 py-1 text-xs font-bold text-white"
      style={{ background: tag.color ?? '#f472b6' }}
    >
      {tag.name}
      <button onClick={onAction} aria-label="제거">
        <X className="h-3 w-3" />
      </button>
    </span>
  )
}

import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Plus, Trash2 } from 'lucide-react'
import { api } from '../lib/api'
import { ListRowsSkeleton } from '../components/LoadingPlaceholders'
import type { ListNotesResponse, Note } from '../lib/types'

export function NotesPage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [content, setContent] = useState('')
  const [mood, setMood] = useState('')

  const list = useQuery({
    queryKey: ['notes', tripId],
    queryFn: () => api.get<ListNotesResponse>(`/trips/${tripId}/notes`),
    enabled: !!tripId,
  })

  const createMut = useMutation({
    mutationFn: (input: { content: string; mood?: string }) =>
      api.post<{ note: Note }>(`/trips/${tripId}/notes`, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['notes', tripId] })
      setContent('')
      setMood('')
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/notes/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['notes', tripId] }),
  })

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행 상세
      </Link>
      <h1 className="text-xl font-bold text-slate-800">메모</h1>

      <form
        className="rounded-2xl bg-white p-4 shadow-sm space-y-2"
        onSubmit={(e) => {
          e.preventDefault()
          if (!content.trim()) return
          createMut.mutate({ content, mood: mood || undefined })
        }}
      >
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={3}
          placeholder="오늘의 한 줄..."
          className="w-full rounded-md border border-slate-200 px-3 py-2 focus:border-sakura-400 focus:outline-none"
        />
        <div className="flex items-center gap-2">
          <input
            value={mood}
            onChange={(e) => setMood(e.target.value)}
            placeholder="기분 (happy/excited/...)"
            className="flex-1 rounded-md border border-slate-200 px-3 py-2 focus:border-sakura-400 focus:outline-none"
          />
          <button
            type="submit"
            disabled={createMut.isPending || !content.trim()}
            className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-2 text-white disabled:opacity-50"
          >
            <Plus className="h-4 w-4" /> 추가
          </button>
        </div>
      </form>

      <ul className="space-y-2">
        {list.isLoading && !list.data && <ListRowsSkeleton message="메모를 펼치고 있어요" />}
        {list.data?.notes?.map((n) => (
          <li
            key={n.id}
            className="flex items-start justify-between rounded-xl bg-white p-4 shadow-sm"
          >
            <div className="flex-1">
              {n.mood && <p className="text-xs text-sakura-600">#{n.mood}</p>}
              <p className="mt-1 whitespace-pre-wrap text-slate-800">{n.content}</p>
            </div>
            <button
              onClick={() => deleteMut.mutate(n.id)}
              className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
              aria-label="삭제"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </li>
        ))}
        {list.data?.notes?.length === 0 && (
          <li className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
            아직 메모가 없습니다.
          </li>
        )}
      </ul>
    </section>
  )
}

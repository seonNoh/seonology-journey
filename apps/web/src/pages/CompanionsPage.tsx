import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Plus, Trash2, Users } from 'lucide-react'
import { api } from '../lib/api'
import { ListRowsSkeleton } from '../components/LoadingPlaceholders'
import type { CompanionRole, ListCompanionsResponse } from '../lib/types'

const ROLE_LABEL: Record<CompanionRole, string> = {
  COMPANION_ROLE_UNSPECIFIED: '미지정',
  COMPANION_ROLE_OWNER: '소유자',
  COMPANION_ROLE_EDITOR: '편집자',
  COMPANION_ROLE_VIEWER: '뷰어',
}

const ROLES: CompanionRole[] = ['COMPANION_ROLE_EDITOR', 'COMPANION_ROLE_VIEWER']

export function CompanionsPage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [memberId, setMemberId] = useState('')
  const [role, setRole] = useState<CompanionRole>('COMPANION_ROLE_EDITOR')

  const list = useQuery({
    queryKey: ['companions', tripId],
    queryFn: () => api.get<ListCompanionsResponse>(`/trips/${tripId}/companions`),
    enabled: !!tripId,
  })

  const addMut = useMutation({
    mutationFn: (input: { memberId: string; role: CompanionRole }) =>
      api.post(`/trips/${tripId}/companions`, input),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['companions', tripId] })
      setMemberId('')
    },
  })

  const updateMut = useMutation({
    mutationFn: ({ id, role }: { id: string; role: CompanionRole }) =>
      api.patch(`/trips/${tripId}/companions/${id}`, { role }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['companions', tripId] }),
  })

  const removeMut = useMutation({
    mutationFn: (id: string) => api.del(`/trips/${tripId}/companions/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['companions', tripId] }),
  })

  const items = list.data?.companions ?? []

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <h1 className="flex items-center gap-2 text-xl font-bold text-slate-800">
        <Users className="h-5 w-5 text-sakura-500" /> 동행자
      </h1>

      <form
        className="flex items-end gap-2 rounded-2xl bg-white p-4 shadow-sm"
        onSubmit={(e) => {
          e.preventDefault()
          if (!memberId.trim()) return
          addMut.mutate({ memberId: memberId.trim(), role })
        }}
      >
        <label className="flex-1 block">
          <span className="text-sm text-slate-700">멤버 ID (Keycloak sub)</span>
          <input
            value={memberId}
            onChange={(e) => setMemberId(e.target.value)}
            className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
        <label className="block">
          <span className="text-sm text-slate-700">역할</span>
          <select
            value={role}
            onChange={(e) => setRole(e.target.value as CompanionRole)}
            className="mt-1 rounded-md border border-slate-200 px-3 py-2"
          >
            {ROLES.map((r) => (
              <option key={r} value={r}>
                {ROLE_LABEL[r]}
              </option>
            ))}
          </select>
        </label>
        <button
          type="submit"
          disabled={addMut.isPending}
          className="flex items-center gap-1 rounded-md bg-sakura-500 px-3 py-2 text-white disabled:opacity-50"
        >
          <Plus className="h-4 w-4" /> 추가
        </button>
      </form>

      <ul className="space-y-2">
        {list.isLoading && !list.data && <ListRowsSkeleton message="동행자를 확인하고 있어요" />}
        {items.map((c) => (
          <li
            key={c.memberId}
            className="flex items-center justify-between rounded-xl bg-white p-4 shadow-sm"
          >
            <div className="flex items-center gap-3">
              {c.avatarUrl ? (
                <img src={c.avatarUrl} alt="" className="h-10 w-10 rounded-full object-cover" />
              ) : (
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-sakura-100 text-sm font-bold text-sakura-700">
                  {(c.displayName ?? c.memberId).slice(0, 1).toUpperCase()}
                </div>
              )}
              <div>
                <p className="font-bold text-slate-800">{c.displayName ?? c.memberId}</p>
                <p className="text-xs text-slate-500">{c.memberId}</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <select
                value={c.role}
                onChange={(e) =>
                  updateMut.mutate({
                    id: c.memberId,
                    role: e.target.value as CompanionRole,
                  })
                }
                disabled={c.role === 'COMPANION_ROLE_OWNER'}
                className="rounded-md border border-slate-200 px-2 py-1 text-sm"
              >
                {(c.role === 'COMPANION_ROLE_OWNER'
                  ? (['COMPANION_ROLE_OWNER'] as CompanionRole[])
                  : ROLES
                ).map((r) => (
                  <option key={r} value={r}>
                    {ROLE_LABEL[r]}
                  </option>
                ))}
              </select>
              {c.role !== 'COMPANION_ROLE_OWNER' && (
                <button
                  onClick={() => removeMut.mutate(c.memberId)}
                  className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
                  aria-label="제거"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              )}
            </div>
          </li>
        ))}
        {items.length === 0 && !list.isLoading && (
          <li className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
            동행자가 없습니다.
          </li>
        )}
      </ul>
    </section>
  )
}

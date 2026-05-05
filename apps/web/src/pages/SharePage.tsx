import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, LinkIcon, Trash2 } from 'lucide-react'
import { QRCodeSVG } from 'qrcode.react'
import { api } from '../lib/api'
import type { Share } from '../lib/types'

interface ListSharesResponse {
  shares: Share[]
}

type CompanionRoleOption = 'COMPANION_ROLE_VIEWER' | 'COMPANION_ROLE_EDITOR'

export function SharePage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [expHours, setExpHours] = useState(72)
  const [permission, setPermission] = useState<CompanionRoleOption>('COMPANION_ROLE_VIEWER')

  const list = useQuery({
    queryKey: ['shares', tripId],
    queryFn: () => api.get<ListSharesResponse>(`/trips/${tripId}/shares`),
    enabled: !!tripId,
  })

  const createMut = useMutation({
    mutationFn: () =>
      api.post<Share>(`/trips/${tripId}/shares`, {
        permission,
        expiresInHours: expHours,
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['shares', tripId] }),
  })

  const deleteMut = useMutation({
    mutationFn: (code: string) => api.del(`/shares/${code}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['shares', tripId] }),
  })

  const shares = list.data?.shares ?? []
  const baseUrl = window.location.origin

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <h1 className="flex items-center gap-2 text-xl font-bold text-slate-800">
        <LinkIcon className="h-5 w-5 text-sakura-500" /> 공유 링크
      </h1>

      <div className="rounded-2xl bg-white p-4 shadow-sm space-y-3">
        <div className="grid grid-cols-2 gap-3">
          <label className="block">
            <span className="text-sm text-slate-700">권한</span>
            <select
              value={permission}
              onChange={(e) => setPermission(e.target.value as CompanionRoleOption)}
              className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            >
              <option value="COMPANION_ROLE_VIEWER">조회 전용</option>
              <option value="COMPANION_ROLE_EDITOR">편집 가능</option>
            </select>
          </label>
          <label className="block">
            <span className="text-sm text-slate-700">만료 (시간)</span>
            <input
              type="number"
              min={1}
              max={720}
              value={expHours}
              onChange={(e) => setExpHours(Number(e.target.value))}
              className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
        </div>
        <button
          onClick={() => createMut.mutate()}
          disabled={createMut.isPending}
          className="w-full rounded-md bg-sakura-500 px-3 py-2 text-white hover:bg-sakura-600 disabled:opacity-50"
        >
          {createMut.isPending ? '생성 중…' : '링크 생성'}
        </button>
      </div>

      {shares.length > 0 && (
        <ul className="space-y-3">
          {shares.map((s) => {
            const shareUrl = `${baseUrl}/join/${s.code}`
            return (
              <li key={s.code} className="rounded-2xl bg-white p-4 shadow-sm">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex-1 space-y-2">
                    <p className="text-sm font-mono text-slate-700 break-all">{shareUrl}</p>
                    <p className="text-xs text-slate-500">
                      권한: {s.permission === 'COMPANION_ROLE_EDITOR' ? '편집' : '조회'}
                      {s.expiresAt && ` / 만료: ${new Date(s.expiresAt).toLocaleString()}`}
                    </p>
                    <button
                      onClick={() => navigator.clipboard.writeText(shareUrl)}
                      className="text-xs text-sakura-600 underline"
                    >
                      링크 복사
                    </button>
                  </div>
                  <QRCodeSVG value={shareUrl} size={80} />
                  <button
                    onClick={() => deleteMut.mutate(s.code)}
                    className="rounded-md p-1 text-slate-400 hover:bg-red-50 hover:text-red-500"
                    aria-label="삭제"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </li>
            )
          })}
        </ul>
      )}

      {shares.length === 0 && !list.isLoading && (
        <p className="rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
          아직 공유 링크가 없습니다.
        </p>
      )}
    </section>
  )
}

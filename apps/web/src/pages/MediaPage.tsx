import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ArrowLeft, Trash2, Image as ImageIcon, Upload } from 'lucide-react'
import { api } from '../lib/api'
import { Lightbox } from '../components/Lightbox'
import { GridSkeleton } from '../components/LoadingPlaceholders'
import type { GetMediaUrlResponse, GetUploadUrlResponse, ListMediaResponse } from '../lib/types'

export function MediaPage() {
  const { tripId = '' } = useParams()
  const qc = useQueryClient()
  const [uploading, setUploading] = useState(false)
  const [progress, setProgress] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [caption, setCaption] = useState('')
  const [lightboxIdx, setLightboxIdx] = useState<number | null>(null)

  const list = useQuery({
    queryKey: ['media', tripId],
    queryFn: () => api.get<ListMediaResponse>(`/trips/${tripId}/media`),
    enabled: !!tripId,
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => api.del(`/media/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['media', tripId] }),
  })

  async function onPickFile(file: File) {
    setError(null)
    setUploading(true)
    setProgress(0)
    try {
      const presign = await api.post<GetUploadUrlResponse>(`/trips/${tripId}/media:upload-url`, {
        filename: file.name,
        mimeType: file.type,
        size: file.size,
      })
      await new Promise<void>((resolve, reject) => {
        const xhr = new XMLHttpRequest()
        xhr.upload.addEventListener('progress', (e) => {
          if (e.lengthComputable) setProgress(Math.round((e.loaded / e.total) * 100))
        })
        xhr.addEventListener('load', () => {
          if (xhr.status >= 200 && xhr.status < 300) resolve()
          else reject(new Error(`S3 PUT 실패: ${xhr.status}`))
        })
        xhr.addEventListener('error', () => reject(new Error('Upload network error')))
        xhr.open('PUT', presign.uploadUrl)
        xhr.setRequestHeader('Content-Type', file.type)
        xhr.send(file)
      })
      await api.post(`/trips/${tripId}/media:confirm`, {
        mediaId: presign.mediaId,
        s3Key: presign.s3Key,
        caption: caption || undefined,
      })
      setCaption('')
      qc.invalidateQueries({ queryKey: ['media', tripId] })
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setUploading(false)
      setProgress(0)
    }
  }

  const items = list.data?.items ?? []

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <h1 className="flex items-center gap-2 text-xl font-bold text-slate-800">
        <ImageIcon className="h-5 w-5 text-sakura-500" /> 사진
      </h1>

      <div className="rounded-2xl bg-white p-4 shadow-sm">
        <label className="block">
          <span className="text-sm text-slate-700">캡션 (선택)</span>
          <input
            value={caption}
            onChange={(e) => setCaption(e.target.value)}
            className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
        <label className="mt-3 flex cursor-pointer items-center justify-center gap-2 rounded-md border-2 border-dashed border-sakura-200 bg-sakura-50/40 px-4 py-6 text-sakura-700 hover:bg-sakura-50">
          <Upload className="h-5 w-5" />
          <span>{uploading ? `업로드 중… ${progress}%` : '사진 선택'}</span>
          <input
            type="file"
            accept="image/*"
            disabled={uploading}
            onChange={(e) => {
              const f = e.target.files?.[0]
              if (f) onPickFile(f)
              e.target.value = ''
            }}
            className="hidden"
          />
        </label>
        {uploading && (
          <div className="mt-2 h-2 w-full overflow-hidden rounded-full bg-slate-200">
            <div
              className="h-full rounded-full bg-sakura-500 transition-all"
              style={{ width: `${progress}%` }}
            />
          </div>
        )}
        {error && <p className="mt-2 text-sm text-red-500">{error}</p>}
      </div>

      {list.error && <p className="text-red-500">{(list.error as Error).message}</p>}

      {list.isLoading && !list.data && <GridSkeleton message="사진을 꺼내오고 있어요" />}

      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4">
        {items.map((m, idx) => (
          <MediaThumb
            key={m.id}
            id={m.id}
            caption={m.caption}
            onDelete={() => deleteMut.mutate(m.id)}
            onClick={() => setLightboxIdx(idx)}
          />
        ))}
        {items.length === 0 && !list.isLoading && (
          <p className="col-span-full rounded-xl bg-white p-6 text-center text-slate-500 shadow-sm">
            아직 사진이 없습니다.
          </p>
        )}
      </div>

      {lightboxIdx !== null && (
        <LightboxWrapper
          items={items}
          startIndex={lightboxIdx}
          onClose={() => setLightboxIdx(null)}
        />
      )}
    </section>
  )
}

function MediaThumb({
  id,
  caption,
  onDelete,
  onClick,
}: {
  id: string
  caption?: string
  onDelete: () => void
  onClick: () => void
}) {
  const url = useQuery({
    queryKey: ['media-url', id],
    queryFn: () => api.get<GetMediaUrlResponse>(`/media/${id}/url?thumbnail=true`),
  })
  return (
    <div
      className="group relative overflow-hidden rounded-xl bg-white shadow-sm cursor-pointer"
      onClick={onClick}
    >
      <div className="aspect-square w-full bg-slate-100">
        {url.data?.url && (
          <img src={url.data.url} alt={caption ?? ''} className="h-full w-full object-cover" />
        )}
      </div>
      {caption && <p className="px-2 py-1 text-xs text-slate-600 line-clamp-1">{caption}</p>}
      <button
        onClick={(e) => {
          e.stopPropagation()
          onDelete()
        }}
        className="absolute right-1 top-1 rounded-full bg-white/80 p-1 text-slate-500 opacity-0 transition-opacity hover:text-red-500 group-hover:opacity-100"
        aria-label="삭제"
      >
        <Trash2 className="h-4 w-4" />
      </button>
    </div>
  )
}

function LightboxWrapper({
  items,
  startIndex,
  onClose,
}: {
  items: { id: string; caption?: string }[]
  startIndex: number
  onClose: () => void
}) {
  const [currentIndex, setCurrentIndex] = useState(startIndex)
  const [urls, setUrls] = useState<Record<string, string>>({})

  useEffect(() => {
    items.forEach((m) => {
      if (!urls[m.id]) {
        api.get<GetMediaUrlResponse>(`/media/${m.id}/url`).then((r) => {
          setUrls((prev) => ({ ...prev, [m.id]: r.url }))
        })
      }
    })
  }, [items])

  const images = items.map((m) => ({
    id: m.id,
    src: urls[m.id] ?? '',
    alt: m.caption,
  }))

  return (
    <Lightbox
      images={images}
      currentIndex={currentIndex}
      onClose={onClose}
      onNavigate={setCurrentIndex}
    />
  )
}

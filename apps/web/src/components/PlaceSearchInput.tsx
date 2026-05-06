import { useEffect, useRef, useState } from 'react'
import { MapPin, Search, Loader2 } from 'lucide-react'
import { api } from '../lib/api'
import type { GeoPoint } from '../lib/types'

export interface PlaceSearchResult {
  placeId: string
  name: string
  address: string
  location: GeoPoint
}

interface GeocodeResponse {
  places?: Array<{
    placeId?: string
    name?: string
    address?: string
    location?: GeoPoint
  }>
}

interface Props {
  label?: string
  value: string
  onChange: (v: string) => void
  onSelect: (place: PlaceSearchResult) => void
  placeholder?: string
}

/**
 * OpenStreetMap Nominatim 기반 장소 검색 자동완성. 백엔드의 `/external/geocode`
 * 엔드포인트를 호출하며, 입력 디바운스(300ms)와 외부 클릭 시 닫힘 처리를 한다.
 * 식당/호텔/관광지 등 모든 종류의 POI 가 검색된다.
 */
export function PlaceSearchInput({ label, value, onChange, onSelect, placeholder }: Props) {
  const [query, setQuery] = useState(value)
  const [results, setResults] = useState<PlaceSearchResult[]>([])
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const wrapRef = useRef<HTMLDivElement | null>(null)

  // 외부 value 변경에 동기화 (예: 폼 리셋).
  useEffect(() => {
    setQuery(value)
  }, [value])

  // 디바운스된 검색.
  useEffect(() => {
    const q = query.trim()
    if (q.length < 2) {
      setResults([])
      setLoading(false)
      return
    }
    setLoading(true)
    const handle = window.setTimeout(async () => {
      try {
        const res = await api.get<GeocodeResponse>(`/external/geocode?q=${encodeURIComponent(q)}`)
        const items: PlaceSearchResult[] = (res.places ?? [])
          .filter((p) => p.location && typeof p.location.latitude === 'number')
          .map((p) => ({
            placeId: p.placeId ?? '',
            name: p.name ?? p.address ?? '',
            address: p.address ?? '',
            location: p.location as GeoPoint,
          }))
        setResults(items)
        setOpen(true)
      } catch {
        setResults([])
      } finally {
        setLoading(false)
      }
    }, 300)
    return () => window.clearTimeout(handle)
  }, [query])

  // 바깥 클릭 시 닫기.
  useEffect(() => {
    const onDown = (e: MouseEvent) => {
      if (!wrapRef.current) return
      if (!wrapRef.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', onDown)
    return () => document.removeEventListener('mousedown', onDown)
  }, [])

  return (
    <div ref={wrapRef} className="relative">
      {label && <label className="mb-1 block text-sm font-medium text-slate-700">{label}</label>}
      <div className="relative">
        <span className="pointer-events-none absolute inset-y-0 left-2 flex items-center text-slate-400">
          <Search className="h-4 w-4" />
        </span>
        <input
          type="text"
          value={query}
          placeholder={placeholder ?? '식당, 호텔, 관광지 검색'}
          onChange={(e) => {
            setQuery(e.target.value)
            onChange(e.target.value)
          }}
          onFocus={() => results.length > 0 && setOpen(true)}
          className="w-full rounded-md border border-slate-300 bg-white px-3 py-2 pl-8 text-sm focus:border-sakura-500 focus:outline-none"
        />
        {loading && (
          <span className="absolute inset-y-0 right-2 flex items-center text-slate-400">
            <Loader2 className="h-4 w-4 animate-spin" />
          </span>
        )}
      </div>
      {open && results.length > 0 && (
        <ul className="absolute z-20 mt-1 max-h-72 w-full overflow-y-auto rounded-md border border-slate-200 bg-white shadow-lg">
          {results.map((r) => (
            <li key={r.placeId}>
              <button
                type="button"
                onClick={() => {
                  setQuery(r.name)
                  onChange(r.name)
                  onSelect(r)
                  setOpen(false)
                }}
                className="flex w-full items-start gap-2 px-3 py-2 text-left text-sm hover:bg-sakura-50"
              >
                <MapPin className="mt-0.5 h-4 w-4 flex-shrink-0 text-sakura-500" />
                <span className="flex-1">
                  <span className="block font-medium text-slate-800">{r.name}</span>
                  {r.address && r.address !== r.name && (
                    <span className="block text-xs text-slate-500">{r.address}</span>
                  )}
                </span>
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

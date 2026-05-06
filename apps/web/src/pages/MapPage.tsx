import { useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ArrowLeft } from 'lucide-react'
import { api } from '../lib/api'
import type { ListDaysResponse, ListSchedulesResponse, Trip } from '../lib/types'

interface Marker {
  lat: number
  lng: number
  title: string
}

interface GeocodeHit {
  name: string
  address: string
  latitude: number
  longitude: number
}

interface GeocodeResponse {
  results: GeocodeHit[]
}

const LEAFLET_CSS = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
const LEAFLET_JS = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'

async function ensureLeaflet(): Promise<any> {
  const w = window as any
  if (w.L) return w.L
  if (!document.querySelector(`link[href="${LEAFLET_CSS}"]`)) {
    const link = document.createElement('link')
    link.rel = 'stylesheet'
    link.href = LEAFLET_CSS
    document.head.appendChild(link)
  }
  if (!document.querySelector(`script[src="${LEAFLET_JS}"]`)) {
    await new Promise<void>((resolve, reject) => {
      const s = document.createElement('script')
      s.src = LEAFLET_JS
      s.async = true
      s.onload = () => resolve()
      s.onerror = () => reject(new Error('failed to load leaflet'))
      document.head.appendChild(s)
    })
  } else {
    await new Promise<void>((resolve) => {
      const check = () => ((window as any).L ? resolve() : setTimeout(check, 50))
      check()
    })
  }
  return (window as any).L
}

export function MapPage() {
  const { tripId = '' } = useParams()
  const containerRef = useRef<HTMLDivElement | null>(null)
  const mapRef = useRef<any>(null)
  const [markers, setMarkers] = useState<Marker[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  const trip = useQuery({
    queryKey: ['trip', tripId],
    queryFn: () => api.get<{ trip: Trip }>(`/trips/${tripId}`),
    enabled: !!tripId,
  })

  useEffect(() => {
    if (!tripId) return
    let cancelled = false
    ;(async () => {
      try {
        const daysRes = await api.get<ListDaysResponse>(`/trips/${tripId}/days`)
        const collected: Marker[] = []
        for (const d of daysRes.days ?? []) {
          try {
            const schedRes = await api.get<ListSchedulesResponse>(`/days/${d.id}/schedules`)
            for (const s of schedRes.schedules ?? []) {
              const lat = s.location?.latitude
              const lng = s.location?.longitude
              if (lat !== undefined && lng !== undefined && (lat !== 0 || lng !== 0)) {
                collected.push({
                  lat,
                  lng,
                  title: s.placeName ? `${s.title} · ${s.placeName}` : s.title,
                })
              }
            }
          } catch {
            // 한 day 실패는 무시
          }
        }
        if (collected.length === 0) {
          const dest = trip.data?.trip?.destination
          if (dest) {
            try {
              const geo = await api.get<GeocodeResponse>(
                `/external/geocode?q=${encodeURIComponent(dest)}`,
              )
              const hit = geo.results?.[0]
              if (hit) {
                collected.push({
                  lat: hit.latitude,
                  lng: hit.longitude,
                  title: hit.name || dest,
                })
              }
            } catch {
              // ignore
            }
          }
        }
        if (!cancelled) setMarkers(collected)
      } catch (e) {
        if (!cancelled) setError((e as Error).message)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [tripId, trip.data?.trip?.destination])

  useEffect(() => {
    if (!markers || !containerRef.current) return
    let active = true
    ;(async () => {
      try {
        const L = await ensureLeaflet()
        if (!active || !containerRef.current) return
        if (mapRef.current) {
          mapRef.current.remove()
          mapRef.current = null
        }
        const center: [number, number] =
          markers.length > 0 ? [markers[0].lat, markers[0].lng] : [35.0, 138.0]
        const zoom = markers.length > 0 ? 11 : 5
        const map = L.map(containerRef.current).setView(center, zoom)
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          attribution: '© OpenStreetMap contributors',
          maxZoom: 19,
        }).addTo(map)
        const layer = L.featureGroup()
        markers.forEach((m) => {
          L.marker([m.lat, m.lng]).bindPopup(m.title).addTo(layer)
        })
        if (markers.length > 0) {
          layer.addTo(map)
          if (markers.length > 1) {
            map.fitBounds(layer.getBounds().pad(0.2))
          }
        }
        mapRef.current = map
      } catch (e) {
        if (active) setError((e as Error).message)
      }
    })()
    return () => {
      active = false
      if (mapRef.current) {
        mapRef.current.remove()
        mapRef.current = null
      }
    }
  }, [markers])

  return (
    <section className="space-y-3">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <div className="rounded-2xl bg-white p-4 shadow-sm">
        <h1 className="text-xl font-bold text-slate-800">지도</h1>
        {trip.data?.trip?.title && <p className="text-sm text-slate-500">{trip.data.trip.title}</p>}
      </div>
      {error && <p className="text-red-500 text-sm">{error}</p>}
      {markers && markers.length === 0 && !error && (
        <p className="rounded-xl bg-white p-4 text-sm text-slate-500 shadow-sm">
          좌표가 등록된 일정이 없습니다. 일정 추가/수정에서 장소 검색으로 좌표를 등록하면 지도에
          표시됩니다.
        </p>
      )}
      <div
        ref={containerRef}
        style={{ height: '70vh', width: '100%' }}
        className="rounded-2xl overflow-hidden shadow-sm bg-slate-100"
      />
    </section>
  )
}

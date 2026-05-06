import { useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import { getToken } from '../lib/keycloak'

/**
 * 위치 공유 트래킹 페이지.
 *
 * - 양쪽 모두 동의(시작 버튼)했을 때만 위치를 송신한다.
 * - WebSocket 으로 백엔드 ws hub 에 접속한다.
 *   `wss://<api-host>/api/v1/ws/trips/<tripId>?token=<bearer>`
 * - 자기 위치는 navigator.geolocation.watchPosition 으로 추적해 ws 로 publish.
 * - 동행자 위치는 ws 메시지로 수신해 지도에 마커로 표시.
 */

interface PeerLocation {
  memberId: string
  displayName?: string
  lat: number
  lng: number
  updatedAt: number
}

interface IncomingMessage {
  type: 'location'
  memberId: string
  displayName?: string
  lat: number
  lng: number
  ts?: number
}

const LEAFLET_CSS = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
const LEAFLET_JS = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'
const FALLBACK_CENTER: [number, number] = [35.6895, 139.6917]

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

function buildWsUrl(tripId: string, token: string): string {
  const base = (import.meta.env.VITE_API_BASE ?? 'https://journey-api.seonology.com').replace(
    /^http/,
    'ws',
  )
  return `${base}/api/v1/ws/trips/${encodeURIComponent(tripId)}?token=${encodeURIComponent(token)}`
}

export function TrackingPage() {
  const { tripId = '' } = useParams<{ tripId: string }>()
  const containerRef = useRef<HTMLDivElement | null>(null)
  const mapRef = useRef<any>(null)
  const meMarkerRef = useRef<any>(null)
  const peerMarkersRef = useRef<Record<string, any>>({})
  const wsRef = useRef<WebSocket | null>(null)
  const watchIdRef = useRef<number | null>(null)

  const [active, setActive] = useState(false)
  const [me, setMe] = useState<{ lat: number; lng: number } | null>(null)
  const [peers, setPeers] = useState<Record<string, PeerLocation>>({})
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let alive = true
    ;(async () => {
      try {
        const L = await ensureLeaflet()
        if (!alive || !containerRef.current) return
        if (mapRef.current) return
        const map = L.map(containerRef.current).setView(FALLBACK_CENTER, 12)
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          attribution: '© OpenStreetMap contributors',
          maxZoom: 19,
        }).addTo(map)
        mapRef.current = map
      } catch (e) {
        if (alive) setError((e as Error).message)
      }
    })()
    return () => {
      alive = false
      if (mapRef.current) {
        mapRef.current.remove()
        mapRef.current = null
      }
    }
  }, [])

  useEffect(() => {
    if (!me || !mapRef.current) return
    const L = (window as any).L
    if (!L) return
    if (meMarkerRef.current) {
      meMarkerRef.current.setLatLng([me.lat, me.lng])
    } else {
      meMarkerRef.current = L.circleMarker([me.lat, me.lng], {
        radius: 10,
        color: '#0ea5e9',
        fillColor: '#0ea5e9',
        fillOpacity: 0.9,
      })
        .bindPopup('나')
        .addTo(mapRef.current)
      mapRef.current.setView([me.lat, me.lng], 14)
    }
  }, [me])

  useEffect(() => {
    if (!mapRef.current) return
    const L = (window as any).L
    if (!L) return
    Object.values(peers).forEach((p) => {
      const m = peerMarkersRef.current[p.memberId]
      if (m) {
        m.setLatLng([p.lat, p.lng])
        m.bindPopup(p.displayName ?? p.memberId)
      } else {
        peerMarkersRef.current[p.memberId] = L.circleMarker([p.lat, p.lng], {
          radius: 10,
          color: '#ec4899',
          fillColor: '#ec4899',
          fillOpacity: 0.9,
        })
          .bindPopup(p.displayName ?? p.memberId)
          .addTo(mapRef.current)
      }
    })
  }, [peers])

  function stop() {
    setActive(false)
    if (watchIdRef.current !== null) {
      navigator.geolocation.clearWatch(watchIdRef.current)
      watchIdRef.current = null
    }
    wsRef.current?.close()
    wsRef.current = null
  }

  async function start() {
    setError(null)
    if (!('geolocation' in navigator)) {
      setError('이 브라우저는 위치 정보를 지원하지 않습니다.')
      return
    }
    const token = await getToken()
    if (!token) {
      setError('로그인 정보가 없습니다.')
      return
    }
    let ws: WebSocket
    try {
      ws = new WebSocket(buildWsUrl(tripId, token))
    } catch (err) {
      setError(err instanceof Error ? err.message : '연결 실패')
      return
    }
    wsRef.current = ws
    ws.onmessage = (ev) => {
      try {
        const msg: IncomingMessage = JSON.parse(typeof ev.data === 'string' ? ev.data : '')
        if (msg.type === 'location') {
          setPeers((prev) => ({
            ...prev,
            [msg.memberId]: {
              memberId: msg.memberId,
              displayName: msg.displayName,
              lat: msg.lat,
              lng: msg.lng,
              updatedAt: msg.ts ?? Date.now(),
            },
          }))
        }
      } catch {
        // ignore
      }
    }
    ws.onerror = () => setError('WebSocket 오류가 발생했습니다.')

    setActive(true)

    watchIdRef.current = navigator.geolocation.watchPosition(
      (pos) => {
        const lat = pos.coords.latitude
        const lng = pos.coords.longitude
        setMe({ lat, lng })
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: 'location', lat, lng, ts: Date.now() }))
        }
      },
      (err) => setError(err.message),
      { enableHighAccuracy: true, maximumAge: 5000, timeout: 15000 },
    )
  }

  useEffect(() => {
    return () => stop()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <section className="space-y-3">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <div className="rounded-2xl bg-white p-4 shadow-sm">
        <div className="flex items-center justify-between">
          <h1 className="text-xl font-bold text-slate-800">위치 트래킹</h1>
          {!active ? (
            <button
              onClick={start}
              className="rounded-full bg-sakura-500 px-4 py-2 text-sm font-bold text-white shadow"
            >
              위치 공유 시작
            </button>
          ) : (
            <button
              onClick={stop}
              className="rounded-full bg-slate-700 px-4 py-2 text-sm font-bold text-white shadow"
            >
              중지
            </button>
          )}
        </div>
        <p className="mt-2 text-xs text-slate-500">
          상호 동의 기반 실시간 위치 공유. 시작 버튼을 누른 시점부터 본인 위치가 동행자에게
          공유되며, 중지하면 더 이상 전송되지 않습니다.
        </p>
      </div>

      {error && <p className="rounded-xl bg-red-50 p-3 text-sm text-red-600">{error}</p>}

      <div
        ref={containerRef}
        style={{ height: '60vh', width: '100%' }}
        className="rounded-2xl overflow-hidden shadow-sm bg-slate-100"
      />

      <div className="rounded-2xl bg-white p-4 shadow-sm">
        <h2 className="text-sm font-bold text-slate-700">동행자 위치</h2>
        {Object.values(peers).length === 0 ? (
          <p className="mt-2 text-xs text-slate-500">아직 수신된 위치가 없습니다.</p>
        ) : (
          <ul className="mt-2 space-y-1 text-xs">
            {Object.values(peers).map((p) => (
              <li
                key={p.memberId}
                className="rounded-md border border-sakura-100 bg-white px-2 py-1"
              >
                <span className="font-bold text-slate-800">{p.displayName ?? p.memberId}</span>
                <span className="ml-2 text-slate-500">
                  {p.lat.toFixed(5)}, {p.lng.toFixed(5)}
                </span>
                <span className="ml-2 text-slate-400">
                  {new Date(p.updatedAt).toLocaleTimeString()}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  )
}

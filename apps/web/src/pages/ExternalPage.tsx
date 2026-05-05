import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { ArrowLeft, MapPin, Navigation, Search } from 'lucide-react'
import { api } from '../lib/api'
import { MapPicker } from '../components/MapPicker'

interface PlaceResult {
  name: string
  address: string
  lat: number
  lng: number
  rating?: number
  types?: string[]
}

interface NearbyResponse {
  results: PlaceResult[]
}

interface RouteResult {
  summary: string
  duration: string
  distance: string
  steps: { instruction: string; distance: string }[]
}

interface TransitResponse {
  routes: RouteResult[]
}

export function NearbyPage() {
  const { tripId = '' } = useParams()
  const [coords, setCoords] = useState<{ lat: number; lng: number } | null>(null)
  const [radius, setRadius] = useState(1000)
  const [placeType, setPlaceType] = useState('restaurant')

  const nearbyMut = useMutation({
    mutationFn: () =>
      api.get<NearbyResponse>(
        `/api/v1/nearby?lat=${coords!.lat}&lng=${coords!.lng}&radius=${radius}&type=${placeType}`,
      ),
  })

  const results = nearbyMut.data?.results ?? []

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <h1 className="flex items-center gap-2 text-xl font-bold text-slate-800">
        <MapPin className="h-5 w-5 text-sakura-500" /> 주변 검색
      </h1>

      <div className="rounded-2xl bg-white p-4 shadow-sm space-y-3">
        <MapPicker
          latitude={coords?.lat}
          longitude={coords?.lng}
          onChange={setCoords}
          accessToken={import.meta.env.VITE_MAPBOX_TOKEN}
        />
        {coords && (
          <p className="text-xs text-slate-500">
            선택: {coords.lat.toFixed(5)}, {coords.lng.toFixed(5)}
          </p>
        )}
        <div className="grid grid-cols-2 gap-2">
          <label className="block">
            <span className="text-sm text-slate-700">반경 (m)</span>
            <input
              type="number"
              value={radius}
              onChange={(e) => setRadius(Number(e.target.value))}
              className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            />
          </label>
          <label className="block">
            <span className="text-sm text-slate-700">장소 유형</span>
            <select
              value={placeType}
              onChange={(e) => setPlaceType(e.target.value)}
              className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
            >
              <option value="restaurant">식당</option>
              <option value="cafe">카페</option>
              <option value="tourist_attraction">관광지</option>
              <option value="convenience_store">편의점</option>
              <option value="hotel">호텔</option>
            </select>
          </label>
        </div>
        <button
          onClick={() => nearbyMut.mutate()}
          disabled={!coords || nearbyMut.isPending}
          className="flex w-full items-center justify-center gap-2 rounded-md bg-sakura-500 px-3 py-2 text-white hover:bg-sakura-600 disabled:opacity-50"
        >
          <Search className="h-4 w-4" />
          {nearbyMut.isPending ? '검색 중…' : '검색'}
        </button>
      </div>

      {results.length > 0 && (
        <ul className="space-y-2">
          {results.map((p, i) => (
            <li key={i} className="rounded-xl bg-white p-3 shadow-sm">
              <p className="font-bold text-slate-800">{p.name}</p>
              <p className="text-xs text-slate-500">{p.address}</p>
              {p.rating && <p className="text-xs text-amber-600">평점 {p.rating}</p>}
            </li>
          ))}
        </ul>
      )}

      {nearbyMut.isError && (
        <p className="text-sm text-red-500">{(nearbyMut.error as Error).message}</p>
      )}
    </section>
  )
}

export function TransitPage() {
  const { tripId = '' } = useParams()
  const [origin, setOrigin] = useState('')
  const [destination, setDestination] = useState('')
  const [departureTime, setDepartureTime] = useState('')

  const transitMut = useMutation({
    mutationFn: () => {
      const [oLat, oLng] = origin.split(',').map(Number)
      const [dLat, dLng] = destination.split(',').map(Number)
      let url = `/api/v1/transit?origin_lat=${oLat}&origin_lng=${oLng}&dest_lat=${dLat}&dest_lng=${dLng}`
      if (departureTime) url += `&departure_time=${encodeURIComponent(departureTime)}`
      return api.get<TransitResponse>(url)
    },
  })

  const routes = transitMut.data?.routes ?? []

  return (
    <section className="space-y-4">
      <Link
        to={`/trips/${tripId}`}
        className="inline-flex items-center gap-1 text-sm text-sakura-700"
      >
        <ArrowLeft className="h-4 w-4" /> 여행으로
      </Link>
      <h1 className="flex items-center gap-2 text-xl font-bold text-slate-800">
        <Navigation className="h-5 w-5 text-sakura-500" /> 교통 검색
      </h1>

      <div className="rounded-2xl bg-white p-4 shadow-sm space-y-3">
        <label className="block">
          <span className="text-sm text-slate-700">출발지 (lat,lng)</span>
          <input
            value={origin}
            onChange={(e) => setOrigin(e.target.value)}
            placeholder="35.6895,139.6917"
            className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
        <label className="block">
          <span className="text-sm text-slate-700">도착지 (lat,lng)</span>
          <input
            value={destination}
            onChange={(e) => setDestination(e.target.value)}
            placeholder="35.6762,139.6503"
            className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
        <label className="block">
          <span className="text-sm text-slate-700">출발 시간 (선택)</span>
          <input
            type="datetime-local"
            value={departureTime}
            onChange={(e) => setDepartureTime(e.target.value)}
            className="mt-1 w-full rounded-md border border-slate-200 px-3 py-2"
          />
        </label>
        <button
          onClick={() => transitMut.mutate()}
          disabled={!origin || !destination || transitMut.isPending}
          className="flex w-full items-center justify-center gap-2 rounded-md bg-sakura-500 px-3 py-2 text-white hover:bg-sakura-600 disabled:opacity-50"
        >
          <Search className="h-4 w-4" />
          {transitMut.isPending ? '검색 중…' : '경로 검색'}
        </button>
      </div>

      {routes.length > 0 && (
        <ul className="space-y-3">
          {routes.map((r, i) => (
            <li key={i} className="rounded-xl bg-white p-4 shadow-sm">
              <div className="flex items-center justify-between">
                <p className="font-bold text-slate-800">{r.summary || `경로 ${i + 1}`}</p>
                <span className="text-sm text-slate-500">{r.duration}</span>
              </div>
              <p className="text-xs text-slate-500">{r.distance}</p>
              {r.steps && (
                <ol className="mt-2 space-y-1">
                  {r.steps.map((step, j) => (
                    <li key={j} className="text-xs text-slate-600">
                      {j + 1}. {step.instruction} ({step.distance})
                    </li>
                  ))}
                </ol>
              )}
            </li>
          ))}
        </ul>
      )}

      {transitMut.isError && (
        <p className="text-sm text-red-500">{(transitMut.error as Error).message}</p>
      )}
    </section>
  )
}

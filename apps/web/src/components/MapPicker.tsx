import { useEffect, useRef } from 'react'
import { MapPin } from 'lucide-react'

interface MapPickerProps {
  latitude?: number
  longitude?: number
  onChange?: (coords: { lat: number; lng: number }) => void
  readOnly?: boolean
  className?: string
  accessToken?: string
}

export function MapPicker({
  latitude,
  longitude,
  onChange,
  readOnly = false,
  className = '',
  accessToken,
}: MapPickerProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const mapRef = useRef<any>(null)

  useEffect(() => {
    if (!containerRef.current || !accessToken) return
    if (mapRef.current) return

    // mapbox-gl 은 번들에서 제외되어 있으므로 런타임에 window 전역 또는
    // 동적 import 로 찾는다. 둘 다 없으면 조용히 포기하고 fallback UI 를
    // 그대로 둔다 (아래 조건 분기가 이미 덮고 있음).
    const loadMap = async () => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const globalAny = window as any
      let mapboxgl = globalAny.mapboxgl
      if (!mapboxgl) {
        try {
          const mod = await import(/* @vite-ignore */ 'mapbox-gl')
          mapboxgl = mod.default || mod
        } catch {
          return
        }
      }
      mapboxgl.accessToken = accessToken
      const map = new mapboxgl.Map({
        container: containerRef.current!,
        style: 'mapbox://styles/mapbox/streets-v12',
        center: [longitude ?? 139.6917, latitude ?? 35.6895],
        zoom: 13,
      })

      if (!readOnly) {
        map.on('click', (e: { lngLat: { lat: number; lng: number } }) => {
          onChange?.({ lat: e.lngLat.lat, lng: e.lngLat.lng })
        })
      }

      mapRef.current = map
    }
    loadMap()

    return () => {
      mapRef.current?.remove()
      mapRef.current = null
    }
  }, [accessToken])

  if (!accessToken) {
    return (
      <div
        className={`flex h-48 items-center justify-center rounded-lg border border-dashed border-gray-300 dark:border-gray-600 ${className}`}
      >
        <div className="flex flex-col items-center text-gray-400">
          <MapPin className="h-8 w-8" />
          <span className="mt-1 text-sm">Map (token required)</span>
          {latitude && longitude && (
            <span className="mt-1 text-xs">
              {latitude.toFixed(4)}, {longitude.toFixed(4)}
            </span>
          )}
        </div>
      </div>
    )
  }

  return <div ref={containerRef} className={`h-48 w-full rounded-lg ${className}`} />
}

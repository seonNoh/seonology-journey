interface MasonryGalleryProps {
  images: { id: string; src: string; alt?: string }[]
  columns?: number
  gap?: number
  onImageClick?: (id: string) => void
  className?: string
}

export function MasonryGallery({
  images,
  columns = 3,
  gap = 4,
  onImageClick,
  className = '',
}: MasonryGalleryProps) {
  const cols: (typeof images)[] = Array.from({ length: columns }, () => [])
  images.forEach((img, i) => {
    cols[i % columns].push(img)
  })

  return (
    <div
      className={`grid ${className}`}
      style={{ gridTemplateColumns: `repeat(${columns}, 1fr)`, gap: `${gap * 4}px` }}
    >
      {cols.map((col, colIdx) => (
        <div key={colIdx} className="flex flex-col" style={{ gap: `${gap * 4}px` }}>
          {col.map((img) => (
            <button
              key={img.id}
              type="button"
              onClick={() => onImageClick?.(img.id)}
              className="overflow-hidden rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <img
                src={img.src}
                alt={img.alt || ''}
                className="w-full object-cover transition-transform hover:scale-105"
                loading="lazy"
              />
            </button>
          ))}
        </div>
      ))}
    </div>
  )
}

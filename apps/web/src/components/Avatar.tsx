interface AvatarProps {
  src?: string
  name: string
  size?: 'sm' | 'md' | 'lg'
}

const sizeMap = { sm: 'h-8 w-8 text-xs', md: 'h-10 w-10 text-sm', lg: 'h-14 w-14 text-base' }

export function Avatar({ src, name, size = 'md' }: AvatarProps) {
  const initials = name
    .split(' ')
    .map((s) => s[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()

  if (src) {
    return (
      <img
        src={src}
        alt={name}
        className={`${sizeMap[size]} rounded-full object-cover ring-2 ring-white`}
      />
    )
  }

  return (
    <div
      className={`${sizeMap[size]} rounded-full bg-sakura-100 text-sakura-700 flex items-center justify-center font-medium ring-2 ring-white`}
      aria-label={name}
    >
      {initials}
    </div>
  )
}

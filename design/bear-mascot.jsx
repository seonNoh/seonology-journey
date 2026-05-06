// Reusable bear mascot SVGs — original cute round bear, NOT Rilakkuma.
// Variants take a palette so the same shape can re-skin per option.

const BearMascot = ({ size = 80, palette }) => {
  const { fur, ear, snout, outline = '#3b2a20', cheek = '#f5b9b9' } = palette
  return (
    <svg width={size} height={size} viewBox="0 0 100 100" style={{ display: 'block' }}>
      {/* ears */}
      <circle cx="22" cy="28" r="13" fill={fur} stroke={outline} strokeWidth="1.5" />
      <circle cx="78" cy="28" r="13" fill={fur} stroke={outline} strokeWidth="1.5" />
      <circle cx="22" cy="28" r="7" fill={ear} />
      <circle cx="78" cy="28" r="7" fill={ear} />
      {/* head */}
      <ellipse cx="50" cy="56" rx="34" ry="32" fill={fur} stroke={outline} strokeWidth="1.5" />
      {/* snout */}
      <ellipse cx="50" cy="65" rx="18" ry="14" fill={snout} />
      {/* cheeks */}
      <circle cx="22" cy="62" r="5" fill={cheek} opacity="0.6" />
      <circle cx="78" cy="62" r="5" fill={cheek} opacity="0.6" />
      {/* eyes */}
      <ellipse cx="38" cy="50" rx="2.6" ry="3.4" fill={outline} />
      <ellipse cx="62" cy="50" rx="2.6" ry="3.4" fill={outline} />
      <circle cx="38.7" cy="49" r="0.8" fill="#fff" />
      <circle cx="62.7" cy="49" r="0.8" fill="#fff" />
      {/* nose */}
      <ellipse cx="50" cy="60" rx="3" ry="2.2" fill={outline} />
      {/* mouth */}
      <path
        d="M 50 62 Q 46 68 42 65"
        fill="none"
        stroke={outline}
        strokeWidth="1.4"
        strokeLinecap="round"
      />
      <path
        d="M 50 62 Q 54 68 58 65"
        fill="none"
        stroke={outline}
        strokeWidth="1.4"
        strokeLinecap="round"
      />
    </svg>
  )
}

// tiny chick (yellow round bird) — placeholder for the small companion
const ChickMascot = ({
  size = 40,
  palette = { body: '#FFD66B', beak: '#E89A3C', outline: '#3b2a20' },
}) => {
  const { body, beak, outline } = palette
  return (
    <svg width={size} height={size} viewBox="0 0 100 100" style={{ display: 'block' }}>
      <ellipse cx="50" cy="55" rx="34" ry="30" fill={body} stroke={outline} strokeWidth="1.5" />
      {/* tiny tuft */}
      <path
        d="M 42 22 L 50 14 L 58 22"
        fill="none"
        stroke={outline}
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      {/* eyes */}
      <circle cx="42" cy="50" r="2.2" fill={outline} />
      <circle cx="58" cy="50" r="2.2" fill={outline} />
      {/* beak */}
      <path d="M 46 60 L 54 60 L 50 66 Z" fill={beak} stroke={outline} strokeWidth="1" />
      {/* cheeks */}
      <circle cx="32" cy="58" r="3" fill="#f5b9b9" opacity="0.7" />
      <circle cx="68" cy="58" r="3" fill="#f5b9b9" opacity="0.7" />
    </svg>
  )
}

// small brown bear (companion)
const BrownBearMascot = ({
  size = 40,
  palette = { fur: '#A57156', ear: '#7A4F3A', snout: '#E8C9B0', outline: '#3b2a20' },
}) => {
  return <BearMascot size={size} palette={{ ...palette, cheek: '#f5b9b9' }} />
}

// Sleeping bear face (eyes closed)
const SleepingBear = ({ size = 80, palette }) => {
  const { fur, ear, snout, outline = '#3b2a20', cheek = '#f5b9b9' } = palette
  return (
    <svg width={size} height={size} viewBox="0 0 100 100" style={{ display: 'block' }}>
      <circle cx="22" cy="28" r="13" fill={fur} stroke={outline} strokeWidth="1.5" />
      <circle cx="78" cy="28" r="13" fill={fur} stroke={outline} strokeWidth="1.5" />
      <circle cx="22" cy="28" r="7" fill={ear} />
      <circle cx="78" cy="28" r="7" fill={ear} />
      <ellipse cx="50" cy="56" rx="34" ry="32" fill={fur} stroke={outline} strokeWidth="1.5" />
      <ellipse cx="50" cy="65" rx="18" ry="14" fill={snout} />
      <circle cx="22" cy="62" r="5" fill={cheek} opacity="0.6" />
      <circle cx="78" cy="62" r="5" fill={cheek} opacity="0.6" />
      {/* closed eyes ‿ ‿ */}
      <path
        d="M 34 50 Q 38 53 42 50"
        fill="none"
        stroke={outline}
        strokeWidth="1.6"
        strokeLinecap="round"
      />
      <path
        d="M 58 50 Q 62 53 66 50"
        fill="none"
        stroke={outline}
        strokeWidth="1.6"
        strokeLinecap="round"
      />
      <ellipse cx="50" cy="60" rx="3" ry="2.2" fill={outline} />
      <path
        d="M 47 64 Q 50 67 53 64"
        fill="none"
        stroke={outline}
        strokeWidth="1.4"
        strokeLinecap="round"
      />
    </svg>
  )
}

Object.assign(window, { BearMascot, ChickMascot, BrownBearMascot, SleepingBear })

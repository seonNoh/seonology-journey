// Expanded mascot set for Sakura Bear theme — original characters.
// Three personae: MainBear (beige bear), MiniBear (small brown), Chick (yellow chick).
// Each has multiple expressions.

const SB_BEAR_PALETTE = {
  fur: '#F4D6B8',
  ear: '#D9A179',
  snout: '#FBEAD5',
  outline: '#5A2230',
  cheek: '#FFB3C7',
}
const SB_MINI_PALETTE = {
  fur: '#B07A56',
  ear: '#7E4F33',
  snout: '#E8C9B0',
  outline: '#3B1F1A',
  cheek: '#FFB3C7',
}
const SB_CHICK_PALETTE = {
  body: '#FFD66B',
  beak: '#E89A3C',
  outline: '#5A2230',
}

// expressions: 'happy' | 'sleep' | 'surprise' | 'wink' | 'plain'
function bearEyes(expr, outline) {
  switch (expr) {
    case 'sleep':
      return (
        <>
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
        </>
      )
    case 'happy':
      return (
        <>
          <path
            d="M 34 51 Q 38 47 42 51"
            fill="none"
            stroke={outline}
            strokeWidth="1.8"
            strokeLinecap="round"
          />
          <path
            d="M 58 51 Q 62 47 66 51"
            fill="none"
            stroke={outline}
            strokeWidth="1.8"
            strokeLinecap="round"
          />
        </>
      )
    case 'wink':
      return (
        <>
          <path
            d="M 34 51 Q 38 47 42 51"
            fill="none"
            stroke={outline}
            strokeWidth="1.8"
            strokeLinecap="round"
          />
          <ellipse cx="62" cy="50" rx="2.6" ry="3.4" fill={outline} />
          <circle cx="62.7" cy="49" r="0.8" fill="#fff" />
        </>
      )
    case 'surprise':
      return (
        <>
          <circle cx="38" cy="50" r="3" fill={outline} />
          <circle cx="62" cy="50" r="3" fill={outline} />
          <circle cx="38.7" cy="49" r="1" fill="#fff" />
          <circle cx="62.7" cy="49" r="1" fill="#fff" />
        </>
      )
    default:
      return (
        <>
          <ellipse cx="38" cy="50" rx="2.6" ry="3.4" fill={outline} />
          <ellipse cx="62" cy="50" rx="2.6" ry="3.4" fill={outline} />
          <circle cx="38.7" cy="49" r="0.8" fill="#fff" />
          <circle cx="62.7" cy="49" r="0.8" fill="#fff" />
        </>
      )
  }
}

function bearMouth(expr, outline) {
  if (expr === 'surprise') {
    return <ellipse cx="50" cy="66" rx="2.5" ry="3" fill={outline} />
  }
  if (expr === 'happy') {
    return (
      <path
        d="M 44 64 Q 50 70 56 64"
        fill="none"
        stroke={outline}
        strokeWidth="1.6"
        strokeLinecap="round"
      />
    )
  }
  return (
    <>
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
    </>
  )
}

const SBBear = ({ size = 80, expr = 'plain', palette = SB_BEAR_PALETTE, accessory }) => {
  const { fur, ear, snout, outline, cheek } = palette
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 100 100"
      style={{ display: 'block', overflow: 'visible' }}
    >
      {/* ears */}
      <circle cx="22" cy="28" r="13" fill={fur} stroke={outline} strokeWidth="1.5" />
      <circle cx="78" cy="28" r="13" fill={fur} stroke={outline} strokeWidth="1.5" />
      <circle cx="22" cy="28" r="7" fill={ear} />
      <circle cx="78" cy="28" r="7" fill={ear} />
      {/* head */}
      <ellipse cx="50" cy="56" rx="34" ry="32" fill={fur} stroke={outline} strokeWidth="1.5" />
      <ellipse cx="50" cy="65" rx="18" ry="14" fill={snout} />
      <circle cx="22" cy="62" r="5" fill={cheek} opacity="0.65" />
      <circle cx="78" cy="62" r="5" fill={cheek} opacity="0.65" />
      {bearEyes(expr, outline)}
      <ellipse cx="50" cy="60" rx="3" ry="2.2" fill={outline} />
      {bearMouth(expr, outline)}
      {accessory === 'flower' && (
        <g transform="translate(72,8)">
          <circle r="4" fill="#F25C7A" />
          <circle cx="-4" cy="-2" r="3" fill="#FFA3B8" />
          <circle cx="4" cy="-2" r="3" fill="#FFA3B8" />
          <circle cx="-2" cy="3" r="3" fill="#FFA3B8" />
          <circle cx="2" cy="3" r="3" fill="#FFA3B8" />
          <circle r="1.5" fill="#fff" />
        </g>
      )}
      {accessory === 'sleepcap' && (
        <g>
          <path
            d="M 18 14 Q 50 -8 82 14 L 78 22 Q 50 6 22 22 Z"
            fill="#F25C7A"
            stroke={outline}
            strokeWidth="1.2"
          />
          <circle cx="86" cy="10" r="4" fill="#fff" stroke={outline} strokeWidth="1.2" />
        </g>
      )}
      {accessory === 'camera' && (
        <g transform="translate(64,72)">
          <rect x="-10" y="-6" width="20" height="14" rx="3" fill="#5A2230" />
          <circle r="4" fill="#FFB3C7" />
          <circle r="2" fill="#5A2230" />
        </g>
      )}
    </svg>
  )
}

const SBChick = ({ size = 40, expr = 'plain' }) => {
  const { body, beak, outline } = SB_CHICK_PALETTE
  const eyes =
    expr === 'happy' ? (
      <>
        <path
          d="M 38 49 Q 42 46 46 49"
          fill="none"
          stroke={outline}
          strokeWidth="1.6"
          strokeLinecap="round"
        />
        <path
          d="M 54 49 Q 58 46 62 49"
          fill="none"
          stroke={outline}
          strokeWidth="1.6"
          strokeLinecap="round"
        />
      </>
    ) : (
      <>
        <circle cx="42" cy="50" r="2.2" fill={outline} />
        <circle cx="58" cy="50" r="2.2" fill={outline} />
      </>
    )
  return (
    <svg width={size} height={size} viewBox="0 0 100 100" style={{ display: 'block' }}>
      <ellipse cx="50" cy="55" rx="34" ry="30" fill={body} stroke={outline} strokeWidth="1.5" />
      <path
        d="M 42 22 L 50 14 L 58 22"
        fill="none"
        stroke={outline}
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      {eyes}
      <path d="M 46 60 L 54 60 L 50 66 Z" fill={beak} stroke={outline} strokeWidth="1" />
      <circle cx="32" cy="58" r="3" fill="#FFB3C7" opacity="0.7" />
      <circle cx="68" cy="58" r="3" fill="#FFB3C7" opacity="0.7" />
    </svg>
  )
}

const SBMini = ({ size = 40, expr = 'plain' }) => (
  <SBBear size={size} expr={expr} palette={SB_MINI_PALETTE} />
)

// Sakura petal (decorative)
const SBPetal = ({ size = 14, color = '#F25C7A', rotate = 0, opacity = 0.7 }) => (
  <svg
    width={size}
    height={size}
    viewBox="0 0 20 20"
    style={{ transform: `rotate(${rotate}deg)`, opacity, display: 'block' }}
  >
    <path d="M10 2 C 14 4 14 10 10 12 C 6 10 6 4 10 2 Z" fill={color} />
    <circle cx="10" cy="11" r="1.2" fill="#fff" opacity="0.6" />
  </svg>
)

// Bear paw print (small footprint)
const SBPaw = ({ size = 16, color = '#FFA3B8' }) => (
  <svg width={size} height={size} viewBox="0 0 20 20" style={{ display: 'block' }}>
    <ellipse cx="10" cy="13" rx="5" ry="4" fill={color} />
    <circle cx="5" cy="7" r="2" fill={color} />
    <circle cx="10" cy="5" r="2" fill={color} />
    <circle cx="15" cy="7" r="2" fill={color} />
    <circle cx="3" cy="11" r="1.5" fill={color} />
    <circle cx="17" cy="11" r="1.5" fill={color} />
  </svg>
)

Object.assign(window, {
  SBBear,
  SBChick,
  SBMini,
  SBPetal,
  SBPaw,
  SB_BEAR_PALETTE,
  SB_MINI_PALETTE,
  SB_CHICK_PALETTE,
})

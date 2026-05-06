// Sakura Bear — Design System tokens & shared atoms.
// Colors mirror existing android theme (Sakura/Warm) but introduce mascot palette.

const SB = {
  // Sakura
  s50: '#FFF5F7',
  s100: '#FFE3EA',
  s200: '#FFC7D4',
  s300: '#FFA3B8',
  s400: '#FF7896',
  s500: '#F25C7A',
  s600: '#D94560',
  s700: '#B3354C',
  s800: '#8A283A',
  s900: '#601C29',
  // Warm neutrals
  w50: '#FAFAF9',
  w100: '#F5F5F4',
  w200: '#E7E5E4',
  w300: '#D6D3D1',
  w400: '#A8A29E',
  w500: '#78716C',
  w600: '#57534E',
  w700: '#44403C',
  w800: '#292524',
  w900: '#1C1917',
  // Sky
  sky100: '#E0F2FE',
  sky500: '#0EA5E9',
  sky700: '#0369A1',
  // Bear-specific cream tones
  cream: '#FFFBF5',
  beige: '#FBF0E0',
  beigeIn: '#F4D6B8',
  // Semantic
  success: '#10B981',
  successBg: '#D1FAE5',
  warn: '#F59E0B',
  warnBg: '#FEF3C7',
  error: '#EF4444',
}

const SB_FONT = `'Gowun Dodum', 'Quicksand', system-ui, sans-serif`

// Pill / chip
const SBChip = ({ children, bg = SB.s100, fg = SB.s700, size = 'md', icon }) => (
  <span
    style={{
      display: 'inline-flex',
      alignItems: 'center',
      gap: 4,
      padding: size === 'sm' ? '2px 8px' : '4px 10px',
      borderRadius: 999,
      background: bg,
      color: fg,
      fontSize: size === 'sm' ? 10 : 11,
      fontWeight: 700,
      whiteSpace: 'nowrap',
    }}
  >
    {icon}
    {children}
  </span>
)

// Card with optional petal corner
const SBCard = ({ children, padding = 14, radius = 18, bg = '#fff', border = SB.s100, style }) => (
  <div
    style={{
      background: bg,
      border: `1px solid ${border}`,
      borderRadius: radius,
      padding,
      boxShadow: '0 2px 0 rgba(96,28,41,0.04), 0 1px 2px rgba(96,28,41,0.04)',
      ...style,
    }}
  >
    {children}
  </div>
)

// Status pill helpers
const STATUS = {
  계획중: { bg: SB.s100, fg: SB.s700 },
  여행중: { bg: '#E0F2FE', fg: SB.sky700 },
  완료: { bg: SB.w200, fg: SB.w700 },
  보관: { bg: SB.w100, fg: SB.w500 },
}

const SBStatusPill = ({ label, size }) => {
  const { bg, fg } = STATUS[label] || STATUS['계획중']
  return (
    <SBChip bg={bg} fg={fg} size={size}>
      {label}
    </SBChip>
  )
}

// Generic icon (24x24 stroke)
const SBIcon = ({ d, size = 18, color = 'currentColor', stroke = 1.7 }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d={d} stroke={color} strokeWidth={stroke} strokeLinecap="round" strokeLinejoin="round" />
  </svg>
)

const SB_ICONS = {
  home: 'M3 11.5L12 4l9 7.5V20a1 1 0 01-1 1h-5v-6h-6v6H4a1 1 0 01-1-1z',
  trip: 'M3 7l9-4 9 4v10l-9 4-9-4z M3 7l9 4 9-4 M12 11v10',
  cam: 'M4 8h3l2-2h6l2 2h3v11H4zM12 17a3.5 3.5 0 100-7 3.5 3.5 0 000 7z',
  bear: 'M7.5 8a4.5 4.5 0 119 0M5 5.5a2.2 2.2 0 100-4.4 2.2 2.2 0 000 4.4zm14 0a2.2 2.2 0 100-4.4 2.2 2.2 0 000 4.4zM5 14a7 7 0 1014 0',
  pin: 'M12 21s-7-7.5-7-12a7 7 0 1114 0c0 4.5-7 12-7 12zm0-9a3 3 0 100-6 3 3 0 000 6z',
  cal: 'M4 6h16v14H4zM4 10h16M9 3v4M15 3v4',
  list: 'M4 6h16M4 12h16M4 18h16',
  food: 'M5 4v8a3 3 0 003 3v6m12-17l-2 5 2 1v11M9 4v5l-2 1',
  bed: 'M3 18v-6a2 2 0 012-2h6V6h6a4 4 0 014 4v8',
  note: 'M5 4h11l3 3v13H5z M9 9h6 M9 13h6 M9 17h4',
  wallet: 'M3 7h15a3 3 0 013 3v8a2 2 0 01-2 2H5a2 2 0 01-2-2zM3 7a2 2 0 012-2h11v2 M16 13h3',
  back: 'M15 5l-7 7 7 7',
  plus: 'M12 5v14M5 12h14',
  more: 'M5 12h.01M12 12h.01M19 12h.01',
  search: 'M11 19a8 8 0 100-16 8 8 0 000 16zm10 2l-4.35-4.35',
  star: 'M12 3l2.7 5.5 6 .9-4.4 4.3 1 6L12 17l-5.4 2.8 1-6-4.4-4.3 6-.9z',
  heart: 'M12 21s-7-4.6-9.5-9A5.5 5.5 0 0112 6a5.5 5.5 0 019.5 6c-2.5 4.4-9.5 9-9.5 9z',
  bell: 'M6 8a6 6 0 1112 0v5l2 3H4l2-3z M10 19a2 2 0 004 0',
  flight: 'M3 16l9-13 2 8 7 1-9 4-2 7-3-5-4-2z',
  clock: 'M12 21a9 9 0 100-18 9 9 0 000 18zM12 7v5l3 2',
  check: 'M5 12l5 5L20 7',
  edit: 'M4 20h4l10-10-4-4L4 16zM14 6l4 4',
  filter: 'M4 5h16l-6 8v6l-4-2v-4z',
  receipt: 'M5 3h14v18l-3-2-2 2-2-2-2 2-2-2-3 2zM8 8h8M8 12h8M8 16h5',
  download: 'M12 4v12m0 0l-4-4m4 4l4-4M4 20h16',
  globe: 'M12 21a9 9 0 100-18 9 9 0 000 18zM3 12h18 M12 3a14 14 0 010 18 M12 3a14 14 0 000 18',
  user: 'M12 12a4 4 0 100-8 4 4 0 000 8zM4 21a8 8 0 0116 0',
  logout: 'M9 4H5v16h4 M16 8l4 4-4 4 M20 12H10',
  mood: 'M12 21a9 9 0 100-18 9 9 0 000 18zM8 14s1.5 2 4 2 4-2 4-2 M9 9h.01 M15 9h.01',
  sun: 'M12 4v2 M12 18v2 M4 12H2 M22 12h-2 M5 5l1.5 1.5 M17.5 17.5L19 19 M5 19l1.5-1.5 M17.5 6.5L19 5 M12 17a5 5 0 100-10 5 5 0 000 10z',
}

// Bottom tab bar
const SBTabBar = ({ active = 'home', onChange = () => {} }) => {
  const items = [
    { id: 'home', l: '홈', d: SB_ICONS.home },
    { id: 'trips', l: '여행', d: SB_ICONS.trip },
    { id: 'photos', l: '사진', d: SB_ICONS.cam },
    { id: 'me', l: '나', d: SB_ICONS.bear },
  ]
  return (
    <div
      style={{
        position: 'absolute',
        left: 0,
        right: 0,
        bottom: 0,
        background: '#fff',
        borderTop: `1px solid ${SB.s100}`,
        display: 'flex',
        justifyContent: 'space-around',
        padding: '8px 4px 12px',
      }}
    >
      {items.map((it) => {
        const a = active === it.id
        return (
          <button
            key={it.id}
            onClick={() => onChange(it.id)}
            style={{
              background: 'none',
              border: 'none',
              cursor: 'pointer',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              gap: 3,
              padding: '4px 14px',
              borderRadius: 14,
              fontFamily: SB_FONT,
              color: a ? SB.s600 : SB.w500,
            }}
          >
            <div
              style={{
                padding: '5px 12px',
                borderRadius: 99,
                background: a ? SB.s100 : 'transparent',
              }}
            >
              <SBIcon d={it.d} size={20} color={a ? SB.s600 : SB.w500} stroke={a ? 2 : 1.7} />
            </div>
            <div style={{ fontSize: 10, fontWeight: 700 }}>{it.l}</div>
          </button>
        )
      })}
    </div>
  )
}

// Top app bar (with back / title / actions)
const SBAppBar = ({ title, onBack, right, transparent = false }) => (
  <div
    style={{
      height: 52,
      padding: '0 8px 0 4px',
      display: 'flex',
      alignItems: 'center',
      gap: 4,
      background: transparent ? 'transparent' : SB.s50,
      fontFamily: SB_FONT,
    }}
  >
    {onBack ? (
      <button
        onClick={onBack}
        style={{
          width: 40,
          height: 40,
          borderRadius: 99,
          border: 'none',
          background: 'transparent',
          cursor: 'pointer',
          display: 'grid',
          placeItems: 'center',
          color: SB.s600,
        }}
      >
        <SBIcon d={SB_ICONS.back} size={22} />
      </button>
    ) : (
      <div style={{ width: 12 }} />
    )}
    <div style={{ flex: 1, fontSize: 17, fontWeight: 700, color: SB.s700 }}>{title}</div>
    {right}
  </div>
)

// Generic floating action button
const SBFab = ({ icon, onClick, label }) => (
  <button
    onClick={onClick}
    style={{
      position: 'absolute',
      right: 16,
      bottom: 90,
      height: 56,
      padding: label ? '0 18px 0 14px' : 0,
      minWidth: 56,
      borderRadius: 99,
      background: SB.s500,
      color: '#fff',
      border: 'none',
      cursor: 'pointer',
      display: 'flex',
      alignItems: 'center',
      gap: 6,
      boxShadow: '0 8px 20px rgba(242,92,122,0.4), 0 2px 6px rgba(242,92,122,0.3)',
      fontFamily: SB_FONT,
      fontWeight: 700,
      fontSize: 14,
    }}
  >
    <SBIcon d={icon} size={22} color="#fff" stroke={2.2} />
    {label}
  </button>
)

// Photo placeholder
const SBPhoto = ({
  h = 120,
  w = '100%',
  label = 'photo',
  tint = SB.s100,
  stroke = SB.s200,
  radius = 14,
}) => (
  <div
    style={{
      width: w,
      height: h,
      borderRadius: radius,
      background: `repeating-linear-gradient(45deg, ${tint} 0 8px, ${stroke} 8px 9px)`,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      color: SB.s700,
      fontFamily: 'ui-monospace, Menlo, monospace',
      fontSize: 10,
      letterSpacing: 1,
      textTransform: 'uppercase',
      fontWeight: 600,
    }}
  >
    {label}
  </div>
)

// Section title (used by everything)
const SBSection = ({ icon, title, count, action, children }) => (
  <div>
    <div
      style={{
        padding: '8px 18px 8px',
        display: 'flex',
        alignItems: 'center',
        gap: 8,
      }}
    >
      {icon && (
        <div
          style={{
            width: 28,
            height: 28,
            borderRadius: 10,
            background: SB.s100,
            display: 'grid',
            placeItems: 'center',
            color: SB.s600,
          }}
        >
          <SBIcon d={icon} size={16} />
        </div>
      )}
      <div style={{ flex: 1, fontSize: 14, fontWeight: 700, color: SB.s900 }}>
        {title}
        {count != null && (
          <span style={{ fontSize: 12, color: SB.w500, marginLeft: 6, fontWeight: 600 }}>
            {count}
          </span>
        )}
      </div>
      {action}
    </div>
    {children}
  </div>
)

Object.assign(window, {
  SB,
  SB_FONT,
  SB_ICONS,
  SBChip,
  SBCard,
  SBStatusPill,
  SBIcon,
  SBTabBar,
  SBAppBar,
  SBFab,
  SBPhoto,
  SBSection,
})

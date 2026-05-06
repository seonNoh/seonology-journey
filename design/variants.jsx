// 5 mascot-bear travel-app candidates for Seonology Journey.
// Each variant: { id, name, palette, font, render: (size) => trip-list screen JSX }

// ─────────────────────────────────────────────────────────────
// Shared atoms (use unique style-object names per file)
// ─────────────────────────────────────────────────────────────
const Stamp = ({ children, bg, fg, border, rotate = 0, size = 'md' }) => (
  <div
    style={{
      display: 'inline-flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: size === 'sm' ? '3px 8px' : '4px 10px',
      borderRadius: 999,
      background: bg,
      color: fg,
      border: border ? `1.2px dashed ${border}` : 'none',
      fontSize: size === 'sm' ? 10 : 11,
      fontWeight: 600,
      transform: `rotate(${rotate}deg)`,
      whiteSpace: 'nowrap',
    }}
  >
    {children}
  </div>
)

const TabIcon = ({ d, active, color, mutedColor }) => (
  <svg width="22" height="22" viewBox="0 0 24 24" fill="none">
    <path
      d={d}
      stroke={active ? color : mutedColor}
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
)

const ICONS = {
  home: 'M3 11.5L12 4l9 7.5V20a1 1 0 01-1 1h-5v-6h-6v6H4a1 1 0 01-1-1z',
  trips: 'M5 7h14v12H5zM5 7l2-3h10l2 3M9 11h6',
  cam: 'M4 8h3l2-2h6l2 2h3v11H4zM12 17a3.5 3.5 0 100-7 3.5 3.5 0 000 7z',
  bear: 'M8 9a4 4 0 118 0M6 7a2 2 0 100-4 2 2 0 000 4zm12 0a2 2 0 100-4 2 2 0 000 4zM6 14a6 6 0 1012 0',
  pin: 'M12 21s-7-7.5-7-12a7 7 0 1114 0c0 4.5-7 12-7 12zm0-9a3 3 0 100-6 3 3 0 000 6z',
  cal: 'M4 6h16v14H4zM4 10h16M9 3v4M15 3v4',
}

// Cute placeholder photo (subtle stripes + label)
const PhotoSlot = ({
  w = '100%',
  h = 90,
  label = 'photo',
  tint = '#E8DFD3',
  stroke = '#C9BCA9',
}) => (
  <div
    style={{
      width: w,
      height: h,
      borderRadius: 14,
      background: `repeating-linear-gradient(45deg, ${tint} 0 8px, ${stroke}55 8px 9px)`,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      color: '#7a6a55',
      fontFamily: 'ui-monospace, Menlo, monospace',
      fontSize: 10,
      letterSpacing: 1,
      textTransform: 'uppercase',
    }}
  >
    {label}
  </div>
)

// ─────────────────────────────────────────────────────────────
// VARIANT 1 — Honey Beige (warm cream, sticker cards)
// ─────────────────────────────────────────────────────────────
const honeyPalette = {
  bg: '#FBF3E4',
  surface: '#FFFFFF',
  ink: '#3B2A20',
  mute: '#8A7560',
  fur: '#E9C9A0',
  ear: '#C99466',
  snout: '#F4E1C5',
  outline: '#3B2A20',
  accent: '#E89A3C',
  accentSoft: '#FCE6C2',
  divider: '#EFE3CC',
}
const honeyFont = `'Gowun Dodum', 'Quicksand', system-ui, sans-serif`

function HoneyBeigeScreen() {
  const p = honeyPalette
  return (
    <div
      style={{
        background: p.bg,
        minHeight: '100%',
        fontFamily: honeyFont,
        color: p.ink,
        paddingBottom: 80,
      }}
    >
      {/* Top */}
      <div style={{ padding: '14px 18px 8px', display: 'flex', alignItems: 'center', gap: 10 }}>
        <BearMascot size={44} palette={p} />
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 11, color: p.mute, fontWeight: 500 }}>こんにちは, さん</div>
          <div style={{ fontSize: 18, fontWeight: 700 }}>오늘도 살랑살랑 여행 ♪</div>
        </div>
        <div
          style={{
            width: 36,
            height: 36,
            borderRadius: 12,
            background: p.accentSoft,
            display: 'grid',
            placeItems: 'center',
          }}
        >
          <TabIcon d={ICONS.cal} active color={p.accent} mutedColor={p.accent} />
        </div>
      </div>

      {/* Hero card with bear+chick */}
      <div
        style={{
          margin: '8px 16px',
          padding: 16,
          borderRadius: 22,
          background: p.accentSoft,
          position: 'relative',
          overflow: 'hidden',
          border: `1.5px solid ${p.divider}`,
        }}
      >
        <div style={{ fontSize: 11, color: p.mute, fontWeight: 600, letterSpacing: 1 }}>
          NEXT TRIP
        </div>
        <div style={{ fontSize: 22, fontWeight: 700, marginTop: 4 }}>교토 벚꽃 여행</div>
        <div style={{ fontSize: 12, color: p.mute, marginTop: 2 }}>04.12 → 04.16 · D-7</div>
        <div style={{ display: 'flex', gap: 6, marginTop: 10 }}>
          <Stamp bg="#fff" fg={p.ink} border={p.outline} rotate={-3}>
            5days
          </Stamp>
          <Stamp bg={p.accent} fg="#fff">
            京都
          </Stamp>
        </div>
        <div style={{ position: 'absolute', right: -6, bottom: -6 }}>
          <BearMascot size={88} palette={p} />
        </div>
        <div style={{ position: 'absolute', right: 70, bottom: 8 }}>
          <ChickMascot size={32} />
        </div>
      </div>

      {/* Quick actions */}
      <div
        style={{
          padding: '12px 16px 4px',
          display: 'grid',
          gridTemplateColumns: 'repeat(4,1fr)',
          gap: 10,
        }}
      >
        {[
          ['일정', ICONS.cal],
          ['지도', ICONS.pin],
          ['사진', ICONS.cam],
          ['메모', ICONS.bear],
        ].map(([l, d]) => (
          <div
            key={l}
            style={{
              background: p.surface,
              borderRadius: 16,
              padding: '12px 6px',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              gap: 4,
              border: `1.5px solid ${p.divider}`,
            }}
          >
            <div
              style={{
                width: 32,
                height: 32,
                borderRadius: 10,
                background: p.accentSoft,
                display: 'grid',
                placeItems: 'center',
              }}
            >
              <TabIcon d={d} active color={p.accent} mutedColor={p.accent} />
            </div>
            <div style={{ fontSize: 11, fontWeight: 600 }}>{l}</div>
          </div>
        ))}
      </div>

      {/* Section */}
      <div style={{ padding: '14px 18px 8px', display: 'flex', alignItems: 'baseline' }}>
        <div style={{ fontSize: 14, fontWeight: 700, flex: 1 }}>다가오는 여행 ✿</div>
        <div style={{ fontSize: 11, color: p.mute }}>전체보기</div>
      </div>

      {/* Trip cards (sticker style) */}
      {[
        {
          t: '오사카 먹방투어',
          d: '05.20 → 05.24',
          tag: '계획중',
          tagBg: '#FCE6C2',
          tagFg: '#A55A0E',
        },
        {
          t: '제주 한달살이',
          d: '06.01 → 06.30',
          tag: '여행중',
          tagBg: '#D6E8C8',
          tagFg: '#3F6B1E',
        },
      ].map((x, i) => (
        <div
          key={i}
          style={{
            margin: '0 16px 10px',
            padding: 12,
            borderRadius: 18,
            background: p.surface,
            border: `1.5px solid ${p.divider}`,
            display: 'flex',
            gap: 12,
            alignItems: 'center',
            boxShadow: '0 2px 0 rgba(59,42,32,0.06)',
          }}
        >
          <PhotoSlot
            w={64}
            h={64}
            label={i === 0 ? 'osaka' : 'jeju'}
            tint="#F1E4CC"
            stroke="#C9BCA9"
          />
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
              <div style={{ fontSize: 14, fontWeight: 700, flex: 1 }}>{x.t}</div>
              <Stamp bg={x.tagBg} fg={x.tagFg} size="sm">
                {x.tag}
              </Stamp>
            </div>
            <div style={{ fontSize: 11, color: p.mute, marginTop: 2 }}>{x.d}</div>
            <div style={{ display: 'flex', gap: 4, marginTop: 8 }}>
              <div style={{ width: 18, height: 18, borderRadius: 6, background: p.accentSoft }} />
              <div style={{ width: 18, height: 18, borderRadius: 6, background: p.divider }} />
              <div
                style={{
                  width: 18,
                  height: 18,
                  borderRadius: 6,
                  background: p.accent,
                  opacity: 0.4,
                }}
              />
            </div>
          </div>
        </div>
      ))}

      {/* Bottom nav */}
      <BottomNav
        palette={{
          bg: '#fff',
          activeBg: p.accentSoft,
          activeFg: p.accent,
          mute: p.mute,
          border: p.divider,
        }}
      />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// VARIANT 2 — Milk & Cocoa (minimal, monochrome browns)
// ─────────────────────────────────────────────────────────────
const milkPalette = {
  bg: '#FAF6F1',
  surface: '#FFFFFF',
  ink: '#2C1F17',
  mute: '#9A8678',
  fur: '#D9B795',
  ear: '#9F6B47',
  snout: '#EFD9BE',
  outline: '#2C1F17',
  accent: '#6B4427',
  accentSoft: '#EFE2D2',
  divider: '#EFE6DA',
}

function MilkCocoaScreen() {
  const p = milkPalette
  return (
    <div
      style={{
        background: p.bg,
        minHeight: '100%',
        fontFamily: honeyFont,
        color: p.ink,
        paddingBottom: 80,
      }}
    >
      {/* Top — minimal */}
      <div style={{ padding: '20px 20px 10px' }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <div style={{ fontSize: 11, color: p.mute, letterSpacing: 2, fontWeight: 600 }}>
              SEONOLOGY
            </div>
            <div style={{ fontSize: 26, fontWeight: 700, marginTop: 4, lineHeight: 1.1 }}>
              Journey
            </div>
          </div>
          <BearMascot size={48} palette={p} />
        </div>
      </div>

      {/* Big featured trip */}
      <div
        style={{
          margin: '8px 16px 16px',
          borderRadius: 20,
          overflow: 'hidden',
          background: p.surface,
          border: `1.5px solid ${p.divider}`,
        }}
      >
        <PhotoSlot h={140} label="kyoto · sakura" tint="#EFE2D2" stroke="#C7B49B" />
        <div style={{ padding: 14 }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
            <div style={{ width: 6, height: 6, borderRadius: 99, background: p.accent }} />
            <div style={{ fontSize: 10, color: p.mute, letterSpacing: 1.5, fontWeight: 600 }}>
              NEXT · D-7
            </div>
          </div>
          <div style={{ fontSize: 18, fontWeight: 700, marginTop: 4 }}>교토 벚꽃 여행</div>
          <div style={{ fontSize: 12, color: p.mute, marginTop: 2 }}>04.12 — 04.16 · 5박 6일</div>
          <div style={{ display: 'flex', gap: 8, marginTop: 12, alignItems: 'center' }}>
            <div
              style={{
                flex: 1,
                height: 4,
                borderRadius: 99,
                background: p.divider,
                overflow: 'hidden',
              }}
            >
              <div style={{ width: '35%', height: '100%', background: p.accent }} />
            </div>
            <div style={{ fontSize: 11, color: p.mute, fontWeight: 600 }}>준비 35%</div>
          </div>
        </div>
      </div>

      {/* List header */}
      <div style={{ padding: '4px 20px 8px', display: 'flex', alignItems: 'center' }}>
        <div style={{ fontSize: 13, fontWeight: 700, flex: 1 }}>모든 여행</div>
        <div style={{ fontSize: 11, color: p.mute }}>↓ 최근순</div>
      </div>

      {/* List items — minimal rows */}
      {[
        { t: '오사카 먹방투어', d: '05.20 → 05.24', tag: '계획중' },
        { t: '제주 한달살이', d: '06.01 → 06.30', tag: '여행중' },
        { t: '도쿄 윈터 쇼핑', d: '12.20 → 12.27', tag: '완료' },
      ].map((x, i) => (
        <div
          key={i}
          style={{
            margin: '0 16px 8px',
            padding: '12px 14px',
            borderRadius: 14,
            background: p.surface,
            border: `1px solid ${p.divider}`,
            display: 'flex',
            alignItems: 'center',
            gap: 12,
          }}
        >
          <div
            style={{
              width: 36,
              height: 36,
              borderRadius: 10,
              background: p.accentSoft,
              display: 'grid',
              placeItems: 'center',
              fontSize: 14,
              fontWeight: 700,
              color: p.accent,
            }}
          >
            {String(i + 1).padStart(2, '0')}
          </div>
          <div style={{ flex: 1 }}>
            <div style={{ fontSize: 13, fontWeight: 700 }}>{x.t}</div>
            <div style={{ fontSize: 11, color: p.mute, marginTop: 2 }}>{x.d}</div>
          </div>
          <div style={{ fontSize: 10, color: p.mute, fontWeight: 600 }}>{x.tag}</div>
        </div>
      ))}

      <BottomNav
        palette={{
          bg: '#fff',
          activeBg: p.accentSoft,
          activeFg: p.accent,
          mute: p.mute,
          border: p.divider,
        }}
      />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// VARIANT 3 — Cafe Lounge (paper texture, stamps, journal)
// ─────────────────────────────────────────────────────────────
const cafePalette = {
  bg: '#F2E8D5',
  surface: '#FFFBF1',
  ink: '#3B2A20',
  mute: '#8A7560',
  fur: '#D4A574',
  ear: '#A06B3F',
  snout: '#F0DBB8',
  outline: '#3B2A20',
  accent: '#8B3A2B',
  accentSoft: '#E8D4C4',
  divider: '#D9C9AF',
}

function CafeLoungeScreen() {
  const p = cafePalette
  return (
    <div
      style={{
        background: p.bg,
        minHeight: '100%',
        fontFamily: honeyFont,
        color: p.ink,
        paddingBottom: 80,
        backgroundImage: `radial-gradient(${p.divider}66 1px, transparent 1px)`,
        backgroundSize: '14px 14px',
      }}
    >
      {/* Header — journal title */}
      <div style={{ padding: '16px 18px 6px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <BearMascot size={36} palette={p} />
          <div style={{ flex: 1 }}>
            <div style={{ fontSize: 11, color: p.mute, fontStyle: 'italic' }}>
              my travel journal
            </div>
            <div style={{ fontSize: 18, fontWeight: 700 }}>くまの旅ノート</div>
          </div>
          <Stamp bg="transparent" fg={p.accent} border={p.accent} rotate={6}>
            EST. 2026
          </Stamp>
        </div>
        <div style={{ height: 1, background: p.divider, marginTop: 10 }} />
      </div>

      {/* Polaroid hero */}
      <div
        style={{
          margin: '10px 16px',
          padding: 10,
          borderRadius: 6,
          background: p.surface,
          transform: 'rotate(-1.5deg)',
          boxShadow: '0 4px 0 rgba(59,42,32,0.08), 0 8px 20px rgba(59,42,32,0.06)',
        }}
      >
        <PhotoSlot h={130} label="next · kyoto" tint={p.accentSoft} stroke={p.divider} />
        <div style={{ padding: '10px 4px 4px', display: 'flex', alignItems: 'center', gap: 8 }}>
          <div style={{ flex: 1 }}>
            <div style={{ fontSize: 14, fontWeight: 700 }}>교토 벚꽃 여행</div>
            <div style={{ fontSize: 11, color: p.mute, marginTop: 2 }}>04.12 — 04.16</div>
          </div>
          <Stamp bg={p.accent} fg="#fff" rotate={-4}>
            D-7
          </Stamp>
        </div>
      </div>

      {/* Stamp row */}
      <div style={{ padding: '6px 18px 10px', display: 'flex', gap: 8, flexWrap: 'wrap' }}>
        <Stamp bg="transparent" fg={p.accent} border={p.accent} rotate={-3}>
          ✓ 항공권
        </Stamp>
        <Stamp bg="transparent" fg={p.mute} border={p.mute} rotate={2}>
          ○ 숙소
        </Stamp>
        <Stamp bg="transparent" fg={p.mute} border={p.mute} rotate={-1}>
          ○ 환전
        </Stamp>
        <Stamp bg={p.accent} fg="#fff" rotate={4}>
          3 / 8
        </Stamp>
      </div>

      {/* Notebook lines section */}
      <div
        style={{
          margin: '4px 16px 0',
          padding: '14px 16px',
          borderRadius: 12,
          background: p.surface,
          border: `1.5px solid ${p.divider}`,
          backgroundImage: `repeating-linear-gradient(transparent 0 27px, ${p.divider}88 27px 28px)`,
        }}
      >
        <div style={{ fontSize: 13, fontWeight: 700, marginBottom: 8 }}>지난 여행</div>
        {[
          ['도쿄 윈터 쇼핑', '12.27', '🌟🌟🌟🌟'],
          ['부산 바다 산책', '10.04', '🌟🌟🌟'],
          ['속초 캠핑', '08.15', '🌟🌟🌟🌟🌟'],
        ].map(([t, d, s]) => (
          <div key={t} style={{ display: 'flex', alignItems: 'center', height: 28, fontSize: 12 }}>
            <div style={{ flex: 1, fontWeight: 600 }}>{t}</div>
            <div style={{ color: p.mute, marginRight: 8 }}>{d}</div>
            <div style={{ fontSize: 10 }}>{s}</div>
          </div>
        ))}
      </div>

      <BottomNav
        palette={{
          bg: p.surface,
          activeBg: p.accentSoft,
          activeFg: p.accent,
          mute: p.mute,
          border: p.divider,
        }}
      />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// VARIANT 4 — Sakura Bear (keeps existing pink palette + bear)
// ─────────────────────────────────────────────────────────────
const sakuraPalette = {
  bg: '#FFF5F7',
  surface: '#FFFFFF',
  ink: '#601C29',
  mute: '#A87680',
  fur: '#F4D6B8',
  ear: '#D9A179',
  snout: '#FBEAD5',
  outline: '#601C29',
  accent: '#F25C7A',
  accentSoft: '#FFE3EA',
  divider: '#FFD1DD',
}

function SakuraBearScreen() {
  const p = sakuraPalette
  return (
    <div
      style={{
        background: `linear-gradient(180deg, ${p.bg} 0%, #FAFAF9 60%)`,
        minHeight: '100%',
        fontFamily: honeyFont,
        color: p.ink,
        paddingBottom: 80,
      }}
    >
      {/* Petal header */}
      <div style={{ padding: '14px 18px 8px', position: 'relative' }}>
        {/* petals */}
        {[
          [8, 14, -10],
          [320, 30, 18],
          [60, 70, 40],
          [280, 90, -5],
        ].map(([x, y, r], i) => (
          <svg
            key={i}
            width="14"
            height="14"
            viewBox="0 0 20 20"
            style={{
              position: 'absolute',
              left: x,
              top: y,
              transform: `rotate(${r}deg)`,
              opacity: 0.6,
            }}
          >
            <path d="M10 2 C 14 4 14 10 10 12 C 6 10 6 4 10 2 Z" fill={p.accent} />
          </svg>
        ))}
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <BearMascot size={44} palette={p} />
          <div style={{ flex: 1 }}>
            <div style={{ fontSize: 11, color: p.mute, fontWeight: 600 }}>HELLO ♡</div>
            <div style={{ fontSize: 18, fontWeight: 700 }}>오늘의 여행 기록</div>
          </div>
          <ChickMascot size={32} />
        </div>
      </div>

      {/* Hero — pink soft card */}
      <div
        style={{
          margin: '6px 16px',
          padding: 16,
          borderRadius: 24,
          background: `linear-gradient(135deg, ${p.accentSoft} 0%, #FFF 100%)`,
          border: `1.5px solid ${p.divider}`,
          position: 'relative',
          overflow: 'hidden',
        }}
      >
        <div style={{ fontSize: 11, color: p.accent, fontWeight: 700, letterSpacing: 1 }}>
          NEXT · D-7
        </div>
        <div style={{ fontSize: 22, fontWeight: 700, marginTop: 4 }}>교토 벚꽃 여행</div>
        <div style={{ fontSize: 12, color: p.mute, marginTop: 2 }}>04.12 → 04.16</div>
        <div style={{ display: 'flex', gap: 6, marginTop: 12 }}>
          <div
            style={{
              background: '#fff',
              padding: '6px 10px',
              borderRadius: 99,
              fontSize: 10,
              fontWeight: 700,
              color: p.accent,
              border: `1px solid ${p.divider}`,
            }}
          >
            ✿ 5박 6일
          </div>
          <div
            style={{
              background: p.accent,
              padding: '6px 10px',
              borderRadius: 99,
              fontSize: 10,
              fontWeight: 700,
              color: '#fff',
            }}
          >
            京都
          </div>
        </div>
        <div style={{ position: 'absolute', right: -4, bottom: -8 }}>
          <SleepingBear size={86} palette={p} />
        </div>
      </div>

      {/* Quick actions — pill row */}
      <div style={{ padding: '12px 16px', display: 'flex', gap: 8, overflow: 'hidden' }}>
        {['일정', '지도', '사진', '메모', '지출'].map((l, i) => (
          <div
            key={l}
            style={{
              padding: '8px 14px',
              borderRadius: 99,
              background: i === 0 ? p.accent : '#fff',
              color: i === 0 ? '#fff' : p.ink,
              fontSize: 11,
              fontWeight: 700,
              border: `1px solid ${i === 0 ? p.accent : p.divider}`,
              whiteSpace: 'nowrap',
            }}
          >
            {l}
          </div>
        ))}
      </div>

      {/* Section title */}
      <div style={{ padding: '6px 18px 6px', fontSize: 13, fontWeight: 700 }}>다가오는 여행</div>

      {/* Trip cards */}
      {[
        {
          t: '오사카 먹방투어',
          d: '05.20 → 05.24',
          tag: '계획중',
          tagBg: p.accentSoft,
          tagFg: p.accent,
        },
        {
          t: '제주 한달살이',
          d: '06.01 → 06.30',
          tag: '여행중',
          tagBg: '#E0F2FE',
          tagFg: '#0369A1',
        },
      ].map((x, i) => (
        <div
          key={i}
          style={{
            margin: '0 16px 10px',
            padding: 14,
            borderRadius: 18,
            background: p.surface,
            border: `1px solid ${p.divider}`,
            display: 'flex',
            gap: 12,
            alignItems: 'center',
          }}
        >
          <div
            style={{
              width: 50,
              height: 50,
              borderRadius: 14,
              background: p.accentSoft,
              display: 'grid',
              placeItems: 'center',
            }}
          >
            <BrownBearMascot size={36} />
          </div>
          <div style={{ flex: 1 }}>
            <div style={{ fontSize: 14, fontWeight: 700 }}>{x.t}</div>
            <div style={{ fontSize: 11, color: p.mute, marginTop: 2 }}>{x.d}</div>
          </div>
          <Stamp bg={x.tagBg} fg={x.tagFg} size="sm">
            {x.tag}
          </Stamp>
        </div>
      ))}

      <BottomNav
        palette={{
          bg: '#fff',
          activeBg: p.accentSoft,
          activeFg: p.accent,
          mute: p.mute,
          border: p.divider,
        }}
      />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// VARIANT 5 — Forest Picnic (sage green + mustard)
// ─────────────────────────────────────────────────────────────
const forestPalette = {
  bg: '#F2EFE4',
  surface: '#FBF8EE',
  ink: '#2E3924',
  mute: '#7A8568',
  fur: '#D4B896',
  ear: '#8F6B49',
  snout: '#EDDFC3',
  outline: '#2E3924',
  accent: '#5C7A3C',
  accentSoft: '#DCE5C8',
  divider: '#D9D2BC',
  mustard: '#D4A23B',
}

function ForestPicnicScreen() {
  const p = forestPalette
  return (
    <div
      style={{
        background: p.bg,
        minHeight: '100%',
        fontFamily: honeyFont,
        color: p.ink,
        paddingBottom: 80,
      }}
    >
      {/* Header */}
      <div style={{ padding: '16px 18px 8px', display: 'flex', alignItems: 'center', gap: 10 }}>
        <div
          style={{
            width: 44,
            height: 44,
            borderRadius: 14,
            background: p.accentSoft,
            display: 'grid',
            placeItems: 'center',
          }}
        >
          <BearMascot size={36} palette={p} />
        </div>
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 11, color: p.mute, fontWeight: 600, letterSpacing: 1 }}>
            FIELD NOTES
          </div>
          <div style={{ fontSize: 18, fontWeight: 700 }}>오늘의 모험 ⛺</div>
        </div>
        <div
          style={{
            width: 36,
            height: 36,
            borderRadius: 12,
            background: p.mustard,
            display: 'grid',
            placeItems: 'center',
            color: '#fff',
          }}
        >
          <TabIcon d={ICONS.cam} active color="#fff" mutedColor="#fff" />
        </div>
      </div>

      {/* Hero — banner with green + leaf badge */}
      <div
        style={{
          margin: '6px 16px',
          borderRadius: 20,
          overflow: 'hidden',
          background: p.accent,
          color: '#fff',
          position: 'relative',
        }}
      >
        <div style={{ padding: '16px 16px 18px' }}>
          <div style={{ fontSize: 11, fontWeight: 700, letterSpacing: 1.5, opacity: 0.9 }}>
            UPCOMING · D-7
          </div>
          <div style={{ fontSize: 22, fontWeight: 700, marginTop: 4 }}>지리산 가을캠핑</div>
          <div style={{ fontSize: 12, opacity: 0.85, marginTop: 2 }}>10.18 → 10.20 · 2박 3일</div>
          <div style={{ display: 'flex', gap: 6, marginTop: 12 }}>
            <div
              style={{
                background: 'rgba(255,255,255,0.18)',
                padding: '5px 10px',
                borderRadius: 99,
                fontSize: 10,
                fontWeight: 700,
                backdropFilter: 'blur(4px)',
              }}
            >
              ⛺ 2명
            </div>
            <div
              style={{
                background: p.mustard,
                padding: '5px 10px',
                borderRadius: 99,
                fontSize: 10,
                fontWeight: 700,
                color: '#3B2A20',
              }}
            >
              HONEY POT
            </div>
          </div>
        </div>
        <div style={{ position: 'absolute', right: 8, bottom: 0 }}>
          <BearMascot
            size={80}
            palette={{ ...p, fur: p.fur, ear: p.ear, snout: p.snout, outline: '#1f2618' }}
          />
        </div>
        {/* leaves */}
        <svg
          width="90"
          height="60"
          viewBox="0 0 90 60"
          style={{ position: 'absolute', left: -10, top: -8, opacity: 0.25 }}
        >
          <path d="M20 30 Q 30 5 60 18 Q 50 40 20 30 Z" fill="#fff" />
        </svg>
      </div>

      {/* Two col cards */}
      <div
        style={{
          padding: '14px 16px 6px',
          display: 'grid',
          gridTemplateColumns: '1fr 1fr',
          gap: 10,
        }}
      >
        {[
          { t: '오사카 먹방', d: '05.20', tag: '계획', big: false },
          { t: '제주 한달', d: '06.01', tag: '여행중', big: true },
        ].map((x, i) => (
          <div
            key={i}
            style={{
              background: p.surface,
              borderRadius: 16,
              padding: 12,
              border: `1.5px solid ${p.divider}`,
            }}
          >
            <PhotoSlot h={70} label={x.t.split(' ')[0]} tint={p.accentSoft} stroke={p.divider} />
            <div style={{ marginTop: 8, fontSize: 13, fontWeight: 700 }}>{x.t}</div>
            <div style={{ fontSize: 10, color: p.mute, marginTop: 2 }}>{x.d}</div>
            <div
              style={{
                marginTop: 6,
                display: 'inline-block',
                padding: '3px 8px',
                borderRadius: 99,
                background: x.big ? p.accent : p.accentSoft,
                color: x.big ? '#fff' : p.accent,
                fontSize: 9,
                fontWeight: 700,
              }}
            >
              {x.tag}
            </div>
          </div>
        ))}
      </div>

      {/* Stat row */}
      <div
        style={{
          margin: '8px 16px',
          padding: '12px 14px',
          borderRadius: 16,
          background: p.surface,
          border: `1.5px solid ${p.divider}`,
          display: 'flex',
          alignItems: 'center',
          gap: 12,
        }}
      >
        <ChickMascot size={36} />
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 12, fontWeight: 700 }}>이번 달 모험 4번 ✿</div>
          <div style={{ fontSize: 10, color: p.mute, marginTop: 2 }}>지난 달보다 +2회</div>
        </div>
        <div style={{ fontSize: 22, fontWeight: 700, color: p.accent }}>4</div>
      </div>

      <BottomNav
        palette={{
          bg: p.surface,
          activeBg: p.accentSoft,
          activeFg: p.accent,
          mute: p.mute,
          border: p.divider,
        }}
      />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// Shared bottom nav
// ─────────────────────────────────────────────────────────────
function BottomNav({ palette }) {
  const items = [
    { d: ICONS.home, l: '홈', a: true },
    { d: ICONS.trips, l: '여행' },
    { d: ICONS.cam, l: '사진' },
    { d: ICONS.bear, l: '나' },
  ]
  return (
    <div
      style={{
        position: 'absolute',
        left: 0,
        right: 0,
        bottom: 0,
        background: palette.bg,
        borderTop: `1px solid ${palette.border}`,
        display: 'flex',
        justifyContent: 'space-around',
        padding: '8px 4px 10px',
      }}
    >
      {items.map((it, i) => (
        <div
          key={i}
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: 3,
            padding: '4px 10px',
            borderRadius: 12,
            background: it.a ? palette.activeBg : 'transparent',
          }}
        >
          <TabIcon d={it.d} active={it.a} color={palette.activeFg} mutedColor={palette.mute} />
          <div
            style={{ fontSize: 9, fontWeight: 700, color: it.a ? palette.activeFg : palette.mute }}
          >
            {it.l}
          </div>
        </div>
      ))}
    </div>
  )
}

const VARIANTS = [
  {
    id: 'honey',
    name: 'Honey Beige',
    subtitle: '따뜻한 크림 + 스티커 카드',
    render: HoneyBeigeScreen,
  },
  {
    id: 'milk',
    name: 'Milk & Cocoa',
    subtitle: '미니멀, 우유빛 + 코코아',
    render: MilkCocoaScreen,
  },
  {
    id: 'cafe',
    name: 'Cafe Lounge',
    subtitle: '저널/폴라로이드/스탬프 무드',
    render: CafeLoungeScreen,
  },
  {
    id: 'sakura',
    name: 'Sakura Bear',
    subtitle: '기존 핑크 유지 + 곰 마스코트',
    render: SakuraBearScreen,
  },
  {
    id: 'forest',
    name: 'Forest Picnic',
    subtitle: '세이지 + 머스타드, 야외 컨셉',
    render: ForestPicnicScreen,
  },
]

Object.assign(window, {
  HoneyBeigeScreen,
  MilkCocoaScreen,
  CafeLoungeScreen,
  SakuraBearScreen,
  ForestPicnicScreen,
  VARIANTS,
  BottomNav,
  Stamp,
  TabIcon,
  PhotoSlot,
  ICONS,
  honeyPalette,
  milkPalette,
  cafePalette,
  sakuraPalette,
  forestPalette,
})

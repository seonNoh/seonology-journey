// Sakura Bear — Six core screens.
// Each screen is a self-contained component rendering inside a 412x892 phone area.
// Mock data lives at top of file.

const SB_TRIPS = [
  {
    id: 't1',
    title: '교토 벚꽃 여행',
    destination: '京都, 일본',
    start: '2026.04.12',
    end: '2026.04.16',
    days: 5,
    status: '계획중',
    dDay: -7,
    prep: 0.35,
    photo: 'kyoto · sakura',
    cover: SB.s100,
  },
  {
    id: 't2',
    title: '오사카 먹방투어',
    destination: '大阪, 일본',
    start: '2026.05.20',
    end: '2026.05.24',
    days: 5,
    status: '계획중',
    dDay: -41,
    prep: 0.1,
    photo: 'osaka · food',
  },
  {
    id: 't3',
    title: '제주 한달살이',
    destination: '제주도',
    start: '2025.06.01',
    end: '2025.06.30',
    days: 30,
    status: '여행중',
    dDay: 0,
    prep: 0.85,
    photo: 'jeju · sea',
  },
  {
    id: 't4',
    title: '도쿄 윈터 쇼핑',
    destination: '東京, 일본',
    start: '2025.12.20',
    end: '2025.12.27',
    days: 8,
    status: '완료',
    dDay: 50,
    prep: 1,
    photo: 'tokyo · winter',
  },
]

const SB_DAYS = [
  {
    id: 'd1',
    n: 1,
    date: '04.12 (일)',
    region: '오사카 → 교토',
    sched: 4,
    meals: 3,
    hasAcc: true,
    summary: '간사이공항 도착, 교토 호텔 체크인',
  },
  {
    id: 'd2',
    n: 2,
    date: '04.13 (월)',
    region: '교토 동부',
    sched: 5,
    meals: 3,
    hasAcc: true,
    summary: '기요미즈데라, 기온 산책, 야사카 신사',
  },
  {
    id: 'd3',
    n: 3,
    date: '04.14 (화)',
    region: '아라시야마',
    sched: 6,
    meals: 3,
    hasAcc: true,
    summary: '대나무숲, 도게츠교, 텐류지',
  },
  {
    id: 'd4',
    n: 4,
    date: '04.15 (수)',
    region: '교토 북부',
    sched: 4,
    meals: 3,
    hasAcc: true,
    summary: '금각사, 료안지, 니조성',
  },
  {
    id: 'd5',
    n: 5,
    date: '04.16 (목)',
    region: '교토 → 오사카',
    sched: 2,
    meals: 2,
    hasAcc: false,
    summary: '후시미이나리, 공항 이동',
  },
]

const SB_SCHEDULES = [
  {
    id: 's1',
    start: '09:30',
    end: '11:00',
    title: '기요미즈데라 (清水寺)',
    place: '교토시 히가시야마구',
    notes: '벚꽃 시즌엔 입장권 예매 추천. 야경도 예쁨',
  },
  { id: 's2', start: '11:30', end: '13:00', title: '산넨자카 거리 산책', place: '清水坂 → 二寧坂' },
  {
    id: 's3',
    start: '14:00',
    end: '15:30',
    title: '기온 거리 + 야사카 신사',
    place: '東山区 祇園町',
  },
  { id: 's4', start: '16:00', end: '17:30', title: '카페 % Arabica 키요미즈점', place: 'コーヒー' },
]

const SB_MEALS = [
  { id: 'm1', type: '아침', name: '※호텔 조식', cost: '0 JPY' },
  {
    id: 'm2',
    type: '점심',
    name: '오멘 우동 (おめん)',
    cost: '1,800 JPY',
    review: '새우튀김 우동 진리. 줄 길어서 11시 도착 추천',
  },
  { id: 'm3', type: '저녁', name: '교토 가이세키 정식', cost: '8,500 JPY' },
]

const SB_ACC = {
  name: '교토 그래나리 호텔 ★4',
  address: '교토시 시모교구 카라스마도리 가락쿠지초',
  checkin: '15:00',
  checkout: '11:00',
  cost: '24,000 JPY / 박',
}

const SB_NOTES = [
  {
    id: 'n1',
    mood: '설렘',
    date: '04.12',
    content: '드디어 출국! 칸사이공항 도착하니 벚꽃이 살랑살랑.',
  },
  {
    id: 'n2',
    mood: '맛있음',
    date: '04.12',
    content: '저녁에 먹은 가이세키 코스가 진짜 예술... 다음에 또 와야지.',
  },
  {
    id: 'n3',
    mood: '평온',
    date: '04.13',
    content: '기요미즈데라 야경 보면서 멍하니 30분. 벚꽃 잎이 바람에 우수수.',
  },
  {
    id: 'n4',
    mood: '피곤',
    date: '04.13',
    content: '하루 22,000보. 발 아파서 호텔 욕조에 1시간 누워있음 ㅋㅋ',
  },
]

const SB_EXPENSES = [
  { id: 'e1', cat: '교통', desc: '간사이공항 → 교토 (하루카)', amount: '3,640 JPY', date: '04.12' },
  { id: 'e2', cat: '숙박', desc: '교토 그래나리 호텔 (4박)', amount: '96,000 JPY', date: '04.12' },
  { id: 'e3', cat: '식사', desc: '오멘 우동 점심', amount: '1,800 JPY', date: '04.12' },
  { id: 'e4', cat: '체험', desc: '기요미즈데라 입장권', amount: '400 JPY', date: '04.13' },
  { id: 'e5', cat: '쇼핑', desc: '기온 기념품 (부채)', amount: '2,200 JPY', date: '04.13' },
  { id: 'e6', cat: '식사', desc: '가이세키 코스 저녁', amount: '8,500 JPY', date: '04.12' },
]

const CAT_META = {
  교통: { color: SB.sky500, bg: SB.sky100, icon: SB_ICONS.flight },
  식사: { color: '#D97706', bg: '#FED7AA', icon: SB_ICONS.food },
  숙박: { color: SB.s600, bg: SB.s100, icon: SB_ICONS.bed },
  체험: { color: '#7C3AED', bg: '#EDE9FE', icon: SB_ICONS.star },
  쇼핑: { color: '#059669', bg: '#D1FAE5', icon: SB_ICONS.wallet },
  기타: { color: SB.w500, bg: SB.w100, icon: SB_ICONS.note },
}

// ─────────────────────────────────────────────────────────────
// 1) Trip List (Home)
// ─────────────────────────────────────────────────────────────
function ScreenTripList({ onOpenTrip = () => {} }) {
  const next = SB_TRIPS[0]
  return (
    <div
      style={{
        background: `linear-gradient(180deg, ${SB.s50} 0%, ${SB.cream} 220px, ${SB.w50} 100%)`,
        minHeight: '100%',
        fontFamily: SB_FONT,
        color: SB.w900,
        paddingBottom: 92,
        position: 'relative',
      }}
    >
      {/* petals decoration */}
      <div
        style={{
          position: 'absolute',
          inset: '0 0 auto 0',
          height: 240,
          pointerEvents: 'none',
          overflow: 'hidden',
        }}
      >
        <div style={{ position: 'absolute', left: 16, top: 24 }}>
          <SBPetal size={12} rotate={-10} />
        </div>
        <div style={{ position: 'absolute', left: 320, top: 16, opacity: 0.7 }}>
          <SBPetal size={10} color={SB.s400} rotate={20} />
        </div>
        <div style={{ position: 'absolute', left: 60, top: 84 }}>
          <SBPetal size={14} color={SB.s300} rotate={45} />
        </div>
        <div style={{ position: 'absolute', left: 280, top: 120, opacity: 0.6 }}>
          <SBPetal size={10} rotate={-30} />
        </div>
        <div style={{ position: 'absolute', left: 200, top: 60, opacity: 0.5 }}>
          <SBPetal size={8} color={SB.s400} rotate={70} />
        </div>
      </div>

      {/* header */}
      <div
        style={{
          padding: '14px 18px 8px',
          display: 'flex',
          alignItems: 'center',
          gap: 10,
          position: 'relative',
        }}
      >
        <SBBear size={48} expr="happy" accessory="flower" />
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 11, color: SB.s600, fontWeight: 700, letterSpacing: 1 }}>
            HELLO ♡ ようこそ
          </div>
          <div style={{ fontSize: 18, fontWeight: 700, color: SB.s900 }}>
            오늘도 살랑살랑 여행 ♪
          </div>
        </div>
        <button
          style={{
            width: 40,
            height: 40,
            borderRadius: 12,
            background: '#fff',
            border: `1px solid ${SB.s100}`,
            display: 'grid',
            placeItems: 'center',
            cursor: 'pointer',
          }}
        >
          <SBIcon d={SB_ICONS.bell} size={18} color={SB.s600} />
        </button>
      </div>

      {/* hero / next trip */}
      <div
        onClick={() => onOpenTrip(next.id)}
        style={{
          margin: '4px 16px 12px',
          borderRadius: 24,
          padding: 16,
          cursor: 'pointer',
          background: `linear-gradient(135deg, ${SB.s100} 0%, #fff 65%)`,
          border: `1.5px solid ${SB.s200}`,
          position: 'relative',
          overflow: 'hidden',
          boxShadow: '0 6px 20px rgba(242,92,122,0.10)',
        }}
      >
        <div style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
          <div style={{ width: 6, height: 6, borderRadius: 99, background: SB.s500 }} />
          <div style={{ fontSize: 11, color: SB.s600, fontWeight: 700, letterSpacing: 1 }}>
            NEXT TRIP · D{next.dDay}
          </div>
        </div>
        <div style={{ fontSize: 22, fontWeight: 700, marginTop: 4 }}>{next.title}</div>
        <div style={{ fontSize: 12, color: SB.w500, marginTop: 2 }}>
          {next.start} → {next.end} · {next.days}박 6일
        </div>
        <div style={{ display: 'flex', gap: 6, marginTop: 12, alignItems: 'center' }}>
          <SBChip icon={<SBIcon d={SB_ICONS.pin} size={10} />}>{next.destination}</SBChip>
          <SBChip bg={SB.s500} fg="#fff">
            ★ HONEY POT
          </SBChip>
        </div>
        <div
          style={{ display: 'flex', gap: 8, marginTop: 12, alignItems: 'center', maxWidth: 230 }}
        >
          <div
            style={{
              flex: 1,
              height: 5,
              borderRadius: 99,
              background: SB.s100,
              overflow: 'hidden',
            }}
          >
            <div
              style={{
                width: `${next.prep * 100}%`,
                height: '100%',
                background: `linear-gradient(90deg, ${SB.s400}, ${SB.s500})`,
              }}
            />
          </div>
          <div style={{ fontSize: 11, color: SB.s700, fontWeight: 700 }}>
            준비 {Math.round(next.prep * 100)}%
          </div>
        </div>
        <div style={{ position: 'absolute', right: -2, bottom: -6 }}>
          <SBBear size={92} expr="happy" accessory="camera" />
        </div>
      </div>

      {/* quick actions */}
      <div style={{ padding: '4px 16px 8px', display: 'flex', gap: 8, overflowX: 'auto' }}>
        {[
          { l: '+ 새 여행', d: SB_ICONS.plus, primary: true },
          { l: '일정', d: SB_ICONS.cal },
          { l: '지도', d: SB_ICONS.pin },
          { l: '사진', d: SB_ICONS.cam },
          { l: '메모', d: SB_ICONS.note },
        ].map((a, i) => (
          <button
            key={i}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 6,
              padding: '8px 14px',
              borderRadius: 99,
              background: a.primary ? SB.s500 : '#fff',
              color: a.primary ? '#fff' : SB.s700,
              border: a.primary ? 'none' : `1px solid ${SB.s100}`,
              fontSize: 12,
              fontWeight: 700,
              fontFamily: SB_FONT,
              cursor: 'pointer',
              whiteSpace: 'nowrap',
              flexShrink: 0,
            }}
          >
            <SBIcon d={a.d} size={14} color={a.primary ? '#fff' : SB.s600} />
            {a.l}
          </button>
        ))}
      </div>

      {/* trip list */}
      <SBSection
        title="다가오는 여행"
        count={`(${SB_TRIPS.filter((t) => t.status !== '완료').length})`}
        action={
          <button
            style={{
              background: 'none',
              border: 'none',
              fontSize: 11,
              color: SB.w500,
              fontFamily: SB_FONT,
            }}
          >
            전체보기 ›
          </button>
        }
      ></SBSection>

      <div style={{ padding: '0 16px', display: 'flex', flexDirection: 'column', gap: 10 }}>
        {SB_TRIPS.slice(1).map((t, i) => (
          <div
            key={t.id}
            onClick={() => onOpenTrip(t.id)}
            style={{
              padding: 14,
              borderRadius: 18,
              background: '#fff',
              cursor: 'pointer',
              border: `1px solid ${SB.s100}`,
              boxShadow: '0 1px 2px rgba(96,28,41,0.04)',
              display: 'flex',
              gap: 12,
              alignItems: 'center',
            }}
          >
            <div
              style={{
                width: 56,
                height: 56,
                borderRadius: 16,
                background: SB.s100,
                display: 'grid',
                placeItems: 'center',
                flexShrink: 0,
              }}
            >
              {i === 0 ? (
                <SBMini size={42} expr="happy" />
              ) : i === 1 ? (
                <SBChick size={42} expr="happy" />
              ) : (
                <SBBear size={42} expr="plain" />
              )}
            </div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <div style={{ fontSize: 14, fontWeight: 700, flex: 1 }}>{t.title}</div>
                <SBStatusPill label={t.status} size="sm" />
              </div>
              <div
                style={{
                  fontSize: 11,
                  color: SB.w500,
                  marginTop: 3,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 4,
                }}
              >
                <SBIcon d={SB_ICONS.pin} size={11} color={SB.s500} />
                {t.destination}
              </div>
              <div
                style={{
                  fontSize: 11,
                  color: SB.w500,
                  marginTop: 2,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 4,
                }}
              >
                <SBIcon d={SB_ICONS.cal} size={11} color={SB.s500} />
                {t.start} → {t.end}
              </div>
            </div>
          </div>
        ))}
      </div>

      <SBFab icon={SB_ICONS.plus} label="새 여행" />
      <SBTabBar active="home" />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// 2) Trip Detail
// ─────────────────────────────────────────────────────────────
function ScreenTripDetail({
  onBack = () => {},
  onOpenDay = () => {},
  onOpenNotes = () => {},
  onOpenExpenses = () => {},
}) {
  const t = SB_TRIPS[0]
  return (
    <div
      style={{
        background: SB.w50,
        minHeight: '100%',
        fontFamily: SB_FONT,
        color: SB.w900,
        paddingBottom: 92,
        position: 'relative',
      }}
    >
      <SBAppBar
        title="여행 상세"
        onBack={onBack}
        right={
          <button
            style={{
              width: 40,
              height: 40,
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              color: SB.s600,
            }}
          >
            <SBIcon d={SB_ICONS.more} size={22} />
          </button>
        }
      />

      {/* hero header */}
      <div
        style={{
          margin: '0 16px 12px',
          borderRadius: 22,
          overflow: 'hidden',
          background: `linear-gradient(135deg, ${SB.s100} 0%, ${SB.s50} 60%, #fff 100%)`,
          border: `1.5px solid ${SB.s200}`,
          position: 'relative',
        }}
      >
        <SBPhoto h={120} label="cover · sakura kyoto" tint={SB.s100} stroke={SB.s200} radius={0} />
        <div style={{ position: 'absolute', right: 8, top: 64 }}>
          <SBBear size={84} expr="wink" accessory="flower" />
        </div>
        <div style={{ padding: '14px 16px 16px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
            <div style={{ fontSize: 22, fontWeight: 700, flex: 1, color: SB.s900 }}>{t.title}</div>
            <SBStatusPill label={t.status} />
          </div>
          <div style={{ display: 'flex', gap: 8, marginTop: 6, flexWrap: 'wrap', maxWidth: 240 }}>
            <SBChip icon={<SBIcon d={SB_ICONS.pin} size={10} />}>{t.destination}</SBChip>
            <SBChip bg="#fff" fg={SB.s700} icon={<SBIcon d={SB_ICONS.cal} size={10} />}>
              {t.days}박 6일
            </SBChip>
          </div>
          <div style={{ marginTop: 10, fontSize: 12, color: SB.w600 }}>
            {t.start} → {t.end}
          </div>
        </div>
      </div>

      {/* progress strip */}
      <div
        style={{
          margin: '0 16px 12px',
          padding: 12,
          borderRadius: 14,
          background: '#fff',
          border: `1px solid ${SB.s100}`,
          display: 'flex',
          alignItems: 'center',
          gap: 12,
        }}
      >
        <SBChick size={36} expr="happy" />
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 12, fontWeight: 700 }}>여행 준비 {Math.round(t.prep * 100)}%</div>
          <div
            style={{
              marginTop: 4,
              height: 5,
              borderRadius: 99,
              background: SB.s100,
              overflow: 'hidden',
            }}
          >
            <div
              style={{
                width: `${t.prep * 100}%`,
                height: '100%',
                background: `linear-gradient(90deg, ${SB.s400}, ${SB.s500})`,
              }}
            />
          </div>
          <div style={{ marginTop: 4, fontSize: 10, color: SB.w500 }}>
            항공권 ✓ · 숙소 ✓ · 환전 · 짐싸기
          </div>
        </div>
      </div>

      {/* sub-nav cards */}
      <div
        style={{ padding: '0 16px 4px', display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10 }}
      >
        {[
          { l: '메모', sub: `${SB_NOTES.length}건`, d: SB_ICONS.note, onClick: onOpenNotes },
          { l: '지출', sub: '113,140 JPY', d: SB_ICONS.wallet, onClick: onOpenExpenses },
        ].map((q, i) => (
          <button
            key={i}
            onClick={q.onClick}
            style={{
              padding: 14,
              borderRadius: 16,
              background: SB.s50,
              border: `1px solid ${SB.s100}`,
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: 10,
              fontFamily: SB_FONT,
            }}
          >
            <div
              style={{
                width: 36,
                height: 36,
                borderRadius: 12,
                background: '#fff',
                display: 'grid',
                placeItems: 'center',
                color: SB.s600,
              }}
            >
              <SBIcon d={q.d} size={18} />
            </div>
            <div style={{ textAlign: 'left' }}>
              <div style={{ fontSize: 13, fontWeight: 700, color: SB.s900 }}>{q.l}</div>
              <div style={{ fontSize: 10, color: SB.w500, marginTop: 1 }}>{q.sub}</div>
            </div>
          </button>
        ))}
      </div>

      <SBSection
        icon={SB_ICONS.cal}
        title="일정"
        count={`${SB_DAYS.length}일`}
        action={
          <button
            style={{
              background: 'none',
              border: 'none',
              color: SB.s600,
              fontWeight: 700,
              fontSize: 11,
              fontFamily: SB_FONT,
              cursor: 'pointer',
            }}
          >
            + 추가
          </button>
        }
      />

      <div style={{ padding: '0 16px', display: 'flex', flexDirection: 'column', gap: 8 }}>
        {SB_DAYS.map((d) => (
          <div
            key={d.id}
            onClick={() => onOpenDay(d.id)}
            style={{
              padding: 12,
              borderRadius: 14,
              background: '#fff',
              cursor: 'pointer',
              border: `1px solid ${SB.s100}`,
              display: 'flex',
              alignItems: 'center',
              gap: 12,
            }}
          >
            <div
              style={{
                width: 44,
                height: 44,
                borderRadius: 99,
                background: SB.s100,
                color: SB.s700,
                display: 'grid',
                placeItems: 'center',
                flexShrink: 0,
                border: `1.5px dashed ${SB.s300}`,
              }}
            >
              <div style={{ textAlign: 'center', lineHeight: 1 }}>
                <div style={{ fontSize: 9, fontWeight: 700, opacity: 0.7 }}>DAY</div>
                <div style={{ fontSize: 14, fontWeight: 700 }}>{d.n}</div>
              </div>
            </div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{ display: 'flex', alignItems: 'baseline', gap: 6 }}>
                <div style={{ fontSize: 13, fontWeight: 700 }}>{d.date}</div>
                <div style={{ fontSize: 11, color: SB.w500 }}>· {d.region}</div>
              </div>
              <div
                style={{
                  fontSize: 11,
                  color: SB.w600,
                  marginTop: 2,
                  whiteSpace: 'nowrap',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                }}
              >
                {d.summary}
              </div>
              <div style={{ display: 'flex', gap: 6, marginTop: 6 }}>
                <SBChip bg={SB.w100} fg={SB.w600} size="sm">
                  일정 {d.sched}
                </SBChip>
                <SBChip bg={SB.w100} fg={SB.w600} size="sm">
                  식사 {d.meals}
                </SBChip>
                {d.hasAcc && (
                  <SBChip bg={SB.s100} fg={SB.s600} size="sm">
                    숙박
                  </SBChip>
                )}
              </div>
            </div>
            <SBIcon d="M9 6l6 6-6 6" color={SB.w400} size={16} />
          </div>
        ))}
      </div>

      <div style={{ height: 24 }} />
      <SBTabBar active="trips" />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// 3) Day Detail
// ─────────────────────────────────────────────────────────────
function ScreenDayDetail({ onBack = () => {} }) {
  return (
    <div
      style={{
        background: SB.w50,
        minHeight: '100%',
        fontFamily: SB_FONT,
        color: SB.w900,
        paddingBottom: 92,
        position: 'relative',
      }}
    >
      <SBAppBar
        title="Day 2 · 04.13 (월)"
        onBack={onBack}
        right={
          <button
            style={{
              width: 40,
              height: 40,
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              color: SB.s600,
            }}
          >
            <SBIcon d={SB_ICONS.edit} size={20} />
          </button>
        }
      />

      {/* day header */}
      <div
        style={{
          margin: '0 16px 12px',
          padding: 14,
          borderRadius: 18,
          background: `linear-gradient(135deg, ${SB.s100}, ${SB.s50})`,
          border: `1.5px solid ${SB.s200}`,
          position: 'relative',
          overflow: 'hidden',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <SBBear size={48} expr="happy" accessory="camera" />
          <div>
            <div style={{ fontSize: 11, color: SB.s600, fontWeight: 700, letterSpacing: 1 }}>
              DAY 2 · 교토 동부
            </div>
            <div style={{ fontSize: 16, fontWeight: 700, marginTop: 2 }}>
              기요미즈데라 + 기온 산책
            </div>
          </div>
        </div>
        <div style={{ display: 'flex', gap: 6, marginTop: 10 }}>
          <SBChip bg="#fff" fg={SB.s700}>
            ☀ 17°C
          </SBChip>
          <SBChip bg="#fff" fg={SB.s700}>
            걸음 22,000+
          </SBChip>
          <SBChip bg={SB.s500} fg="#fff">
            📷 12장
          </SBChip>
        </div>
      </div>

      {/* schedules */}
      <SBSection icon={SB_ICONS.list} title="일정" count={`${SB_SCHEDULES.length}건`} />
      <div style={{ padding: '0 16px', position: 'relative' }}>
        <div
          style={{
            position: 'absolute',
            left: 38,
            top: 6,
            bottom: 6,
            width: 2,
            background: `repeating-linear-gradient(${SB.s200} 0 4px, transparent 4px 8px)`,
          }}
        />
        {SB_SCHEDULES.map((s, i) => (
          <div
            key={s.id}
            style={{ display: 'flex', gap: 12, marginBottom: 10, position: 'relative' }}
          >
            <div
              style={{
                width: 50,
                paddingTop: 12,
                fontSize: 10,
                color: SB.s700,
                fontWeight: 700,
                textAlign: 'center',
                flexShrink: 0,
              }}
            >
              {s.start}
            </div>
            <div style={{ width: 14, position: 'relative', flexShrink: 0 }}>
              <div
                style={{
                  position: 'absolute',
                  left: 4,
                  top: 14,
                  width: 10,
                  height: 10,
                  borderRadius: 99,
                  background: SB.s500,
                  border: '2px solid #fff',
                  boxShadow: `0 0 0 2px ${SB.s200}`,
                }}
              />
            </div>
            <div
              style={{
                flex: 1,
                padding: 12,
                borderRadius: 14,
                background: '#fff',
                border: `1px solid ${SB.s100}`,
              }}
            >
              <div style={{ fontSize: 13, fontWeight: 700 }}>{s.title}</div>
              <div
                style={{
                  fontSize: 11,
                  color: SB.w500,
                  marginTop: 2,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 4,
                }}
              >
                <SBIcon d={SB_ICONS.pin} size={11} color={SB.s500} />
                {s.place}
              </div>
              {s.notes && (
                <div
                  style={{
                    marginTop: 6,
                    padding: 8,
                    background: SB.s50,
                    borderRadius: 8,
                    fontSize: 11,
                    color: SB.w700,
                  }}
                >
                  {s.notes}
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* meals */}
      <SBSection icon={SB_ICONS.food} title="식사" count={`${SB_MEALS.length}건`} />
      <div style={{ padding: '0 16px', display: 'flex', flexDirection: 'column', gap: 8 }}>
        {SB_MEALS.map((m) => (
          <div
            key={m.id}
            style={{
              padding: 12,
              borderRadius: 14,
              background: '#fff',
              border: `1px solid ${SB.s100}`,
              display: 'flex',
              gap: 12,
              alignItems: 'center',
            }}
          >
            <div
              style={{
                width: 40,
                height: 40,
                borderRadius: 12,
                background: SB.s50,
                display: 'grid',
                placeItems: 'center',
              }}
            >
              <SBIcon d={SB_ICONS.food} size={18} color={SB.s600} />
            </div>
            <div style={{ flex: 1 }}>
              <div style={{ display: 'flex', alignItems: 'baseline', gap: 6 }}>
                <SBChip bg={SB.s100} fg={SB.s700} size="sm">
                  {m.type}
                </SBChip>
                <div style={{ fontSize: 13, fontWeight: 700 }}>{m.name}</div>
              </div>
              {m.review && (
                <div style={{ fontSize: 11, color: SB.w500, marginTop: 4 }}>{m.review}</div>
              )}
            </div>
            <div style={{ fontSize: 11, fontWeight: 700, color: SB.s700 }}>{m.cost}</div>
          </div>
        ))}
      </div>

      {/* accommodation */}
      <SBSection icon={SB_ICONS.bed} title="숙소" />
      <div style={{ padding: '0 16px 16px' }}>
        <div
          style={{
            padding: 14,
            borderRadius: 14,
            background: '#fff',
            border: `1px solid ${SB.s100}`,
          }}
        >
          <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
            <SBPhoto w={56} h={56} label="hotel" tint={SB.s100} stroke={SB.s200} radius={12} />
            <div style={{ flex: 1 }}>
              <div style={{ fontSize: 14, fontWeight: 700 }}>{SB_ACC.name}</div>
              <div
                style={{
                  fontSize: 11,
                  color: SB.w500,
                  marginTop: 2,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 4,
                }}
              >
                <SBIcon d={SB_ICONS.pin} size={11} color={SB.s500} />
                {SB_ACC.address}
              </div>
            </div>
          </div>
          <div
            style={{
              marginTop: 10,
              padding: 8,
              background: SB.s50,
              borderRadius: 10,
              display: 'flex',
              justifyContent: 'space-between',
              fontSize: 11,
              color: SB.w700,
            }}
          >
            <span>
              체크인 <b style={{ color: SB.s700 }}>{SB_ACC.checkin}</b>
            </span>
            <span>
              체크아웃 <b style={{ color: SB.s700 }}>{SB_ACC.checkout}</b>
            </span>
            <span style={{ color: SB.s700, fontWeight: 700 }}>{SB_ACC.cost}</span>
          </div>
        </div>
      </div>

      <SBFab icon={SB_ICONS.plus} label="기록 추가" />
      <SBTabBar active="trips" />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// 4) Notes
// ─────────────────────────────────────────────────────────────
const MOOD_META = {
  설렘: { color: '#F25C7A', bg: SB.s100, emoji: '✨' },
  맛있음: { color: '#D97706', bg: '#FED7AA', emoji: '🍙' },
  평온: { color: '#0EA5E9', bg: SB.sky100, emoji: '🌿' },
  피곤: { color: '#7C3AED', bg: '#EDE9FE', emoji: '💤' },
}

function ScreenNotes({ onBack = () => {} }) {
  return (
    <div
      style={{
        background: SB.w50,
        minHeight: '100%',
        fontFamily: SB_FONT,
        color: SB.w900,
        paddingBottom: 92,
        position: 'relative',
      }}
    >
      <SBAppBar
        title="여행 메모"
        onBack={onBack}
        right={
          <div style={{ display: 'flex', gap: 4 }}>
            <button
              style={{
                width: 40,
                height: 40,
                border: 'none',
                background: 'transparent',
                cursor: 'pointer',
                color: SB.s600,
              }}
            >
              <SBIcon d={SB_ICONS.search} size={20} />
            </button>
            <button
              style={{
                width: 40,
                height: 40,
                border: 'none',
                background: 'transparent',
                cursor: 'pointer',
                color: SB.s600,
              }}
            >
              <SBIcon d={SB_ICONS.filter} size={20} />
            </button>
          </div>
        }
      />

      {/* mood filter pills */}
      <div style={{ padding: '4px 16px 12px', display: 'flex', gap: 6, overflowX: 'auto' }}>
        {['전체', ...Object.keys(MOOD_META)].map((m, i) => (
          <div
            key={m}
            style={{
              padding: '6px 12px',
              borderRadius: 99,
              background: i === 0 ? SB.s500 : '#fff',
              color: i === 0 ? '#fff' : SB.w700,
              border: i === 0 ? 'none' : `1px solid ${SB.s100}`,
              fontSize: 11,
              fontWeight: 700,
              whiteSpace: 'nowrap',
              display: 'flex',
              alignItems: 'center',
              gap: 4,
            }}
          >
            {i > 0 && <span>{MOOD_META[m].emoji}</span>}#{m}
          </div>
        ))}
      </div>

      {/* sticky note style cards */}
      <div style={{ padding: '0 16px', display: 'flex', flexDirection: 'column', gap: 12 }}>
        {SB_NOTES.map((n, i) => {
          const meta = MOOD_META[n.mood]
          const rotate = i % 2 === 0 ? -0.6 : 0.4
          return (
            <div
              key={n.id}
              style={{
                padding: 14,
                borderRadius: 18,
                background: '#fff',
                border: `1px solid ${SB.s100}`,
                borderLeft: `4px solid ${meta.color}`,
                transform: `rotate(${rotate}deg)`,
                boxShadow: '0 2px 8px rgba(96,28,41,0.06)',
              }}
            >
              <div style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <SBChip bg={meta.bg} fg={meta.color} size="sm">
                  {meta.emoji} {n.mood}
                </SBChip>
                <div style={{ fontSize: 10, color: SB.w400 }}>{n.date}</div>
                <div style={{ flex: 1 }} />
                <SBIcon d={SB_ICONS.heart} size={14} color={SB.s400} />
              </div>
              <div style={{ fontSize: 13, color: SB.w800, marginTop: 8, lineHeight: 1.55 }}>
                {n.content}
              </div>
            </div>
          )
        })}

        {/* paw decoration */}
        <div
          style={{ display: 'flex', gap: 8, padding: 16, justifyContent: 'center', opacity: 0.35 }}
        >
          <SBPaw size={14} color={SB.s400} />
          <SBPaw size={14} color={SB.s300} />
          <SBPaw size={14} color={SB.s400} />
        </div>
      </div>

      <SBFab icon={SB_ICONS.edit} label="새 메모" />
      <SBTabBar active="trips" />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// 5) Expenses
// ─────────────────────────────────────────────────────────────
function ScreenExpenses({ onBack = () => {} }) {
  const total = 113140
  const byCat = {
    교통: 3640,
    숙박: 96000,
    식사: 10300,
    체험: 400,
    쇼핑: 2200,
  }
  const max = Math.max(...Object.values(byCat))

  return (
    <div
      style={{
        background: SB.w50,
        minHeight: '100%',
        fontFamily: SB_FONT,
        color: SB.w900,
        paddingBottom: 92,
        position: 'relative',
      }}
    >
      <SBAppBar
        title="지출"
        onBack={onBack}
        right={
          <button
            style={{
              width: 40,
              height: 40,
              border: 'none',
              background: 'transparent',
              cursor: 'pointer',
              color: SB.s600,
            }}
          >
            <SBIcon d={SB_ICONS.download} size={20} />
          </button>
        }
      />

      {/* total card */}
      <div
        style={{
          margin: '0 16px 12px',
          padding: 16,
          borderRadius: 22,
          background: `linear-gradient(135deg, ${SB.s500} 0%, ${SB.s400} 100%)`,
          color: '#fff',
          position: 'relative',
          overflow: 'hidden',
          boxShadow: '0 8px 24px rgba(242,92,122,0.25)',
        }}
      >
        <div style={{ fontSize: 11, opacity: 0.85, fontWeight: 700, letterSpacing: 1 }}>
          총 지출 · 4박 6일
        </div>
        <div style={{ fontSize: 32, fontWeight: 700, marginTop: 4 }}>
          {total.toLocaleString()} <span style={{ fontSize: 14, opacity: 0.85 }}>JPY</span>
        </div>
        <div style={{ fontSize: 12, opacity: 0.85, marginTop: 2 }}>
          ≈ ₩ 1,068,000 · 1인당 평균 22,628 JPY
        </div>
        <div style={{ position: 'absolute', right: 8, bottom: -4 }}>
          <SBBear size={68} expr="surprise" />
        </div>
      </div>

      {/* category breakdown */}
      <SBSection title="카테고리별" />
      <div style={{ padding: '0 16px 4px', display: 'flex', flexDirection: 'column', gap: 8 }}>
        {Object.entries(byCat).map(([cat, val]) => {
          const meta = CAT_META[cat]
          const pct = val / max
          return (
            <div key={cat} style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
              <div
                style={{
                  width: 28,
                  height: 28,
                  borderRadius: 8,
                  background: meta.bg,
                  display: 'grid',
                  placeItems: 'center',
                  color: meta.color,
                  flexShrink: 0,
                }}
              >
                <SBIcon d={meta.icon} size={14} />
              </div>
              <div style={{ flex: 1 }}>
                <div style={{ display: 'flex', alignItems: 'baseline' }}>
                  <div style={{ fontSize: 12, fontWeight: 700, flex: 1 }}>{cat}</div>
                  <div style={{ fontSize: 12, fontWeight: 700, color: meta.color }}>
                    {val.toLocaleString()} JPY
                  </div>
                </div>
                <div
                  style={{
                    marginTop: 4,
                    height: 5,
                    borderRadius: 99,
                    background: SB.w100,
                    overflow: 'hidden',
                  }}
                >
                  <div style={{ width: `${pct * 100}%`, height: '100%', background: meta.color }} />
                </div>
              </div>
            </div>
          )
        })}
      </div>

      <SBSection title="기록" count={`${SB_EXPENSES.length}건`} />
      <div style={{ padding: '0 16px', display: 'flex', flexDirection: 'column', gap: 6 }}>
        {SB_EXPENSES.map((e) => {
          const meta = CAT_META[e.cat] || CAT_META['기타']
          return (
            <div
              key={e.id}
              style={{
                padding: 12,
                borderRadius: 14,
                background: '#fff',
                border: `1px solid ${SB.s100}`,
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
                  background: meta.bg,
                  display: 'grid',
                  placeItems: 'center',
                  color: meta.color,
                  flexShrink: 0,
                }}
              >
                <SBIcon d={meta.icon} size={16} />
              </div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ display: 'flex', gap: 6, alignItems: 'baseline' }}>
                  <SBChip bg={meta.bg} fg={meta.color} size="sm">
                    {e.cat}
                  </SBChip>
                  <div style={{ fontSize: 10, color: SB.w400 }}>{e.date}</div>
                </div>
                <div
                  style={{
                    fontSize: 12,
                    color: SB.w800,
                    marginTop: 3,
                    whiteSpace: 'nowrap',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                  }}
                >
                  {e.desc}
                </div>
              </div>
              <div style={{ fontSize: 13, fontWeight: 700, color: SB.s700 }}>{e.amount}</div>
            </div>
          )
        })}
      </div>

      <SBFab icon={SB_ICONS.plus} label="지출" />
      <SBTabBar active="trips" />
    </div>
  )
}

// ─────────────────────────────────────────────────────────────
// 6) Login
// ─────────────────────────────────────────────────────────────
function ScreenLogin() {
  return (
    <div
      style={{
        background: `radial-gradient(circle at 50% 30%, ${SB.s100} 0%, ${SB.s50} 40%, ${SB.cream} 100%)`,
        minHeight: '100%',
        fontFamily: SB_FONT,
        color: SB.w900,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        padding: '40px 32px 24px',
        position: 'relative',
        overflow: 'hidden',
      }}
    >
      {/* floating petals */}
      {[
        [30, 80, -20, 14],
        [340, 60, 30, 12],
        [60, 200, 60, 16],
        [310, 300, -30, 14],
        [40, 380, 10, 10],
        [330, 460, -50, 12],
      ].map(([x, y, r, s], i) => (
        <div key={i} style={{ position: 'absolute', left: x, top: y, opacity: 0.55 }}>
          <SBPetal size={s} rotate={r} color={i % 2 ? SB.s400 : SB.s300} />
        </div>
      ))}

      <div style={{ marginTop: 60 }} />
      {/* hero mascot scene */}
      <div
        style={{
          position: 'relative',
          width: 220,
          height: 180,
          display: 'flex',
          alignItems: 'flex-end',
          justifyContent: 'center',
        }}
      >
        <div style={{ position: 'absolute', left: 0, bottom: 8 }}>
          <SBChick size={56} expr="happy" />
        </div>
        <SBBear size={148} expr="happy" accessory="flower" />
        <div style={{ position: 'absolute', right: -4, bottom: 24 }}>
          <SBMini size={70} expr="happy" />
        </div>
        {/* ground */}
        <div
          style={{
            position: 'absolute',
            left: 0,
            right: 0,
            bottom: 0,
            height: 14,
            borderRadius: '50% 50% 0 0',
            background: SB.s100,
            zIndex: -1,
          }}
        />
      </div>

      <div
        style={{ marginTop: 24, fontSize: 13, color: SB.s600, fontWeight: 700, letterSpacing: 2 }}
      >
        SEONOLOGY · JOURNEY
      </div>
      <h1
        style={{
          margin: '8px 0 0',
          fontSize: 30,
          fontWeight: 700,
          color: SB.s900,
          textAlign: 'center',
          lineHeight: 1.25,
        }}
      >
        旅の記録,
        <br />
        곰돌이와 함께.
      </h1>
      <p
        style={{
          margin: '12px 0 0',
          fontSize: 13,
          color: SB.w600,
          textAlign: 'center',
          lineHeight: 1.6,
          maxWidth: 280,
        }}
      >
        여행 계획·일정·사진·메모·지출까지
        <br />한 권의 폭신폭신한 노트에.
      </p>

      <div style={{ flex: 1 }} />

      <button
        style={{
          width: '100%',
          height: 54,
          borderRadius: 18,
          background: SB.s500,
          color: '#fff',
          fontFamily: SB_FONT,
          fontSize: 15,
          fontWeight: 700,
          border: 'none',
          cursor: 'pointer',
          boxShadow: '0 8px 20px rgba(242,92,122,0.3)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: 8,
        }}
      >
        <SBIcon d="M9 12l2 2 4-4m7 2a9 9 0 11-18 0 9 9 0 0118 0z" size={18} color="#fff" />
        Keycloak으로 로그인
      </button>

      <button
        style={{
          marginTop: 10,
          width: '100%',
          height: 48,
          borderRadius: 16,
          background: '#fff',
          color: SB.s700,
          fontFamily: SB_FONT,
          fontSize: 13,
          fontWeight: 700,
          border: `1px solid ${SB.s100}`,
          cursor: 'pointer',
        }}
      >
        둘러보기 (게스트)
      </button>

      <div style={{ marginTop: 14, fontSize: 10, color: SB.w400, textAlign: 'center' }}>
        계속하시면 <span style={{ color: SB.s600, fontWeight: 700 }}>이용약관</span> 및{' '}
        <span style={{ color: SB.s600, fontWeight: 700 }}>개인정보처리방침</span>에 동의하게 됩니다.
      </div>
    </div>
  )
}

Object.assign(window, {
  ScreenTripList,
  ScreenTripDetail,
  ScreenDayDetail,
  ScreenNotes,
  ScreenExpenses,
  ScreenLogin,
})

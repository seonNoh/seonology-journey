import { Plane } from 'lucide-react'

/**
 * CuteLoader - 벚꽃(sakura) 테마에 맞춘 귀엽고 캐주얼한 로딩 UI.
 *
 * End-to-end 자동 재배포 플로우 확인용 더미 터치: 이 주석이 배포된
 * 파드에도 있으면 CI → ghcr → Image Updater → ArgoCD sync 가 모두
 * 동작한 것이다.
 *
 * 재사용 블록 3 가지를 제공:
 *  - <CuteLoader />          : inline 스피너 (여행 중인 비행기 + 흔들리는 벚꽃)
 *  - <CuteLoaderBlock />     : 카드 안에 넣는 세로 중앙 정렬 블록
 *  - <FlyingPlane />         : 일직선으로 날아가는 비행기 애니메이션 (텍스트 아래 장식용)
 *
 * Tailwind 의 animate-* 유틸을 쓰기보다, 프로젝트 전반의 "느린 pulse"
 * 느낌을 유지하기 위해 인라인 <style> 로 keyframes 를 선언한다. 짧고
 * 독립적이라 기존 CSS 번들을 건드리지 않는 쪽이 깔끔하다.
 */
export function CuteLoader({
  message,
  size = 'md',
}: {
  message?: string
  size?: 'sm' | 'md' | 'lg'
}) {
  const sizeMap = {
    sm: { plane: 'h-4 w-4', petal: 12, gap: 'gap-1.5', text: 'text-xs' },
    md: { plane: 'h-5 w-5', petal: 14, gap: 'gap-2', text: 'text-sm' },
    lg: { plane: 'h-6 w-6', petal: 16, gap: 'gap-3', text: 'text-base' },
  }[size]

  return (
    <div className={`inline-flex items-center ${sizeMap.gap} text-sakura-700`} role="status">
      <style>{cuteLoaderKeyframes}</style>
      <span className="relative inline-flex items-center justify-center">
        <Petal size={sizeMap.petal} className="cute-petal cute-petal-a" />
        <Petal size={sizeMap.petal} className="cute-petal cute-petal-b" />
        <Plane className={`${sizeMap.plane} cute-plane text-sakura-500`} aria-hidden="true" />
      </span>
      {message && <span className={`${sizeMap.text} text-slate-600`}>{message}</span>}
      <span className="sr-only">로딩 중</span>
    </div>
  )
}

export function CuteLoaderBlock({ message = '준비 중이에요' }: { message?: string }) {
  return (
    <div
      className="flex flex-col items-center justify-center gap-3 rounded-2xl bg-white/80 px-6 py-10 shadow-sm"
      role="status"
      aria-live="polite"
    >
      <style>{cuteLoaderKeyframes}</style>
      <span className="relative inline-flex h-12 w-12 items-center justify-center">
        <Petal size={18} className="cute-petal cute-petal-a" />
        <Petal size={18} className="cute-petal cute-petal-b" />
        <Petal size={14} className="cute-petal cute-petal-c" />
        <Plane className="h-6 w-6 cute-plane text-sakura-500" aria-hidden="true" />
      </span>
      <p className="text-sm text-slate-600">{message}</p>
    </div>
  )
}

/** 얇은 상단 진행바 느낌의 라인 로더. 페이지 전환 없이 섹션만 새로고침될 때 사용. */
export function FlyingPlane({ className = '' }: { className?: string }) {
  return (
    <div
      className={`relative h-6 w-full overflow-hidden rounded-full bg-sakura-50 ${className}`}
      role="status"
      aria-label="불러오는 중"
    >
      <style>{cuteLoaderKeyframes}</style>
      <Plane
        className="absolute top-1/2 h-4 w-4 -translate-y-1/2 cute-plane-fly text-sakura-500"
        aria-hidden="true"
      />
    </div>
  )
}

function Petal({ size, className = '' }: { size: number; className?: string }) {
  // Hand-drawn sakura petal: one scoop + a subtle notch at the tip.
  return (
    <svg
      className={className}
      viewBox="0 0 24 24"
      width={size}
      height={size}
      fill="currentColor"
      aria-hidden="true"
    >
      <path
        d="M12 2c2.8 2.6 5.5 5.5 5.5 9.2 0 3.6-2.9 6.3-5.5 8.3-2.6-2-5.5-4.7-5.5-8.3C6.5 7.5 9.2 4.6 12 2Zm0 6.5c-.6.5-1.1 1.2-1.1 2.1 0 .8.5 1.5 1.1 2.1.6-.6 1.1-1.3 1.1-2.1 0-.9-.5-1.6-1.1-2.1Z"
        className="text-sakura-300"
      />
    </svg>
  )
}

// Keep keyframes in a string so any Loader usage ships with its
// animations even if a page is mounted before the main CSS bundle.
const cuteLoaderKeyframes = `
  @keyframes cute-plane-bob {
    0%, 100% { transform: translateY(0) rotate(-8deg); }
    50%      { transform: translateY(-3px) rotate(8deg); }
  }
  @keyframes cute-petal-fall-a {
    0%   { transform: translate(-12px, -14px) rotate(0deg); opacity: 0; }
    20%  { opacity: 1; }
    100% { transform: translate(14px, 18px) rotate(240deg); opacity: 0; }
  }
  @keyframes cute-petal-fall-b {
    0%   { transform: translate(10px, -16px) rotate(0deg); opacity: 0; }
    25%  { opacity: 1; }
    100% { transform: translate(-14px, 16px) rotate(-200deg); opacity: 0; }
  }
  @keyframes cute-petal-fall-c {
    0%   { transform: translate(0, -18px) rotate(0deg); opacity: 0; }
    30%  { opacity: 0.9; }
    100% { transform: translate(6px, 22px) rotate(360deg); opacity: 0; }
  }
  @keyframes cute-plane-fly {
    0%   { left: -12%; transform: translateY(-50%) rotate(-4deg); }
    55%  { transform: translateY(-60%) rotate(6deg); }
    100% { left: 100%; transform: translateY(-40%) rotate(-2deg); }
  }
  .cute-plane      { animation: cute-plane-bob 1.6s ease-in-out infinite; transform-origin: center; }
  .cute-petal      { position: absolute; color: #ffc7d4; }
  .cute-petal-a    { animation: cute-petal-fall-a 2.4s ease-in infinite; }
  .cute-petal-b    { animation: cute-petal-fall-b 2.8s ease-in infinite 0.3s; }
  .cute-petal-c    { animation: cute-petal-fall-c 3.2s ease-in infinite 0.6s; }
  .cute-plane-fly  { animation: cute-plane-fly 2.8s ease-in-out infinite; }
`

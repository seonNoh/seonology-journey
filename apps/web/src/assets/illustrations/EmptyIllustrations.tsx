export function NoTripsIllustration({ className = '' }: { className?: string }) {
  return (
    <svg
      width="200"
      height="160"
      viewBox="0 0 200 160"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      aria-hidden="true"
    >
      {/* Suitcase */}
      <rect
        x="60"
        y="50"
        width="80"
        height="60"
        rx="8"
        stroke="#94a3b8"
        strokeWidth="2"
        fill="#f1f5f9"
      />
      <rect
        x="75"
        y="40"
        width="50"
        height="14"
        rx="4"
        stroke="#94a3b8"
        strokeWidth="2"
        fill="none"
      />
      <line
        x1="100"
        y1="70"
        x2="100"
        y2="90"
        stroke="#cbd5e1"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <line
        x1="90"
        y1="80"
        x2="110"
        y2="80"
        stroke="#cbd5e1"
        strokeWidth="2"
        strokeLinecap="round"
      />
      {/* Ground */}
      <ellipse cx="100" cy="120" rx="50" ry="6" fill="#e2e8f0" />
      {/* Plane trail */}
      <path
        d="M30 30 Q60 20 80 35"
        stroke="#94a3b8"
        strokeWidth="1.5"
        strokeDasharray="4 3"
        fill="none"
      />
      {/* Plane */}
      <path d="M25 32l8-3-2 5-6-2z" fill="#64748b" />
    </svg>
  )
}

export function NoScheduleIllustration({ className = '' }: { className?: string }) {
  return (
    <svg
      width="200"
      height="160"
      viewBox="0 0 200 160"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      aria-hidden="true"
    >
      {/* Calendar */}
      <rect
        x="50"
        y="40"
        width="100"
        height="80"
        rx="8"
        stroke="#94a3b8"
        strokeWidth="2"
        fill="#f1f5f9"
      />
      <line x1="50" y1="60" x2="150" y2="60" stroke="#94a3b8" strokeWidth="2" />
      <circle cx="70" cy="50" r="3" fill="#94a3b8" />
      <circle cx="130" cy="50" r="3" fill="#94a3b8" />
      {/* Empty rows */}
      <rect x="60" y="70" width="30" height="6" rx="3" fill="#e2e8f0" />
      <rect x="60" y="85" width="50" height="6" rx="3" fill="#e2e8f0" />
      <rect x="60" y="100" width="40" height="6" rx="3" fill="#e2e8f0" />
      {/* Dashed circle */}
      <circle
        cx="130"
        cy="95"
        r="15"
        stroke="#cbd5e1"
        strokeWidth="1.5"
        strokeDasharray="4 3"
        fill="none"
      />
      <line
        x1="130"
        y1="88"
        x2="130"
        y2="102"
        stroke="#cbd5e1"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <line
        x1="123"
        y1="95"
        x2="137"
        y2="95"
        stroke="#cbd5e1"
        strokeWidth="2"
        strokeLinecap="round"
      />
    </svg>
  )
}

export function NoPhotosIllustration({ className = '' }: { className?: string }) {
  return (
    <svg
      width="200"
      height="160"
      viewBox="0 0 200 160"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={className}
      aria-hidden="true"
    >
      {/* Camera body */}
      <rect
        x="55"
        y="50"
        width="90"
        height="65"
        rx="10"
        stroke="#94a3b8"
        strokeWidth="2"
        fill="#f1f5f9"
      />
      {/* Lens */}
      <circle cx="100" cy="82" r="20" stroke="#94a3b8" strokeWidth="2" fill="none" />
      <circle cx="100" cy="82" r="12" stroke="#cbd5e1" strokeWidth="1.5" fill="#e2e8f0" />
      {/* Flash */}
      <rect
        x="80"
        y="44"
        width="20"
        height="8"
        rx="3"
        fill="#e2e8f0"
        stroke="#94a3b8"
        strokeWidth="1.5"
      />
      {/* Shutter button */}
      <circle cx="125" cy="44" r="5" fill="#94a3b8" />
      {/* Ground shadow */}
      <ellipse cx="100" cy="125" rx="45" ry="5" fill="#e2e8f0" />
    </svg>
  )
}

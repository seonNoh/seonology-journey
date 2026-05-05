import { Plane, LogOut, User } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export function AppShell({ children }: { children: React.ReactNode }) {
  const auth = useAuth()
  return (
    <div className="min-h-screen bg-gradient-to-b from-sakura-50 to-white">
      <header
        className="border-b border-sakura-100 bg-white/80 backdrop-blur sticky top-0 z-10"
        role="banner"
      >
        <div className="mx-auto max-w-5xl flex items-center gap-3 px-4 py-3">
          <Link
            to="/"
            className="flex items-center gap-2 text-sakura-700 font-bold"
            aria-label="Seonology Journey Home"
          >
            <Plane className="h-5 w-5" aria-hidden="true" />
            Seonology Journey
          </Link>
          <nav
            className="flex-1 flex items-center gap-4 text-sm text-slate-700"
            aria-label="Main navigation"
          >
            <Link to="/trips" className="hover:text-sakura-600">
              여행
            </Link>
          </nav>
          <div className="flex items-center gap-2 text-sm">
            {auth.authenticated ? (
              <>
                <span className="flex items-center gap-1 text-slate-600">
                  <User className="h-4 w-4" />
                  {auth.username ?? 'me'}
                </span>
                <button
                  onClick={auth.logout}
                  className="flex items-center gap-1 rounded-md border border-sakura-200 px-2 py-1 text-sakura-700 hover:bg-sakura-50"
                >
                  <LogOut className="h-4 w-4" /> 로그아웃
                </button>
              </>
            ) : (
              <button
                onClick={auth.login}
                className="rounded-md bg-sakura-500 px-3 py-1.5 text-white hover:bg-sakura-600"
              >
                로그인
              </button>
            )}
          </div>
        </div>
      </header>
      <main id="main-content" className="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  )
}

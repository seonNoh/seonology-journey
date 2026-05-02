import { useAuth } from '../hooks/useAuth'

export function HomePage() {
  const auth = useAuth()
  return (
    <section className="space-y-6">
      <div className="rounded-2xl bg-white p-8 shadow-sm">
        <h1 className="text-2xl font-bold text-sakura-700">旅の記録, ここから.</h1>
        <p className="mt-2 text-slate-600">
          여행 계획부터 일정·식사·숙박·지출·기록까지 한 곳에서 관리.
        </p>
        {!auth.authenticated && (
          <button
            onClick={auth.login}
            className="mt-4 rounded-md bg-sakura-500 px-4 py-2 text-white hover:bg-sakura-600"
          >
            로그인하고 시작하기
          </button>
        )}
        {auth.authenticated && (
          <a
            href="/trips"
            className="mt-4 inline-block rounded-md bg-sakura-500 px-4 py-2 text-white hover:bg-sakura-600"
          >
            내 여행 보기
          </a>
        )}
      </div>
    </section>
  )
}

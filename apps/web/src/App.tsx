import { Plane } from 'lucide-react'

export function App() {
  return (
    <main className="min-h-screen flex items-center justify-center">
      <section className="rounded-3xl bg-white shadow-xl px-10 py-12 text-center">
        <Plane className="mx-auto mb-4 h-12 w-12 text-sakura-500" />
        <h1 className="text-3xl font-bold text-sakura-700">Seonology Journey</h1>
        <p className="mt-3 text-slate-600">旅の記録, ここから.</p>
      </section>
    </main>
  )
}

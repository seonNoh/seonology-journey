import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { App } from './App'
import './lib/i18n'
import './index.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      // 기본 staleTime 을 1 분으로 올려 화면 전환마다 다시 네트워크를 치지 않게 한다.
      // 개별 쿼리 (e.g. media URL) 에서 상황에 맞는 staleTime 으로 덮어쓴다.
      staleTime: 60_000,
      gcTime: 5 * 60_000,
    },
  },
})

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </React.StrictMode>,
)

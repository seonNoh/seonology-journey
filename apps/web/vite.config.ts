import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: { port: 5173 },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    css: true,
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      // mapbox-gl 은 런타임에 토큰이 있을 때만 동적으로 로드하므로 번들에서
      // 제외한다. index.html 에서 CDN 으로 로드되거나, 토큰이 없으면
      // MapPicker 가 fallback 플레이스홀더를 렌더링한다.
      external: ['mapbox-gl'],
      output: {
        manualChunks: {
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'auth-vendor': ['keycloak-js'],
          'charts-vendor': ['recharts'],
          'dnd-vendor': ['@dnd-kit/core', '@dnd-kit/sortable', '@dnd-kit/utilities'],
          'query-vendor': ['@tanstack/react-query', 'zustand'],
          'motion-vendor': ['framer-motion', 'lottie-react'],
          'i18n-vendor': ['i18next', 'react-i18next'],
        },
      },
    },
  },
})

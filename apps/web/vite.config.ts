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

import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

export default defineConfig({
  plugins: [react()],
  base: '/app/',
  server: { port: 5173, host: true },
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        landing: resolve(__dirname, 'src/landing/landing.ts'),
      },
      output: {
        entryFileNames: (chunk) => (chunk.name === 'landing' ? 'assets/landing.js' : 'assets/[name]-[hash].js'),
      },
    },
  },
})

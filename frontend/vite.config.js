import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    // This is the key fix for containerized environments.
    // It tells the HMR client to connect directly to the host machine,
    // which then forwards the connection to the container.
    hmr: {
      host: 'localhost',
      port: 5173,
      protocol: 'ws',
    },
  },
})
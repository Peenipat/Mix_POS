import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      'tailwindcss/version.js': path.resolve(__dirname, 'src/tailwind-version.js'),
      '@': path.resolve(__dirname, 'src'), 
       "@object/shared": path.resolve(__dirname, "../shared/src"), 
    },
  },
  css: {
    postcss: './postcss.config.cjs',
  },
})

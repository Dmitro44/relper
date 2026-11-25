import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/auth': 'http://localhost:3000',
      '/users': 'http://localhost:3000',
      '/locations': 'http://localhost:3000',
      '/places': 'http://localhost:3000',
      '/tips': 'http://localhost:3000',
      '/comments': 'http://localhost:3000',
      '/tags': 'http://localhost:3000',
      '/tasks': 'http://localhost:3000',
      '/audit': 'http://localhost:3000',
    },
  },
})

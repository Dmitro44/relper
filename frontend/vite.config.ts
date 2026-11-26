import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

const apiTarget = 'http://localhost:3000'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/auth': apiTarget,
      '/users': apiTarget,
      '/locations': apiTarget,
      '/places': apiTarget,
      '/place-suggestions': apiTarget,
      '/tips': apiTarget,
      '/comments': apiTarget,
      '/tags': apiTarget,
      '/tasks': apiTarget,
      '/audit': {
        target: apiTarget,
        bypass(req) {
          return req.headers.accept?.includes('text/html') ? '/index.html' : undefined
        },
      },
    },
  },
})

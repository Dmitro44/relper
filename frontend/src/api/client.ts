import createClient from 'openapi-fetch'
import type { paths } from './schema'
import { getToken } from './token'

export const api = createClient<paths>({ baseUrl: '' })

api.use({
  onRequest({ request }) {
    const token = getToken()
    if (token) {
      request.headers.set('Authorization', `Bearer ${token}`)
    }
    return request
  },
})

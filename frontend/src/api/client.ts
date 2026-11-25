import createClient from 'openapi-fetch'
import { refreshAccessToken } from './refresh'
import type { paths } from './schema'
import { getToken, setToken } from './token'

const rawFetch = window.fetch.bind(window)
const authPaths = new Set(['/auth/login', '/auth/register', '/auth/refresh', '/auth/logout'])

async function fetchWithRefresh(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  const request = new Request(input, init)
  const response = await rawFetch(request.clone())
  const path = new URL(request.url).pathname

  if (response.status !== 401 || authPaths.has(path)) {
    return response
  }

  const token = await refreshAccessToken()
  if (!token) {
    setToken('')
    window.dispatchEvent(new Event('relper:auth-expired'))
    return response
  }

  const retry = request.clone()
  retry.headers.set('Authorization', `Bearer ${token}`)
  return rawFetch(retry)
}

export const api = createClient<paths>({ baseUrl: '', fetch: fetchWithRefresh })

api.use({
  onRequest({ request }) {
    const token = getToken()
    if (token) {
      request.headers.set('Authorization', `Bearer ${token}`)
    }
    return request
  },
})

import { setToken } from './token'

let pending: Promise<string | null> | null = null

async function request(): Promise<string | null> {
  try {
    const response = await fetch('/auth/refresh', {
      method: 'POST',
      credentials: 'same-origin',
    })
    if (!response.ok) {
      return null
    }

    const body = (await response.json()) as { access_token?: string }
    if (!body.access_token) {
      return null
    }

    setToken(body.access_token)
    return body.access_token
  } catch {
    return null
  }
}

export function refreshAccessToken(): Promise<string | null> {
  if (!pending) {
    pending = request().finally(() => {
      pending = null
    })
  }
  return pending
}

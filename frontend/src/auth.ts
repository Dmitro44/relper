import { reactive } from 'vue'
import { api } from './api/client'
import { errorMessage } from './api/error'
import { getToken, setToken } from './api/token'
import type { Me } from './types'

interface AuthState {
  user: Me | null
  ready: boolean
  busy: boolean
  error: string
}

export const auth = reactive<AuthState>({
  user: null,
  ready: false,
  busy: false,
  error: '',
})

let pendingLoad: Promise<void> | null = null

export function isAdmin(): boolean {
  return auth.user?.roles?.includes('admin') ?? false
}

export async function loadMe(): Promise<void> {
  if (pendingLoad) {
    return pendingLoad
  }

  pendingLoad = (async () => {
    if (!getToken()) {
      auth.user = null
      auth.ready = true
      return
    }

    const { data, error } = await api.GET('/auth/me')
    if (error || !data) {
      setToken('')
      auth.user = null
      auth.error = errorMessage(error)
      auth.ready = true
      return
    }

    auth.user = data
    auth.error = ''
    auth.ready = true
  })()

  try {
    await pendingLoad
  } finally {
    pendingLoad = null
  }
}

export async function login(email: string, password: string): Promise<void> {
  auth.busy = true
  auth.error = ''

  try {
    const { data, error } = await api.POST('/auth/login', {
      body: { email, password },
    })
    if (error || !data) {
      auth.error = errorMessage(error)
      return
    }

    setToken(data.access_token)
    await loadMe()
  } finally {
    auth.busy = false
  }
}

export async function register(email: string, password: string, nickname: string): Promise<void> {
  auth.busy = true
  auth.error = ''

  try {
    const { data, error } = await api.POST('/auth/register', {
      body: { email, password, nickname },
    })
    if (error || !data) {
      auth.error = errorMessage(error)
      return
    }

    setToken(data.access_token)
    await loadMe()
  } finally {
    auth.busy = false
  }
}

export function logout(): void {
  setToken('')
  auth.user = null
  auth.error = ''
}

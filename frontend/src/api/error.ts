export function errorMessage(error: unknown): string {
  if (!error) {
    return 'Не удалось выполнить запрос'
  }

  if (typeof error === 'string') {
    return error
  }

  if (typeof error === 'object') {
    const value = error as { detail?: unknown; message?: unknown; title?: unknown }
    if (typeof value.detail === 'string' && value.detail) {
      return value.detail
    }
    if (typeof value.message === 'string' && value.message) {
      return value.message
    }
    if (typeof value.title === 'string' && value.title) {
      return value.title
    }
  }

  return 'Не удалось выполнить запрос'
}

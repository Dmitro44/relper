export const placeCategories = ['food', 'transport', 'healthcare', 'services', 'hobby', 'shopping']

export const placeCategoryLabels: Record<string, string> = {
  food: 'Еда',
  transport: 'Транспорт',
  healthcare: 'Медицина',
  services: 'Сервисы',
  hobby: 'Хобби',
  shopping: 'Покупки',
}

export const tipCategories = ['documents', 'banking', 'housing', 'integration', 'daily_life']

export const tipCategoryLabels: Record<string, string> = {
  documents: 'Документы',
  banking: 'Банки',
  housing: 'Жильё',
  integration: 'Адаптация',
  daily_life: 'Быт',
}

export function label(map: Record<string, string>, value: string): string {
  return map[value] ?? value
}

export function formatDate(value: string): string {
  return new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

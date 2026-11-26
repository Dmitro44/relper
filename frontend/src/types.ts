import type { components } from './api/schema'

export type Location = components['schemas']['LocationItem']
export type Place = components['schemas']['PlaceWithRatingItem']
export type Comment = components['schemas']['CommentItem']
export type Tip = components['schemas']['TipItem']
export type LikeSummary = components['schemas']['LikeSummary']
export type AuditEvent = components['schemas']['AuditItem']
export type Me = components['schemas']['MeBody']
export type Suggestion = components['schemas']['SuggestionItem']
export type SuggestionStatus = 'pending' | 'approved' | 'rejected'

export interface TipInput {
  title: string
  content: string
  category: string
}

export interface SuggestionInput {
  name: string
  category: string
  address: string
  description: string
  coordinates: { lat: number; lng: number }
}

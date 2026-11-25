import type { components } from './api/schema'

export type Location = components['schemas']['LocationItem']
export type Place = components['schemas']['PlaceWithRatingItem']
export type Comment = components['schemas']['CommentItem']
export type Tip = components['schemas']['TipItem']
export type AuditEvent = components['schemas']['AuditItem']
export type Me = components['schemas']['MeBody']

export interface TipInput {
  title: string
  content: string
  category: string
}

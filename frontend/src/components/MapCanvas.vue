<script setup lang="ts">
import * as maplibregl from 'maplibre-gl'
import maplibreWorkerUrl from 'maplibre-gl/dist/maplibre-gl-worker.mjs?url'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { Place } from '../types'

maplibregl.setWorkerUrl(maplibreWorkerUrl)

const props = defineProps<{
  places: Place[]
  selectedId: string | null
  center: { lat: number; lng: number }
}>()

const emit = defineEmits<{
  select: [id: string]
  close: []
  anchor: [point: { x: number; y: number } | null]
}>()

const container = ref<HTMLDivElement | null>(null)
let map: maplibregl.Map | null = null
let markers: maplibregl.Marker[] = []

const mapStyle = 'https://tiles.openfreemap.org/styles/bright'

function emitAnchor(): void {
  if (!map || !props.selectedId) {
    emit('anchor', null)
    return
  }

  const place = props.places.find((item) => item.id === props.selectedId)
  if (!place) {
    emit('anchor', null)
    return
  }

  const point = map.project([place.coordinates.lng, place.coordinates.lat])
  emit('anchor', { x: point.x, y: point.y })
}

function renderMarkers(): void {
  if (!map) {
    return
  }

  markers.forEach((marker) => marker.remove())
  markers = props.places.map((place) => {
    const element = document.createElement('button')
    element.type = 'button'
    element.className = place.id === props.selectedId ? 'map-pin is-selected' : 'map-pin'
    element.setAttribute('aria-label', place.name)
    element.addEventListener('click', (event) => {
      event.stopPropagation()
      emit('select', place.id)
    })

    const dot = document.createElement('span')
    dot.className = 'map-pin-dot'
    element.appendChild(dot)

    return new maplibregl.Marker({ element, anchor: 'center' })
      .setLngLat([place.coordinates.lng, place.coordinates.lat])
      .addTo(map!)
  })
}

function focusSelected(): void {
  if (!map || !props.selectedId) {
    return
  }

  const place = props.places.find((item) => item.id === props.selectedId)
  if (!place) {
    return
  }

  map.flyTo({
    center: [place.coordinates.lng, place.coordinates.lat],
    zoom: Math.max(13, map.getZoom()),
  })
}

function resize(): void {
  map?.resize()
  emitAnchor()
}

onMounted(() => {
  if (!container.value) {
    return
  }

  map = new maplibregl.Map({
    container: container.value,
    style: mapStyle,
    center: [props.center.lng, props.center.lat],
    zoom: 11,
  })

  map.addControl(new maplibregl.NavigationControl({ showCompass: false }), 'top-right')
  map.on('load', () => {
    renderMarkers()
    focusSelected()
    map?.resize()
    emitAnchor()
  })
  map.on('move', emitAnchor)
  map.on('zoom', emitAnchor)
  map.on('resize', emitAnchor)
  map.on('click', () => emit('close'))

  renderMarkers()
  requestAnimationFrame(() => map?.resize())
  window.addEventListener('resize', resize)
})

watch(
  () => props.places,
  () => {
    renderMarkers()
    requestAnimationFrame(emitAnchor)
  },
  { deep: true },
)

watch(
  () => props.selectedId,
  () => {
    renderMarkers()
    focusSelected()
    requestAnimationFrame(emitAnchor)
  },
)

watch(
  () => props.center,
  (center) => {
    map?.flyTo({ center: [center.lng, center.lat], zoom: 11 })
  },
  { deep: true },
)

onBeforeUnmount(() => {
  window.removeEventListener('resize', resize)
  markers.forEach((marker) => marker.remove())
  markers = []
  map?.remove()
  map = null
})
</script>

<template>
  <div ref="container" class="map-canvas" aria-label="Карта мест"></div>
</template>

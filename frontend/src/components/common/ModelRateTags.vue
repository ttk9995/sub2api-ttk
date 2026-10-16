<template>
  <div
    v-if="items.length > 0"
    class="flex min-w-0 max-w-full flex-wrap gap-1"
    data-testid="model-rate-tags"
  >
    <span
      v-for="item in items"
      :key="item.model"
      class="inline-flex max-w-full items-center gap-1 rounded-md border border-gray-200 bg-white px-1.5 py-0.5 text-[11px] font-medium leading-4 text-gray-600 shadow-sm dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300"
      :title="item.model + ' ' + formatRate(item.rate) + 'x'"
    >
      <span class="truncate">{{ item.model }}</span>
      <span class="shrink-0 font-semibold text-primary-600 dark:text-primary-400">{{ formatRate(item.rate) }}x</span>
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { GroupModelRate } from '@/types'

interface Props {
  modelRates?: GroupModelRate[]
  rateMultiplier?: number
  userRateMultiplier?: number | null
}

const props = withDefaults(defineProps<Props>(), {
  modelRates: () => [],
  userRateMultiplier: null
})

const items = computed(() => {
  const groupRate = Number(props.rateMultiplier)
  const userRate = props.userRateMultiplier ?? groupRate
  const hasUserRateAdjustment = Number.isFinite(groupRate) && groupRate > 0 && userRate !== groupRate

  return props.modelRates
    .map((item) => {
      const configuredRate = Number(item.rate_multiplier)
      if (!item.model || !Number.isFinite(configuredRate)) return null
      const effectiveRate = hasUserRateAdjustment ? configuredRate * userRate / groupRate : configuredRate
      if (!Number.isFinite(effectiveRate) || Math.abs(effectiveRate - 1) < 0.000001) return null
      return { model: item.model, rate: effectiveRate }
    })
    .filter((item): item is { model: string; rate: number } => item !== null)
})

const formatRate = (value: number): string => {
  const rounded = Number(value.toFixed(4))
  return Object.is(rounded, -0) ? '0' : String(rounded)
}
</script>

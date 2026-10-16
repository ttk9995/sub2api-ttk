<template>
  <div class="flex min-w-0 flex-1 items-start justify-between gap-3">
    <!-- Left: name + description -->
    <div
      class="flex min-w-0 flex-1 flex-col items-start"
      :title="description || undefined"
    >
      <!-- Row 1: platform badge (name bold) -->
      <GroupBadge
        :name="name"
        :platform="platform"
        :subscription-type="subscriptionType"
        :show-rate="false"
        class="groupOptionItemBadge"
      />
      <!-- Row 2: automatic model-rate tags; keep legacy description as fallback. -->
      <div
        v-if="modelRateItems.length > 0"
        class="mt-1.5 flex w-full flex-wrap gap-1"
      >
        <span
          v-for="item in modelRateItems"
          :key="item.model"
          class="inline-flex max-w-full items-center gap-1 rounded-md border border-gray-200 bg-white px-1.5 py-0.5 text-[11px] font-medium leading-4 text-gray-600 shadow-sm dark:border-dark-600 dark:bg-dark-800 dark:text-dark-300"
          :title="`${item.model} ${formatRate(item.rate)}x`"
        >
          <span class="truncate">{{ item.model }}</span>
          <span class="shrink-0 font-semibold text-primary-600 dark:text-primary-400">{{ formatRate(item.rate) }}x</span>
        </span>
      </div>
      <span
        v-else-if="description"
        class="mt-1.5 w-full whitespace-pre-line [overflow-wrap:anywhere] text-left text-xs leading-relaxed text-gray-500 dark:text-gray-400 line-clamp-3"
      >
        {{ description }}
      </span>
    </div>

    <!-- Right: rate pill + checkmark (vertically centered to first row) -->
    <div class="flex shrink-0 items-center gap-2 pt-0.5">
      <div class="flex shrink-0 flex-col items-end gap-1">
        <!-- Rate pill (platform color) -->
        <span v-if="showRatePill" :class="['inline-flex items-center whitespace-nowrap rounded-full px-3 py-1 text-xs font-semibold', ratePillClass]">
          <template v-if="hasCustomRate">
            <span class="mr-1 line-through opacity-50">{{ rateMultiplier }}x</span>
            <span class="font-bold">{{ userRateMultiplier }}x</span>
          </template>
          <template v-else>
            {{ rateMultiplier }}x {{ t('admin.groups.rateLabel') }}
          </template>
        </span>
        <span
          v-if="hasPeakRate"
          class="inline-flex items-center whitespace-nowrap rounded-full bg-amber-50 px-3 py-1 text-xs font-semibold text-amber-700 dark:bg-amber-900/20 dark:text-amber-300"
          :title="peakRateTitle"
        >
          {{ peakRateText }}
        </span>
      </div>
      <!-- Checkmark -->
      <svg
        v-if="showCheckmark && selected"
        class="h-4 w-4 shrink-0 text-primary-600 dark:text-primary-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        stroke-width="2"
      >
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import GroupBadge from './GroupBadge.vue'
import type { GroupModelRate, SubscriptionType, GroupPlatform } from '@/types'
import { useAppStore } from '@/stores/app'
import { formatPeakRateWindow, serverTimezoneLabel } from '@/utils/peak-rate'

const { t } = useI18n()

interface Props {
  name: string
  platform: GroupPlatform
  subscriptionType?: SubscriptionType
  rateMultiplier?: number
  userRateMultiplier?: number | null
  peakRateEnabled?: boolean
  peakStart?: string
  peakEnd?: string
  peakRateMultiplier?: number
  description?: string | null
  modelRates?: GroupModelRate[]
  selected?: boolean
  showCheckmark?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  subscriptionType: 'standard',
  selected: false,
  showCheckmark: true,
  userRateMultiplier: null,
  peakRateEnabled: false,
  modelRates: () => []
})

// Whether user has a custom rate different from default
const hasCustomRate = computed(() => {
  return (
    props.userRateMultiplier !== null &&
    props.userRateMultiplier !== undefined &&
    props.rateMultiplier !== undefined &&
    props.userRateMultiplier !== props.rateMultiplier
  )
})

// The default 1x group rate adds no information; keep the pill for overrides.
const showRatePill = computed(() => {
  if (props.rateMultiplier === undefined) return false
  if (hasCustomRate.value) return true
  return Math.abs(Number(props.rateMultiplier) - 1) >= 0.000001
})

const appStore = useAppStore()

const hasPeakRate = computed(() => {
  return Boolean(props.peakRateEnabled && props.peakStart && props.peakEnd)
})

const peakRateText = computed(() => {
  return formatPeakRateWindow(
    {
      peak_rate_enabled: props.peakRateEnabled,
      peak_start: props.peakStart,
      peak_end: props.peakEnd,
      peak_rate_multiplier: props.peakRateMultiplier
    },
    serverTimezoneLabel(appStore.cachedPublicSettings?.server_utc_offset)
  )
})

const peakRateTitle = computed(() => {
  return t('common.peakRateTooltip', { window: peakRateText.value })
})

const modelRateItems = computed(() => {
  const groupRate = Number(props.rateMultiplier)
  const userRate = props.userRateMultiplier ?? groupRate
  const hasUserRateAdjustment = Number.isFinite(groupRate) && groupRate > 0 && userRate !== groupRate

  return props.modelRates
    .map((item) => {
      const configuredRate = Number(item.rate_multiplier)
      if (!item.model || !Number.isFinite(configuredRate)) return null
      const effectiveRate = hasUserRateAdjustment
        ? configuredRate * userRate / groupRate
        : configuredRate
      if (!Number.isFinite(effectiveRate) || Math.abs(effectiveRate - 1) < 0.000001) return null
      return { model: item.model, rate: effectiveRate }
    })
    .filter((item): item is { model: string; rate: number } => item !== null)
})

const formatRate = (value: number): string => {
  const rounded = Number(value.toFixed(4))
  return Object.is(rounded, -0) ? '0' : String(rounded)
}

// Rate pill color matches platform badge color
const ratePillClass = computed(() => {
  switch (props.platform) {
    case 'anthropic':
      return 'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-400'
    case 'openai':
      return 'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400'
    case 'gemini':
      return 'bg-sky-50 text-sky-700 dark:bg-sky-900/20 dark:text-sky-400'
    default: // antigravity and others
      return 'bg-violet-50 text-violet-700 dark:bg-violet-900/20 dark:text-violet-400'
  }
})
</script>

<style scoped>
/* Bold the group name inside GroupBadge when used in dropdown option */
.groupOptionItemBadge :deep(span.truncate) {
  font-weight: 600;
}
</style>

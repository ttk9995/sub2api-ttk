<template>
  <div class="mt-3 rounded-lg border border-gray-200 p-3 dark:border-dark-400">
    <div class="flex items-start justify-between gap-3">
      <div class="min-w-0">
        <p class="text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('admin.groups.modelPricing.gpt61SolThresholdTitle') }}
        </p>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          {{ t('admin.groups.modelPricing.gpt61SolThresholdHint') }}
        </p>
      </div>
      <Toggle
        :model-value="enabled"
        :aria-label="t('admin.groups.modelPricing.gpt61SolThresholdTitle')"
        @update:model-value="emit('update:enabled', $event)"
      />
    </div>
    <label class="mt-3 block">
      <span class="input-label">{{ t('admin.groups.modelPricing.contextThresholdTokens') }}</span>
      <input
        :value="threshold"
        :disabled="!enabled"
        :required="enabled"
        type="number"
        min="1"
        max="2147483647"
        step="1"
        class="input disabled:cursor-not-allowed disabled:opacity-50"
        @input="updateThreshold"
      />
    </label>
    <p class="mt-2 text-xs text-gray-500 dark:text-gray-400">
      {{ t('admin.groups.modelPricing.gpt61SolThresholdPolicy') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Toggle from '@/components/common/Toggle.vue'

defineProps<{ enabled: boolean; threshold: number }>()
const emit = defineEmits<{
  'update:enabled': [value: boolean]
  'update:threshold': [value: number]
}>()
const { t } = useI18n()

function updateThreshold(event: Event): void {
  emit('update:threshold', (event.target as HTMLInputElement).valueAsNumber)
}
</script>

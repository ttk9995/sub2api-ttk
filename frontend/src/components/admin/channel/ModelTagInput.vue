<template>
  <div>
    <!-- Tags display -->
    <div class="relative" @focusout="handleFocusOut">
      <div class="flex min-h-[2.5rem] flex-wrap gap-1.5 rounded-lg border border-gray-200 bg-white p-2 dark:border-dark-600 dark:bg-dark-800">
      <span
        v-for="(model, idx) in models"
        :key="idx"
        class="inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-sm"
        :class="getPlatformTagClass(props.platform || '')"
      >
        {{ model }}
        <span
          v-if="isCandidateMode && !candidateSet.has(model)"
          class="text-[10px] text-amber-700 dark:text-amber-300"
          :title="t('admin.channels.form.modelCandidateUnavailable')"
        >
          {{ t('admin.channels.form.modelCandidateUnavailable') }}
        </span>
        <button
          type="button"
          @click="removeModel(idx)"
          class="ml-0.5 rounded-full p-0.5 hover:bg-primary-200 dark:hover:bg-primary-800"
        >
          <Icon name="x" size="xs" />
        </button>
      </span>
      <input
        ref="inputRef"
        v-model="inputValue"
        type="text"
        class="flex-1 min-w-[120px] border-none bg-transparent text-sm outline-none placeholder:text-gray-400 dark:text-white"
        :placeholder="inputPlaceholder"
        :role="isCandidateMode ? 'combobox' : undefined"
        :aria-expanded="isCandidateMode ? isCandidateMenuOpen : undefined"
        :aria-controls="isCandidateMode ? candidateListID : undefined"
        :aria-autocomplete="isCandidateMode ? 'list' : undefined"
        @focus="isCandidateMenuOpen = isCandidateMode"
        @keydown.enter="handleEnter"
        @keydown.tab="handleTab"
        @keydown.backspace="handleBackspace"
        @paste="handlePaste"
        @blur="handleInputBlur"
      />
      </div>
      <div
        v-if="isCandidateMode && isCandidateMenuOpen"
        :id="candidateListID"
        role="listbox"
        aria-multiselectable="true"
        class="absolute left-0 right-0 top-full z-30 mt-1 max-h-60 overflow-y-auto rounded-md border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-500 dark:bg-dark-700"
      >
      <div v-if="candidatesLoading" class="px-3 py-2 text-xs text-gray-500 dark:text-gray-400">
        {{ t('admin.channels.form.modelCandidatesLoading') }}
      </div>
      <div v-else-if="candidatesError" class="flex items-center justify-between gap-3 px-3 py-2 text-xs text-red-600 dark:text-red-400">
        <span>{{ t('admin.channels.form.modelCandidatesError') }}</span>
        <button
          type="button"
          class="inline-flex shrink-0 items-center gap-1 rounded px-2 py-1 font-medium text-primary-600 hover:bg-primary-50 dark:text-primary-300 dark:hover:bg-dark-600"
          @mousedown.prevent
          @click="emit('retry')"
        >
          <Icon name="refresh" size="xs" />
          {{ t('admin.channels.form.modelCandidatesRetry') }}
        </button>
      </div>
      <div v-else-if="filteredCandidates.length === 0" class="px-3 py-2 text-xs text-gray-500 dark:text-gray-400">
        {{ inputValue.trim() ? t('admin.channels.form.modelCandidatesNoMatch') : t('admin.channels.form.modelCandidatesEmpty') }}
      </div>
      <button
        v-for="model in filteredCandidates"
        :key="model"
        type="button"
        role="option"
        :aria-selected="models.includes(model)"
        :data-testid="`model-candidate-${model}`"
        class="flex w-full items-center gap-2 px-3 py-1.5 text-left text-sm text-gray-700 hover:bg-gray-50 dark:text-gray-200 dark:hover:bg-dark-600"
        @mousedown.prevent
        @click="toggleCandidate(model)"
      >
        <span
          class="flex h-4 w-4 shrink-0 items-center justify-center rounded border"
          :class="models.includes(model) ? 'border-primary-600 bg-primary-600 text-white' : 'border-gray-300 dark:border-dark-400'"
        >
          <Icon v-if="models.includes(model)" name="check" size="xs" />
        </span>
        <span class="min-w-0 truncate">{{ model }}</span>
      </button>
      </div>
    </div>
    <p class="mt-1 text-xs text-gray-400">
      {{ isCandidateMode ? t('admin.channels.form.modelCandidatesHint') : t('admin.channels.form.modelInputHint', 'Press Enter to add, supports paste for batch import.') }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, useId } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { getPlatformTagClass } from './types'

const { t } = useI18n()

const props = defineProps<{
  models: string[]
  placeholder?: string
  platform?: string
  candidates?: string[]
  candidatesLoading?: boolean
  candidatesError?: boolean
}>()

const emit = defineEmits<{
  'update:models': [models: string[]]
  retry: []
}>()

const inputValue = ref('')
const inputRef = ref<HTMLInputElement>()
const isCandidateMenuOpen = ref(false)
const candidateListID = `model-candidate-list-${useId()}`
const isCandidateMode = computed(() => props.candidates !== undefined)
const candidateSet = computed(() => new Set(props.candidates || []))
const inputPlaceholder = computed(() => {
  if (props.models.length > 0) return ''
  return isCandidateMode.value
    ? t('admin.channels.form.modelCandidatesPlaceholder')
    : props.placeholder || ''
})
const filteredCandidates = computed(() => {
  const query = inputValue.value.trim().toLocaleLowerCase()
  return (props.candidates || []).filter(model => !query || model.toLocaleLowerCase().includes(query))
})

function handleFocusOut(event: FocusEvent) {
  const nextTarget = event.relatedTarget
  if (nextTarget instanceof Node && event.currentTarget instanceof HTMLElement && event.currentTarget.contains(nextTarget)) {
    return
  }
  isCandidateMenuOpen.value = false
}

function handleInputBlur() {
  if (!isCandidateMode.value) addModel()
}

function addModel() {
  const val = inputValue.value.trim()
  if (!val) return
  if (!props.models.includes(val)) {
    emit('update:models', [...props.models, val])
  }
  inputValue.value = ''
}

function handleEnter(event: KeyboardEvent) {
  if (event.isComposing) return
  if (isCandidateMode.value) {
    event.preventDefault()
    const firstCandidate = filteredCandidates.value[0]
    if (!firstCandidate) return
    if (!props.models.includes(firstCandidate)) toggleCandidate(firstCandidate)
    return
  }
  event.preventDefault()
  addModel()
}

function handleTab(event: KeyboardEvent) {
  if (isCandidateMode.value) return
  if (event.isComposing || !inputValue.value.trim()) return
  event.preventDefault()
  addModel()
}

function removeModel(idx: number) {
  const newModels = [...props.models]
  newModels.splice(idx, 1)
  emit('update:models', newModels)
}

function handleBackspace(event: KeyboardEvent) {
  if (event.isComposing) return
  if (inputValue.value === '' && props.models.length > 0) {
    removeModel(props.models.length - 1)
  }
}

function handlePaste(e: ClipboardEvent) {
  if (isCandidateMode.value) return
  e.preventDefault()
  const text = e.clipboardData?.getData('text') || ''
  const items = text.split(/[,\n;]+/).map(s => s.trim()).filter(Boolean)
  if (items.length === 0) return
  const unique = [...new Set([...props.models, ...items])]
  emit('update:models', unique)
  inputValue.value = ''
}

function toggleCandidate(model: string) {
  const nextModels = props.models.includes(model)
    ? props.models.filter(selected => selected !== model)
    : [...props.models, model]
  emit('update:models', nextModels)
}
</script>

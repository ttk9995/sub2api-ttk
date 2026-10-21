import type { GroupPlatform } from '@/types'

export function supportsGroupLongContextPricing(platform: GroupPlatform): boolean {
  return platform === 'openai' || platform === 'composite'
}

export function isValidGroupContextThreshold(value: number): boolean {
  return Number.isInteger(value) && value >= 1 && value <= 2147483647
}

export function groupLongContextPricingPayload(form: {
  platform: GroupPlatform
  gpt61_sol_long_context_enabled: boolean
  gpt61_sol_long_context_threshold: number
}) {
  return {
    gpt61_sol_long_context_enabled:
      supportsGroupLongContextPricing(form.platform) && form.gpt61_sol_long_context_enabled,
    gpt61_sol_long_context_threshold: isValidGroupContextThreshold(form.gpt61_sol_long_context_threshold)
      ? form.gpt61_sol_long_context_threshold
      : 500000,
  }
}

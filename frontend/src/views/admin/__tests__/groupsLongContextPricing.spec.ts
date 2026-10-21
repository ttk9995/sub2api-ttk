import { describe, expect, it } from 'vitest'
import { groupLongContextPricingPayload, isValidGroupContextThreshold } from '../groupsLongContextPricing'

describe('group context threshold validation', () => {
  it('rejects empty, fractional, nonpositive, and overflowing token counts', () => {
    for (const value of [NaN, Infinity, 0, -1, 500000.5, 2147483648]) {
      expect(isValidGroupContextThreshold(value)).toBe(false)
    }
    for (const value of [1, 500000, 640000, 2147483647]) {
      expect(isValidGroupContextThreshold(value)).toBe(true)
    }
  })

  it('preserves the value when disabled and clears enabled state on unsupported platforms', () => {
    expect(groupLongContextPricingPayload({
      platform: 'openai', gpt61_sol_long_context_enabled: false, gpt61_sol_long_context_threshold: 640000,
    })).toEqual({ gpt61_sol_long_context_enabled: false, gpt61_sol_long_context_threshold: 640000 })
    expect(groupLongContextPricingPayload({
      platform: 'anthropic', gpt61_sol_long_context_enabled: true, gpt61_sol_long_context_threshold: 640000,
    }).gpt61_sol_long_context_enabled).toBe(false)
  })
})

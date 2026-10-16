import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import GroupBadge from '../GroupBadge.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key }),
  }
})

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ cachedPublicSettings: null }),
}))

describe('GroupBadge default rate visibility', () => {
  const mountBadge = (rateMultiplier: number) => mount(GroupBadge, {
    props: {
      name: 'Example group',
      platform: 'openai',
      rateMultiplier,
      hideDefaultRate: true,
    },
    global: {
      stubs: {
        PlatformIcon: true,
      },
    },
  })

  it('hides the default 1x label when requested', () => {
    expect(mountBadge(1).text()).not.toContain('1x')
  })

  it('keeps a non-default group rate visible', () => {
    expect(mountBadge(1.5).text()).toContain('1.5x')
  })
})

import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'

import GroupOptionItem from '../GroupOptionItem.vue'

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

describe('GroupOptionItem description layout', () => {
  it('applies multiline and overflow-safe text styles', () => {
    const description = 'First section\nvery-long-unbroken-description-value-that-must-not-overflow'
    const wrapper = mount(GroupOptionItem, {
      props: {
        name: 'Example group',
        platform: 'openai',
        description,
      },
      global: {
        stubs: {
          GroupBadge: true,
        },
      },
    })

    const descriptionElement = wrapper
      .findAll('span')
      .find((element) => element.text() === description)

    expect(descriptionElement).toBeDefined()
    expect(descriptionElement?.classes()).toContain('whitespace-pre-line')
    expect(descriptionElement?.classes()).toContain('[overflow-wrap:anywhere]')
    expect(descriptionElement?.classes()).toContain('line-clamp-3')
    expect(wrapper.find('[title]').attributes('title')).toBe(description)
  })

  it('shows explicit model rates and hides the default 1x tag', () => {
    const wrapper = mount(GroupOptionItem, {
      props: {
        name: 'Example group',
        platform: 'openai',
        rateMultiplier: 1,
        modelRates: [
          { model: 'gpt-5.6-luna', rate_multiplier: 0.3 },
          { model: 'gpt-5.6-mini', rate_multiplier: 1 },
        ],
      },
      global: {
        stubs: {
          GroupBadge: true,
        },
      },
    })

    expect(wrapper.text()).toContain('gpt-5.6-luna')
    expect(wrapper.text()).toContain('0.3x')
    expect(wrapper.text()).not.toContain('gpt-5.6-mini')
    expect(wrapper.text()).not.toContain('1x 倍率')
  })

  it('keeps the rate pill when a group has a non-default rate', () => {
    const wrapper = mount(GroupOptionItem, {
      props: {
        name: 'Premium group',
        platform: 'openai',
        rateMultiplier: 1.5,
      },
      global: {
        stubs: {
          GroupBadge: true,
        },
      },
    })

    expect(wrapper.text()).toContain('1.5x admin.groups.rateLabel')
  })
})

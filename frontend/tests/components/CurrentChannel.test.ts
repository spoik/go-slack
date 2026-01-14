// @vitest-environment happy-dom

import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CurrentChannel from '@/components/CurrentChannel.vue'

describe('CurrentChannel component', () => {
    it('shows a message when no channel is selected', () => {
        const wrapper = mount(CurrentChannel, {
            props: {
                channel: undefined
            }
        })

        expect(wrapper.get('[data-test="empty message"]').text()).toContain('Please select a channel.')
        expect(wrapper.find('h1').exists()).toBe(false)
    })

    it('shows the channel name when a channel is selected', () => {
        const channel = { id: '1', name: 'general' }
        const wrapper = mount(CurrentChannel, {
            props: {
                channel
            }
        })

        expect(wrapper.get('h1').text()).toBe(channel.name)
        expect(wrapper.find('[data-test="empty message"]').exists()).toBe(false)
    })
})

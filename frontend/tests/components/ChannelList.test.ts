// @vitest-environment happy-dom

import { nextTick } from 'vue'
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { getChannels, type Channel } from '@/utils/channel-service'
import ChannelList from '@/components/ChannelList.vue'

vi.mock('@/utils/channel-service')

describe('ChannelList component', () => {
	function mockGetChannels(channels: Channel[]) {
		vi.mocked(getChannels).mockResolvedValue(channels)
	}

	it('shows a loading message while the channels are loading', async () => {
		mockGetChannels([{id: 1, name: "test"}])

		const wrapper = mount(ChannelList)

		expect(wrapper.find('[data-test="loading"]').exists()).toBe(true)

		await nextTick()
		await nextTick()

		expect(wrapper.find('[data-test="loading"]').exists()).toBe(false)
		expect(wrapper.find('[data-test="error"]').exists()).toBe(false)
	})

	it('loads channels from the api and shows them', async () => {
		mockGetChannels([
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
			{ id: '3', name: 'other' },
		])

		const wrapper = mount(ChannelList)
		await nextTick()
		await nextTick()

		expect(wrapper.findAll('li')).toHaveLength(3)
	})

	it('shows an error if the channels failed to load', async () => {
		vi.mocked(getChannels).mockRejectedValue()

		const wrapper = mount(ChannelList)
		expect(wrapper.find('[data-test="error"]').exists()).toBe(false)

		await nextTick()
		await nextTick()

		expect(wrapper.find('[data-test="error"]').exists()).toBe(true)
	})
})

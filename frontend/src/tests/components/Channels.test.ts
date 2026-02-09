// @vitest-environment happy-dom

import { describe, it, expect, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { nextTick } from 'vue'

import Channels from '@/components/Channels.vue'
import { getChannels, type Channel } from '@/utils/channel-service'
import ChannelList from '@/components/ChannelList.vue'
import CurrentChannel from '@/components/CurrentChannel.vue'

vi.mock('@/utils/channel-service')

describe('Channels component', () => {
	it('forwards the selected channel from ChannelList to CurrentChannel and back to ChannelList', async () => {
		const testChannel: Channel = { id: '1', name: 'general' }

		const wrapper = shallowMount(Channels)
		const currentChannl = wrapper.findComponent(CurrentChannel)
		const channelList = wrapper.findComponent(ChannelList)

		expect(currentChannl.props('channel')).toEqual(undefined)
		expect(channelList.props('selectedChannel')).toEqual(undefined)

		channelList.vm.$emit('channelSelected', testChannel)
		await nextTick()

		expect(currentChannl.props('channel')).toEqual(testChannel)
		expect(channelList.props('selectedChannel')).toEqual(testChannel)
	})
})

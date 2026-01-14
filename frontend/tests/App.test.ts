// @vitest-environment happy-dom

import { describe, it, expect, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import App from '@/App.vue'
import { getChannels, type Channel } from '@/utils/channel-service'
import  ChannelList from '@/components/ChannelList.vue'
import  CurrentChannel from '@/components/CurrentChannel.vue'

vi.mock('@/utils/channel-service')

describe('App component', () => {
    function mockGetChannels(channels: Channel[]) {
	vi.mocked(getChannels).mockResolvedValue(channels)
    }

    it('forwards the selected channel from ChannelList to CurrentChannel', async () => {
        const testChannels = [
            { id: '1', name: 'general' },
            { id: '2', name: 'temp' },
        ]
        mockGetChannels(testChannels)

        const wrapper = shallowMount(App)

        // Simulate ChannelList emitting a channelSelected event
        const channelList = wrapper.findComponent(ChannelList)
        await channelList.vm.$emit('channelSelected', testChannels[1])
        
        // Assert that CurrentChannel is passed the channel emitted in the channelSelected event
        const currentChannl = wrapper.findComponent(CurrentChannel)
        expect(currentChannl.props('channel')).toEqual(testChannels[1])
    })
})

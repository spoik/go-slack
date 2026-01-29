// @vitest-environment happy-dom

import { describe, it, expect, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import App from '@/App.vue'
import { getChannels, type Channel } from '@/utils/channel-service'
import  ChannelList from '@/components/ChannelList.vue'
import  CurrentChannel from '@/components/CurrentChannel.vue'
import { nextTick } from 'vue'

vi.mock('@/utils/channel-service')

describe('App component', () => {
    function mockGetChannels(channels: Channel[]) {
	vi.mocked(getChannels).mockResolvedValue(channels)
    }

    it('forwards the selected channel from ChannelList to CurrentChannel and back to ChannelList', async () => {
        const testChannels = [
            { id: '1', name: 'general' },
            { id: '2', name: 'temp' },
        ]
        mockGetChannels(testChannels)

        const wrapper = shallowMount(App)
        const currentChannl = wrapper.findComponent(CurrentChannel)
        const channelList = wrapper.findComponent(ChannelList)

        expect(currentChannl.props('channel')).toEqual(undefined)
        expect(channelList.props('selectedChannel')).toEqual(undefined)

        channelList.vm.$emit('channelSelected', testChannels[1])
        await nextTick()
        
        expect(currentChannl.props('channel')).toEqual(testChannels[1])
        expect(channelList.props('selectedChannel')).toEqual(testChannels[1])
    })
})

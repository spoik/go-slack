// @vitest-environment happy-dom

import { nextTick } from 'vue'
import { describe, it, expect, vi } from 'vitest'
import { mount, shallowMount, VueWrapper, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory, type Router } from 'vue-router'

import { routes } from '@/routes'
import { getChannels, type Channel } from '@/utils/channel-service'
import ChannelList from '@/components/ChannelList.vue'
import CreateChannel from '@/components/CreateChannel.vue'

type WrappedChannelList = VueWrapper<InstanceType<typeof ChannelList>>

vi.mock('@/utils/channel-service')

describe('ChannelList component', () => {
	let router: Router

	beforeEach(() => {
		router = createRouter({
			history: createMemoryHistory(),
			routes
		})
	})

	function mockGetChannels(channels: Channel[]) {
		vi.mocked(getChannels).mockResolvedValue(channels)
	}

	async function initWrapper(route: string = "/", selectedChannel: Channel | null = null): Promise<WrappedChannelList> {
		var props = { selectedChannel: undefined as Channel | undefined }

		if (selectedChannel != null) {
			props.selectedChannel = selectedChannel
		}

		router.push(route)
		await router.isReady()

		return mount(
			ChannelList,
			{
				props,
				global: {
					plugins: [router]
				}
			}
		)
	}

	it('shows a loading message while the channels are loading', async () => {
		mockGetChannels([{ id: "1", name: "test" }])

		const wrapper = await initWrapper()

		expect(wrapper.find('[data-test="loading"]').exists()).toBe(true)

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

		const wrapper = await initWrapper()
		await nextTick()

		expect(wrapper.findAll('li')).toHaveLength(3)
	})

	it('shows an error if the channels failed to load', async () => {
		vi.mocked(getChannels).mockRejectedValue(new Error("error"))

		const wrapper = await initWrapper()
		expect(wrapper.find('[data-test="error"]').exists()).toBe(false)

		await nextTick()

		expect(wrapper.find('[data-test="error"]').exists()).toBe(true)
		expect(wrapper.find('[data-test="error"]').text()).toEqual('Failed to load channels. Please reload to try again.')
	})

	it('changes the browser URL to the page for a specific channel when a channel is clicked', async () => {
		const testChannels: Channel[] = [
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
		]
		mockGetChannels(testChannels)

		const wrapper = await initWrapper()
		await flushPromises()

		expect(router.currentRoute.value.path).toEqual("/")

		await wrapper.findAll('a[data-test="channel"]')[0]?.trigger('click')
		await flushPromises()

		expect(router.currentRoute.value.path).toEqual(`/channels/${testChannels[0]?.id}`)
	})

	it('shows no active channel when selectedChannel prop is empty', async () => {
		const channels: Channel[] = [
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
		]
		mockGetChannels(channels)

		const wrapper = await initWrapper("/", undefined)

		expect(wrapper.findAll('[data-test="channel"][data-test-active]').length).toBe(0)
	})

	it('shows an active channel when selectedChannel prop is provided', async () => {
		const channels: Channel[] = [
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
		]
		mockGetChannels(channels)

		const wrapper = await initWrapper("/", channels[0])
		await nextTick()

		expect(wrapper.findAll('[data-test="channel"][data-test-active]').length).toBe(1)
		expect(wrapper.get('[data-test="channel"][data-test-active]').text()).toEqual(channels[0]?.name)
	})

	it('emits a channelSelected event when the url is for a specific channel', async () => {
		const channels: Channel[] = [
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
		]
		mockGetChannels(channels)

		router.push(`/channels/${channels[0]?.id}`)
		await router.isReady()

		const wrapper = await initWrapper()
		await nextTick()

		expect(wrapper.emitted('channelSelected')).toBeTruthy()
		expect(wrapper.emitted('channelSelected')?.[0]?.[0]).toEqual(channels[0])

		router.push(`/channels/${channels[1]?.id}`)
		await flushPromises()

		expect(wrapper.emitted('channelSelected')?.length).toEqual(2)
		expect(wrapper.emitted('channelSelected')?.[1]?.[0]).toEqual(channels[1])
	})

	describe('when the CreateChannel component emits a newly created channel', () => {
		let wrapper: VueWrapper<InstanceType<typeof ChannelList>>
		let existingChannels: Channel[]

		beforeEach(async () => {
			existingChannels = [
				{ id: "1", name: "A" },
				{ id: "2", name: "C" }
			]
			mockGetChannels(structuredClone(existingChannels))

			wrapper = mount(ChannelList, {
				global: {
					plugins: [router]
				}
			})

			await nextTick()
		})

		async function emitChannelCreated(channel: Channel) {
			const createChannelComponent = wrapper.findComponent(CreateChannel)
			createChannelComponent.vm.$emit('channelCreated', channel)
			await nextTick()
		}

		it('shows the new channel in the list of channels', async () => {
			expect(wrapper.findAll('[data-test="channel"]')).toHaveLength(2)

			const newChannel: Channel = { id: "3", name: "B" }
			await emitChannelCreated(newChannel);

			const channelElements = wrapper.findAll('[data-test="channel"]')
			expect(channelElements).toHaveLength(3)
			expect(channelElements[0]?.text()).toEqual(existingChannels[0]?.name)
			expect(channelElements[1]?.text()).toEqual(newChannel.name)
			expect(channelElements[2]?.text()).toEqual(existingChannels[1]?.name)
		})
	})
})

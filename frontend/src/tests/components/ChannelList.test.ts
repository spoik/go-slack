// @vitest-environment happy-dom

import { nextTick } from 'vue'
import { describe, it, expect, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import { getChannels, createChannel, type Channel } from '@/utils/channel-service'
import ChannelList from '@/components/ChannelList.vue'

vi.mock('@/utils/channel-service')

describe('ChannelList component', () => {
	function mockGetChannels(channels: Channel[]) {
		vi.mocked(getChannels).mockResolvedValue(channels)
	}

	it('shows a loading message while the channels are loading', async () => {
		mockGetChannels([{ id: "1", name: "test" }])

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
		vi.mocked(getChannels).mockRejectedValue(new Error("error"))

		const wrapper = mount(ChannelList)
		expect(wrapper.find('[data-test="error"]').exists()).toBe(false)

		await nextTick()
		await nextTick()

		expect(wrapper.find('[data-test="error"]').exists()).toBe(true)
		expect(wrapper.find('[data-test="error"]').text()).toEqual('Failed to load channels. Please reload to try again.')
	})

	it('emits channelSelected event with channel the when a channel is clicked', async () => {
		const testChannels = [
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
		]
		mockGetChannels(testChannels)

		const wrapper = mount(ChannelList)
		await nextTick()
		await nextTick()

		wrapper.findAll('[data-test="channel"]')[0]?.trigger('click')

		expect(wrapper.emitted('channelSelected')).toBeTruthy()

		const emittedChannel = wrapper.emitted('channelSelected')?.[0]?.[0]
		expect(emittedChannel).toEqual(testChannels[0])
	})

	it('shows no active channel when selectedChannel prop is empty', async () => {
		const channels: Channel[] = [
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
		]
		mockGetChannels(channels)

		const wrapper = mount(ChannelList, {
			props: {
				selectedChannel: undefined
			}
		})
		await nextTick()
		await nextTick()

		expect(wrapper.findAll('[data-test="channel"].active').length).toBe(0)
	})

	it('shows an active channel when selectedChannel prop is provided', async () => {
		const channels: Channel[] = [
			{ id: '1', name: 'general' },
			{ id: '2', name: 'temp' },
		]
		mockGetChannels(channels)

		const wrapper = mount(ChannelList, {
			props: {
				selectedChannel: channels[0]
			}
		})
		await nextTick()
		await nextTick()

		expect(wrapper.findAll('[data-test="channel"].active').length).toBe(1)
		expect(wrapper.get('[data-test="channel"].active').text()).toEqual(channels[0]?.name)
	})

	describe('creating a new channel', async () => {
		let wrapper: VueWrapper<InstanceType<typeof ChannelList>>
		let existingChannels: Channel[]

		beforeEach(async () => {
			existingChannels = [
				{ id: "1", name: "A" },
				{ id: "2", name: "C" }
			]
			mockGetChannels(structuredClone(existingChannels))
			wrapper = mount(ChannelList)
			await nextTick()
		})

		it('shows the create channel form when create channel button is clicked', async () => {
			expect(wrapper.find('[data-test="create channel form"]').exists()).toBe(false)
			expect(wrapper.find('[data-test="create channel button"]').exists()).toBe(true)

			wrapper.get('[data-test="create channel button"]').trigger('click')
			await nextTick()

			expect(wrapper.find('[data-test="create channel form"]').exists()).toBe(true)
			expect(wrapper.find('[data-test="create channel button"]').exists()).toBe(false)
		})

		describe('submitting the form', () => {
			beforeEach(async () => {
				vi.mocked(createChannel).mockResolvedValue({ id: "1", name: "Anything" })
				wrapper.get('[data-test="create channel button"]').trigger('click')
				await nextTick()
			})

			it('makes a request to create a new channel when the channel form is submitted', async () => {
				const channelName = "New channel"
				wrapper.get('[data-test="channel name input"]').setValue(channelName)
				wrapper.get('[data-test="create channel form"]').trigger("submit")
				expect(createChannel).toHaveBeenCalledWith(channelName)
			})

			describe('when the new channel was created successfully', () => {
				let newChannel: Channel
				beforeEach(() => {
					newChannel = { id: "1", name: "B Channel name" }
					vi.mocked(createChannel).mockResolvedValue(newChannel)
				})

				async function submitNewChannel() {
					wrapper.get('[data-test="channel name input"]').setValue(newChannel.name)
					wrapper.get('[data-test="create channel form"]').trigger("submit")
					await nextTick()
				}

				it('hides the new channel form', async () => {
					expect(wrapper.find('[data-test="create channel form"]').exists()).toBe(true)
					expect(wrapper.find('[data-test="create channel button"]').exists()).toBe(false)

					await submitNewChannel()

					expect(wrapper.find('[data-test="create channel form"]').exists()).toBe(false)
					expect(wrapper.find('[data-test="create channel button"]').exists()).toBe(true)
				})

				it('shows the new channel in the list of channels', async () => {
					expect(wrapper.findAll('[data-test="channel"]')).toHaveLength(2)

					await submitNewChannel()

					const channelElements = wrapper.findAll('[data-test="channel"]')
					expect(channelElements).toHaveLength(3)
					expect(channelElements[0]?.text()).toEqual(existingChannels[0]?.name)
					expect(channelElements[1]?.text()).toEqual(newChannel.name)
					expect(channelElements[2]?.text()).toEqual(existingChannels[1]?.name)
				})
			})
		})
	})
})

// @vitest-environment happy-dom

import CreateChannel from '@/components/CreateChannel.vue'
import { mount, VueWrapper } from '@vue/test-utils'
import { nextTick } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

import { createChannel, type Channel } from '@/utils/channel-service'

vi.mock('@/utils/channel-service')

type CreateChannelWrapper = VueWrapper<InstanceType<typeof CreateChannel>>

describe('CreateChannel component', () => {
	let wrapper: CreateChannelWrapper

	beforeEach(() => {
		wrapper = mount(CreateChannel, {
			attachTo: document.body
		})
	})

	afterEach(() => {
		wrapper.unmount()
	})

	function formElement(wrapper: CreateChannelWrapper) {
		return wrapper.find('[data-test="form"]')
	}

	function createChannelButton(wrapper: CreateChannelWrapper) {
		return wrapper.find('[data-test="create channel button"]')
	}

	function createChannelError(wrapper: CreateChannelWrapper) {
		return wrapper.find('[data-test="create channel error"]')
	}

	function channelNameInput(wrapper: CreateChannelWrapper) {
		return wrapper.get('[data-test="channel name input"]')
	}

	it('shows the create channel form when create channel button is clicked', async () => {
		expect(formElement(wrapper).exists()).toBe(false)
		expect(createChannelButton(wrapper).exists()).toBe(true)

		createChannelButton(wrapper).trigger('click')
		await nextTick()

		expect(formElement(wrapper).exists()).toBe(true)
		expect(createChannelButton(wrapper).exists()).toBe(false)
	})

	it('focuses the channel name input when the create channel form is shown', async () => {
		createChannelButton(wrapper).trigger('click')
		await nextTick()
		await nextTick() // watcher + nextTick inside

		expect(channelNameInput(wrapper).element).toBe(document.activeElement)
	})

	it('hides the create channel form when cancel button is clicked', async () => {
		createChannelButton(wrapper).trigger('click')
		await nextTick()
		expect(formElement(wrapper).exists()).toBe(true)

		wrapper.get('[data-test="cancel"]').trigger('click')
		await nextTick()

		expect(formElement(wrapper).exists()).toBe(false)
		expect(createChannelButton(wrapper).exists()).toBe(true)
	})

	it('hides the create channel form when the Escape key is pressed in the channel name input', async () => {
		createChannelButton(wrapper).trigger('click')
		await nextTick()
		expect(formElement(wrapper).exists()).toBe(true)

		channelNameInput(wrapper).trigger('keydown.esc')
		await nextTick()

		expect(formElement(wrapper).exists()).toBe(false)
		expect(createChannelButton(wrapper).exists()).toBe(true)
	})

	describe('submitting the form', () => {
		beforeEach(async () => {
			createChannelButton(wrapper).trigger('click')
			await nextTick()
		})

		async function submitNewChannel(channelName: string = "anything") {
			channelNameInput(wrapper).setValue(channelName)
			formElement(wrapper).trigger("submit")
			await nextTick()
		}

		it('makes a request to create a new channel when the channel form is submitted', async () => {
			vi.mocked(createChannel).mockResolvedValue({ id: "1", name: "Anything" })
			const channelName = "New channel"
			await submitNewChannel(channelName)
			expect(createChannel).toHaveBeenCalledWith(channelName)
		})

		describe('with an empty channel name', () => {
			it('shows an error message', async () => {
				expect(createChannelError(wrapper).exists()).toBe(false)
				await submitNewChannel("")
				expect(createChannelError(wrapper).exists()).toBe(true)
				expect(createChannel).not.toHaveBeenCalled()
				expect(createChannelError(wrapper).text()).toEqual('Please enter a channel name.')
			})
		})

		describe('with a channel name containing all whitespace', () => {
			it('shows an error message', async () => {
				expect(createChannelError(wrapper).exists()).toBe(false)
				await submitNewChannel(" \t\n\r")
				expect(createChannelError(wrapper).exists()).toBe(true)
				expect(createChannel).not.toHaveBeenCalled()
				expect(createChannelError(wrapper).text()).toEqual('Please enter a channel name.')
			})
		})

		describe('when the new channel failed to be created', () => {
			beforeEach(() => {
				vi.mocked(createChannel).mockRejectedValue(new Error("error"))
			})

			it('does not hide the create channel form', async () => {
				expect(formElement(wrapper).exists()).toBe(true)
				await submitNewChannel()
				expect(wrapper.find('[data-test="submit"]').exists()).toBe(true)
			})

			describe('when a generic error is returned', () => {
				it('shows a generic error message', async () => {
					expect(createChannelError(wrapper).exists()).toBe(false)
					await submitNewChannel()
					expect(createChannelError(wrapper).exists()).toBe(true)
					expect(createChannelError(wrapper).text()).toEqual("Failed to create channel. Please try again.")
				})
			})

			describe('when an AxiosError is returned', () => {
				var error: AxiosError
				var errorMessage: string
				beforeEach(() => {
					errorMessage = "Axios error messsage"
					error = new AxiosError(
						'error',
						'error',
						undefined,
						undefined,
						{
							data: errorMessage,
							status: 422,
							statusText: 'Unprocessible entity',
							headers: new AxiosHeaders(),
							config: { headers: new AxiosHeaders() }
						}
					)
					vi.mocked(createChannel).mockRejectedValue(error)
				})

				it('shows the error message in the AxiosError', async () => {
					await submitNewChannel()
					expect(createChannelError(wrapper).exists()).toBe(true)
					expect(createChannelError(wrapper).text()).toEqual(errorMessage)
				})
			})
		})

		describe('when the new channel was created successfully', () => {
			let newChannel: Channel
			beforeEach(() => {
				newChannel = { id: "1", name: "B Channel name" }
				vi.mocked(createChannel).mockResolvedValue(newChannel)
			})

			it('hides the new channel form', async () => {
				expect(formElement(wrapper).exists()).toBe(true)
				expect(createChannelButton(wrapper).exists()).toBe(false)

				await submitNewChannel(newChannel.name)

				expect(formElement(wrapper).exists()).toBe(false)
				expect(createChannelButton(wrapper).exists()).toBe(true)
			})

			it('emits the newly created channel', async () => {
				expect(wrapper.emitted('channelSelected')).toBeFalsy()

				await submitNewChannel(newChannel.name)

				expect(wrapper.emitted('channelCreated')).toBeTruthy()
				const emittedChannel = wrapper.emitted('channelCreated')?.[0]?.[0]
				expect(emittedChannel).toEqual(newChannel)
			})

			describe('when a error from a previous attempt to create a channel is present', () => {
				beforeEach(async () => {
					vi.mocked(createChannel).mockRejectedValue(new Error("error"))
					await submitNewChannel()
				})

				it('removes the old error message', async () => {
					expect(createChannelError(wrapper).exists()).toBe(true)

					vi.mocked(createChannel).mockResolvedValue({ id: '1', name: 'Test' })
					await submitNewChannel()
					await nextTick()

					expect(createChannelError(wrapper).exists()).toBe(false)
				})
			})
		})
	})
})

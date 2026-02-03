// @vitest-environment happy-dom

import { mount, type VueWrapper } from '@vue/test-utils'
import CreateMessage from '@/components/CreateMessage.vue'
import { type Channel, createMessage } from '@/utils/channel-service'
import { nextTick } from 'vue'

type CreateMessageWrapper = VueWrapper<InstanceType<typeof CreateMessage>>

vi.mock('@/utils/channel-service')

describe('CreateMessage component', () => {
	let wrapper: CreateMessageWrapper
	let channel: Channel

	beforeEach(() => {
		channel = { id: "1", name: "Channel name" }
		wrapper = mount(CreateMessage, {
			props: {
				channel
			}
		})
	})

	describe("submitting the form", () => {
		async function submit(message: string) {
			wrapper.find('textarea').setValue(message)
			wrapper.find('form').trigger('submit')
			await nextTick()
		}

		it('submits a new message to the provided channel', async () => {
			const message = "New message"
			await submit(message)
			expect(createMessage).toBeCalledWith(channel.id, message)
		})

		describe('when the message it successfully created', () => {
			let errorMessage: string

			beforeEach(() => {
				errorMessage = "Error message"
				vi.mocked(createMessage).mockRejectedValue(new Error("Error"))
			})

			it('shows the error message', async () => {
				expect(wrapper.find('[data-error]').exists()).toBe(false)
				await submit("message")
				const errorElement = wrapper.get('[data-error]')
				expect(errorElement.text()).toEqual("Unable to create message. Please try again.")
			})
		})
	})
})

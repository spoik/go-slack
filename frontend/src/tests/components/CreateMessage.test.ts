// @vitest-environment happy-dom

import { mount, type VueWrapper } from '@vue/test-utils'
import CreateMessage from '@/components/CreateMessage.vue'
import { type Channel, createMessage, Message } from '@/utils/channel-service'
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
			messageBodyInput().setValue(message)
			formElement().trigger('submit')
			await nextTick()
		}

		function formElement() {
			return wrapper.find('form')
		}

		function messageBodyInput() {
			return wrapper.find('textarea')
		}

		function errorMessageElement() {
			return wrapper.find('[data-error]')
		}

		it('submits a new message to the channel passed in props', async () => {
			const message = "New message"
			await submit(message)
			expect(createMessage).toBeCalledWith(channel.id, message)
		})

		describe('when the message is not successfully created', () => {
			let errorMessage: string

			beforeEach(() => {
				errorMessage = "Error message"
				vi.mocked(createMessage).mockRejectedValue(new Error("Error"))
			})

			it('shows the error message', async () => {
				expect(errorMessageElement().exists()).toBe(false)

				await submit("Message")

				const errorElement = errorMessageElement()
				expect(errorElement.text()).toEqual("Unable to create message. Please try again.")
			})

			it("does not emit the newly created message", async () => {
				expect(wrapper.emitted("messageCreated")).toBeUndefined()
				await submit("Message")
				expect(wrapper.emitted("messageCreated")).toBeUndefined()
			})

			it("does not clear out the message text from the from", async () => {
				const message = "New message text"
				messageBodyInput().setValue(message)
				expect(messageBodyInput().element.value).toEqual(message)

				formElement().trigger("submit")
				await nextTick()

				expect(messageBodyInput().element.value).toEqual(message)
			})
		})

		describe('when the message is successfully created', () => {
			let message: Message

			beforeEach(() => {
				message = new Message("1", "Message body", new Date())
				vi.mocked(createMessage).mockResolvedValueOnce(message)
			})

			it("emits the newly created message", async () => {
				expect(wrapper.emitted("messageCreated")?.length).toBeUndefined()

				await submit("Message")

				expect(wrapper.emitted("messageCreated")?.length).toEqual(1)
				const emittedMessage = wrapper.emitted("messageCreated")?.[0]?.[0]
				expect(emittedMessage).toEqual(message)
			})

			it("clears the new message textarea", async () => {
				const newMessage = "New message body"
				messageBodyInput().setValue(newMessage)
				expect(messageBodyInput().element.value).toEqual(newMessage)

				formElement().trigger("submit")
				await nextTick()
				await nextTick()

				expect(messageBodyInput().element.value).toEqual("")
			})

			it("doesn't show the error message element", async () => {
				expect(errorMessageElement().exists()).toBe(false)
				await submit("Message")
				expect(errorMessageElement().exists()).toBe(false)
			})
		})
	})
})

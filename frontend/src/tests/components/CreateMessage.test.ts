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

	function messageBodyInput() {
		return wrapper.find('textarea')
	}

	beforeEach(() => {
		channel = { id: "1", name: "Channel name" }
		wrapper = mount(CreateMessage, {
			props: {
				channel
			}
		})
	})

	it("pressing enter in text area submits the new message", async () => {
		expect(createMessage).not.toHaveBeenCalled()
		messageBodyInput().trigger("keydown.enter")
		expect(createMessage).toHaveBeenCalled()
	})

	it("pressing shift+enter in text area does not submit the new message", async () => {
		expect(createMessage).not.toHaveBeenCalled()
		messageBodyInput().trigger("keydown", {
			key: 'Enter',
			shiftKey: true
		})
		expect(createMessage).not.toHaveBeenCalled()
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

			describe("then another message is submitted and is successful", () => {
				let message: Message

				beforeEach(async () => {
					await submit("Message")

					message = new Message("1", "Message body", new Date())
					vi.mocked(createMessage).mockResolvedValueOnce(message)
				})

				it("hides the old error message", async () => {
					wrapper.get('[data-error]')
					expect(errorMessageElement().exists()).toBe(true)
					await submit("New message")
					expect(errorMessageElement().exists()).toBe(false)
				})
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

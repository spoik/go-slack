// @vitest-environment happy-dom

import { mount, type VueWrapper } from '@vue/test-utils'
import CreateMessage from '@/components/CreateMessage.vue'
import { type Channel, createMessage } from '@/utils/channel-service'

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

	it('submits a new message to the provided channel', () => {
		const message = "New message"
		wrapper.find('textarea').setValue(message)
		wrapper.find('form').trigger('submit')
		expect(createMessage).toBeCalledWith(channel.id, message)
	})
})

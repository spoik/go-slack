// @vitest-environment happy-dom

import { nextTick } from 'vue'
import { type Mock } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import CurrentChannel from '@/components/CurrentChannel.vue'
import * as channelService from "@/utils/channel-service"
import { Message, type Channel } from "@/utils/channel-service"

describe('CurrentChannel component', () => {
    let getMessagesMocked: Mock<(channelId: string) => Promise<Message[]>>

    beforeEach(() => {
        getMessagesMocked = vi.spyOn(channelService, 'getMessages')
    })

    function mockGetMessages(messages: Message[]) {
        getMessagesMocked.mockResolvedValue(messages)
    }

    describe('when a channel is not provided', () => {
        let wrapper: VueWrapper<InstanceType<typeof CurrentChannel>>

        beforeEach(() => {
            wrapper = mount(CurrentChannel, {
                props: {
                    channel: undefined
                }
            })
        })

        it('shows a message when no channel is selected', () => {
            expect(wrapper.get('[data-test="channel empty message"]').text()).toContain('Please select a channel.')

            expect(wrapper.find('[data-test="error"]').exists()).toBe(false)
            expect(wrapper.find('[data-test="messages empty message"]').exists()).toBe(false)
        })

        it('does not make a request to get the messages in a channel', () => {
            expect(getMessagesMocked).not.toHaveBeenCalled()
        })
    })

    describe('when a channel is provided', () => {
        let channel: Channel

        async function initWrapper(): Promise<VueWrapper<InstanceType<typeof CurrentChannel>>> {
            channel = { id: '1', name: 'general' }

            const wrapper = mount(CurrentChannel, {
                props: {
                    channel
                }
            })

            await nextTick()

            return wrapper
        }

        it('makes a call to get the messages for the channel', () => {
            mockGetMessages([])
            initWrapper()
            expect(getMessagesMocked).toHaveBeenCalledWith(channel.id)
        })

        describe('when the messages failed to load', () => {
            it('shows an error message', async () => {
                getMessagesMocked.mockRejectedValue(channel.id)
                const wrapper = await initWrapper()
                expect(wrapper.find('[data-test="error"]').exists()).toBe(true)
                expect(wrapper.find('[data-test="messages empty message"]').exists()).toBe(false)
            })
        })

        describe('when the channel has no messages', () => {
            it('shows an empty messages', async () => {
                mockGetMessages([])
                const wrapper = await initWrapper()
                expect(wrapper.find('[data-test="messages empty message"]').exists()).toBe(true)
                expect(wrapper.find('[data-test="error"]').exists()).toBe(false)
            })
        })

        describe('when the channel has messages', () => {
            let wrapper: VueWrapper<InstanceType<typeof CurrentChannel>>
            let messages: Message[]

            beforeEach(async () => {
                const date = new Date(2024, 0, 20, 10, 30, 0); // January 20, 2024, 10:30:00
                messages = [
                    new Message("1", "Test 1", date),
                    new Message("2", "Test 2", new Date())
                ]

                mockGetMessages(messages)
                wrapper = await initWrapper()
            })

            it("doesn't show an empty or error messages", () => {
                expect(wrapper.find('[data-test="error"]').exists()).toBe(false)
                expect(wrapper.find('[data-test="messages empty message"]').exists()).toBe(false)
            })

            it("shows the messages", async () => {
                expect(wrapper.find('[data-test="messages"]').exists()).toBe(true)
                expect(wrapper.find('[data-test-message]').exists()).toBe(true)
                expect(wrapper.findAll('[data-test-message]').length).toEqual(2)
                expect(wrapper.get('[data-test-message="1"] [data-test="message text"]').text()).toEqual(messages[0]?.message)
                expect(wrapper.get('[data-test-message="2"] [data-test="message text"]').text()).toEqual(messages[1]?.message)
            })

            it('shows the created at date for each message', () => {
                const timeElement = wrapper.find(`[data-test-message="${messages[0]?.id}"] time`);
                expect(timeElement.exists()).toBe(true);
                expect(timeElement.attributes('datetime')).toEqual(messages[0]?.created_at.toISOString());
                expect(timeElement.text()).toEqual(messages[0]?.created_at.toLocaleString());
            });
        })
    })

    describe('when the channel is changed', () => {
        let channel1: Channel
        let channel1Messages: Message[]

        let channel2: Channel
        let channel2Messages: Message[]

        beforeEach(() => {
            channel1 = { id: '1', name: 'general' }
            channel1Messages = [new Message('1', 'message 1', new Date())]

            channel2 = { id: '2', name: 'test' }
            channel2Messages = [
                new Message('2', 'message 2', new Date()),
                new Message('3', 'message 3', new Date())
            ]
        })

        async function initWrapper(): Promise<VueWrapper<InstanceType<typeof CurrentChannel>>> {
            mockGetMessages(channel1Messages)

            const wrapper = mount(CurrentChannel, {
                props: {
                    channel: channel1
                }
            })

            await nextTick()

            return wrapper
        }

        async function changeChannel(wrapper: VueWrapper<InstanceType<typeof CurrentChannel>>) {
            await wrapper.setProps({ channel: channel2 })
            await nextTick()
        }

        it('loads messages from the new channel', async () => {
            const wrapper = await initWrapper()
            expect(wrapper.findAll('[data-test-message]').length).toEqual(1)
            expect(wrapper.get('[data-test-message="1"] [data-test="message text"]').text()).toEqual(channel1Messages[0]?.message)

            mockGetMessages(channel2Messages)
            await changeChannel(wrapper)

            expect(wrapper.findAll('[data-test-message]').length).toEqual(2)
            expect(wrapper.get('[data-test-message="2"] [data-test="message text"]').text()).toEqual(channel2Messages[0]?.message)
            expect(wrapper.get('[data-test-message="3"] [data-test="message text"]').text()).toEqual(channel2Messages[1]?.message)
        })
    })
})

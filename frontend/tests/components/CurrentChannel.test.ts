// @vitest-environment happy-dom

import { nextTick } from 'vue'
import { describe, it, expect, beforeEach, Mock, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import CurrentChannel from '@/components/CurrentChannel.vue'
import { type Channel, type Message, getMessages } from "@/utils/channel-service"

vi.mock('@/utils/channel-service')

describe('CurrentChannel component', () => {
    let getMessagesMocked: Mock<typeof getMessages>

    beforeEach(() => {
        getMessagesMocked = vi.mocked(getMessages)
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
            expect(wrapper.find('h1').exists()).toBe(false)
        })

        it('does not make a request to get the messages in a channel', () => {
            expect(getMessagesMocked).not.toHaveBeenCalled()
        })
    })

    describe('when a channel is provided', () => {
        const channel: Channel = { id: '1', name: 'general' }

        async function initWrapper(): VueWrapper<InstanceType<typeof CurrentChannel>> {
            const wrapper = mount(CurrentChannel, {
                props: {
                    channel
                }
            })

            await nextTick()

            return wrapper
        }

        it('shows the channel name when a channel is selected', async () => {
            const wrapper = await initWrapper()
            expect(wrapper.get('h1').text()).toBe(channel.name)
            expect(wrapper.find('[data-test="channel empty message"]').exists()).toBe(false)
            expect(wrapper.find('[data-test="error"]').exists()).toBe(false)
        })

        it('makes a call to get the messages for the channe', () => {
            initWrapper()
            expect(getMessagesMocked).toHaveBeenCalledWith(channel.id)
        })

        describe('when the messages failed to load', () => {
            it('shows an error message', async () => {
                getMessagesMocked.mockRejectedValue()
                const wrapper = await initWrapper()
                expect(wrapper.get('[data-test="error"]').exists()).toBe(true)
                expect(wrapper.find('[data-test="messages empty message"]').exists()).toBe(false)
            })
        })

        describe('when the channel has no messages', () => {
            it('shows an empty messages', async () => {
                mockGetMessages([])
                const wrapper = await initWrapper()
                expect(wrapper.get('[data-test="messages empty message"]').exists()).toBe(true)
                expect(wrapper.find('[data-test="error"]').exists()).toBe(false)
            })
        })

        describe('when the channel has messages', () => {
            let wrapper: VueWrapper<InstanceType<typeof CurrentChannel>>
            const messages: Message[] = [
                { id: "1", message: "Test 1" },
                { id: "2", message: "Test 2" },
            ]

            beforeEach(async () => {
                mockGetMessages(messages)
                wrapper = await initWrapper()
            })

            it("doesn't show an empty or error messages", () => {
                expect(wrapper.find('[data-test="error"]').exists()).toBe(false)
                expect(wrapper.find('[data-test="messages empty message"]').exists()).toBe(false)
            })

            it("shows the messages", () => {
                expect(wrapper.get('[data-test="messages"]').exists()).toBe(true)
                expect(wrapper.findAll('[data-test-message]').length).toEqual(2)
                expect(wrapper.get('[data-test-message="1"]').text()).toEqual(messages[0].message)
                expect(wrapper.get('[data-test-message="2"]').text()).toEqual(messages[1].message)
            })
        })
    })

    describe('when the channel is changed', () => {
        const channel1: Channel = { id: '1', name: 'general' }
        const channel1Messages: Message[] = [{ id: '1', message: 'message 1'}]

        const channel2: Channel = { id: '2', name: 'test' }
        const channel2Messages: Message[] = [
            { id: '2', message: 'message 2'},
            { id: '3', message: 'message 3'},
        ]

        async function initWrapper(): VueWrapper<InstanceType<typeof CurrentChannel>> {
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

        it('updates the name of the channel', async () => {
            const wrapper = await initWrapper()
            expect(wrapper.get('h1').text()).toEqual(channel1.name)

            await changeChannel(wrapper)
            expect(wrapper.get('h1').text()).toEqual(channel2.name)
        })

        it('loads messages from the new channel', async () => {
            const wrapper = await initWrapper()
            expect(wrapper.findAll('[data-test-message]').length).toEqual(1)
            expect(wrapper.get('[data-test-message="1"]').text()).toEqual(channel1Messages[0].message)

            mockGetMessages(channel2Messages)
            await changeChannel(wrapper)
            expect(wrapper.findAll('[data-test-message]').length).toEqual(2)
            expect(wrapper.get('[data-test-message="2"]').text()).toEqual(channel2Messages[0].message)
            expect(wrapper.get('[data-test-message="3"]').text()).toEqual(channel2Messages[1].message)
        })
    })
})

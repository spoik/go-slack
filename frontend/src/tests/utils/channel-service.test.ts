import { vi, describe, it, expect } from 'vitest'
import axios from 'axios'
import { getChannels, getMessages, createChannel, createMessage, type Channel, Message } from '@/utils/channel-service'

// Mock axios
vi.mock('axios')

describe('channel-service', () => {
  describe('getChannels', () => {
    it('should call axios.get with the correct URL', async () => {
      const mockChannels: Channel[] = [
        { id: '1', name: 'general' },
        { id: '2', name: 'random' },
      ]
      vi.mocked(axios.get).mockResolvedValue({ data: mockChannels })

      const channels = await getChannels()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:8000/channels')
      expect(channels).toEqual(mockChannels)
    })
  })

  describe('createChannel', () => {
    it('makes a network call to create the channel', async () => {
      const channelName = "Channel name"
      const mockChannel: Channel = { id: '1', name: channelName }
      vi.mocked(axios.post).mockResolvedValue({ data: mockChannel })

      const returnedChannel = await createChannel(channelName)

      const expectedData = { name: channelName }
      expect(axios.post).toHaveBeenCalledWith('http://localhost:8000/channels', expectedData)

      expect(returnedChannel.id).toEqual(mockChannel.id)
      expect(returnedChannel.name).toEqual(mockChannel.name)
    })
  })

  describe('getChannelMessages', () => {
    it('should call axios.get with the correct URL', async () => {
      const channelId = "4"
      const mockMessages: Message[] = [
        new Message('1', 'message 1', new Date()),
        new Message('2', 'message 2', new Date())
      ]

      vi.mocked(axios.get).mockResolvedValue({ data: mockMessages })

      const messages = await getMessages(channelId)

      expect(axios.get).toHaveBeenCalledWith(`http://localhost:8000/channels/${channelId}/messages`)
      expect(messages).toEqual(mockMessages)
    })
  })

  describe('createMessage', () => {
    it('makes a network call to create the channel', async () => {
      const message = "This is the message"
      const channelId = "2"
      const mockMessage = new Message("1", message, new Date())
      vi.mocked(axios.post).mockResolvedValue({ data: mockMessage })

      const returnedMessage = await createMessage(channelId, message)

      expect(axios.post).toHaveBeenCalledWith(
        'http://localhost:8000/channels/2/messages',
        { message }
      )

      expect(returnedMessage.id).toEqual(mockMessage.id)
      expect(returnedMessage.message).toEqual(mockMessage.message)
      expect(returnedMessage.created_at).toEqual(mockMessage.created_at)
    })
  })
})

import { describe, it, expect, vi } from 'vitest'
import axios from 'axios'
import { getChannels, getMessages, type Channel, type Message } from '../../src/utils/channel-service'

// Mock axios
vi.mock('axios')

describe('channel-service', () => {
  describe('getChannels', () => {
    it('should call axios.get with the correct URL', async () => {
      const mockChannels: Channel[] = [
        { id: '1', name: 'general' },
        { id: '2', name: 'random' },
      ]
      axios.get.mockResolvedValue({ data: mockChannels })

      const channels = await getChannels()

      expect(axios.get).toHaveBeenCalledWith('http://localhost:8000/channels')
      expect(channels).toEqual(mockChannels)
    })
  })

  describe('getChannelMessages', () => {
    it('should call axios.get with the correct URL', async () => {
      const channelId = "4"
      const mockMessages: Message[] = [
        { id: '1', message: 'message 1' },
        { id: '2', message: 'message 2' },
      ]
      
      axios.get.mockResolvedValue({ data: mockMessages})

      const messages = await getMessages(channelId)

      expect(axios.get).toHaveBeenCalledWith(`http://localhost:8000/channels/${channelId}/messages`)
      expect(messages).toEqual(mockMessages)
    })
  })
})

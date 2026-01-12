import { describe, it, expect, vi } from 'vitest'
import axios from 'axios'
import { getChannels } from '../../src/utils/channel-service'

// Mock axios
vi.mock('axios')

describe('channel-service', () => {
  it('should call axios.get with the correct URL when fetching channels', async () => {
    const mockChannels = [
      { id: '1', name: 'general' },
      { id: '2', name: 'random' },
    ];
    axios.get.mockResolvedValue({ data: mockChannels })

    const channels = await getChannels()

    expect(axios.get).toHaveBeenCalledWith('http://localhost:8000/channels')
    expect(channels).toEqual(mockChannels)
  });
});

import axios from 'axios'

interface Channel {
	readonly id: string;
	name: string;
}

export async function getChannels(): Promise<Channel[]> {
	const { data } = await axios.get<Channel>('http://localhost:8000/channels')
	return data
}

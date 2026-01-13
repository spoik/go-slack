import axios from 'axios'

export interface Channel {
	readonly id: string;
	name: string;
}

export async function getChannels(): Promise<Channel[]> {
	const { data } = await axios.get<Channel[]>('http://localhost:8000/channels')
	return data
}

export interface Message {
	readonly id: string;
	message: string;
}

export async function getMessages(channelId: string): Promise<Message[]> {
	const { data } = await axios.get<Message[]>(`http://localhost:8000/channels/${channelId}/messages`)
	return data
}

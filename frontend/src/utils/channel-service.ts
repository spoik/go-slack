import 'reflect-metadata'
import axios from 'axios'
import { plainToInstance, Type } from 'class-transformer'

export interface Channel {
	readonly id: string;
	name: string;
}

export async function getChannels(): Promise<Channel[]> {
	const { data } = await axios.get<Channel[]>('http://localhost:8000/channels')
	return data
}

interface MessageData {
	readonly id: string;
	message: string;
	created_at: string;
}

export class Message {
	readonly id: string;
	message: string;

	@Type(() => Date)
	created_at: Date;

	constructor(id: string, message: string, created_at: Date | string) {
		this.id = id
		this.message = message
		this.created_at = typeof created_at === "string" ? new Date(created_at) : created_at
	}

	formattedDate(): string {
		return this.created_at.toLocaleString()
	}
}

export async function getMessages(channelId: string): Promise<Message[]> {
	const { data } = await axios.get<MessageData[]>(`http://localhost:8000/channels/${channelId}/messages`)
	return plainToInstance(Message, data)
}

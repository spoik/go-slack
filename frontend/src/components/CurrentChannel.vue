<script setup lang="ts">
import { toRefs, ref, onMounted, computed, watch } from "vue"
import { type Channel, Message, getMessages } from "@/utils/channel-service"
import CreateMessage from "@/components/CreateMessage.vue"

const props = defineProps<{ channel?: Channel }>()
const { channel } = toRefs(props)
const messages = ref<Message[] | undefined>(undefined)
const error = ref<string | undefined>(undefined)

async function loadMessages() {
    if (channel.value == null) {
        return
    }

    try {
        messages.value = await getMessages(channel.value.id)
    } catch (e) {
        error.value = "An error occurred when loading the channel. Please try again."
    }
}

function newMessageCreated(newMessage: Message) {
    messages.value?.push(newMessage)
}

const hasMessages = computed(() => {
    return messages.value != undefined && messages.value.length > 0
})

watch(channel, loadMessages)

onMounted(loadMessages)
</script>

<template>
    <div class="h-screen flex flex-col">
        <p v-if="channel == null" data-test="channel empty message">Please select a channel.</p>

        <template v-else>
            <div class="min-h-0 overflow-y-auto flex-1 p-5">
                <p v-if="error != undefined" data-test="error">{{ error }}</p>

                <p v-else-if="!hasMessages" data-test="messages empty message">No messages in this channel yet.</p>

                <ul dv-else class="flex flex-col gap-5" data-test="messages">
                    <li v-for="message in messages" :key="message.id" :data-test-message="message.id" class="flex">
                        <p class="grow" data-test="message text">{{ message.message }}</p>

                        <time data-test="datetime" class="text-gray-400" :datetime="message.created_at.toISOString()">
                            {{ message.formattedDate() }}
                        </time>
                    </li>
                </ul>
            </div>

            <div v-if="error == undefined" class="p-5 border-t border-zinc-700">
                <CreateMessage :channel="channel" @message-created="newMessageCreated" />
            </div>
        </template>
    </div>
</template>

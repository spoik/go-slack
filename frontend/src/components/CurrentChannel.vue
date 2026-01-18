<script setup lang="ts">
import { toRefs, ref, onMounted, computed, watch } from "vue"
import { type Channel, type Message, getMessages } from "@/utils/channel-service"

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

const hasMessages = computed((): boolean => {
    return messages.value != undefined && messages.value.length > 0
})

watch(channel, loadMessages)

onMounted(loadMessages)
</script>

<template>
    <div class="p-5 border">
        <p v-if="channel == null" data-test="channel empty message">Please select a channel.</p>

        <template v-else>
            <p v-if="error != undefined" data-test="error">{{ error }}</p>
            <p v-else-if="!hasMessages" data-test="messages empty message">No messages in this channel yet.</p>
            <ul v-else data-test="messages">
                <li v-for="message in messages" :key="message.id" :data-test-message="message.id">
                    {{ message.message }}
                </li>
            </ul>
        </template>
    </div>
</template>

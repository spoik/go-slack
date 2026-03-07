<script setup lang="ts">
import { ref } from "vue"
import { createMessage, Message, type Channel } from '@/utils/channel-service'

const props = defineProps<{ channel: Channel }>()
const emit = defineEmits<{
    messageCreated: [message: Message]
}>()


const message = ref<string>("")
const errorMessage = ref<string | null>(null)

async function createNewMessage() {
    let newMessage: Message

    try {
        newMessage = await createMessage(props.channel.id, message.value)
    } catch (error: any) {
        errorMessage.value = "Unable to create message. Please try again."
        return
    }

    errorMessage.value = null
    message.value = ""
    emit("messageCreated", newMessage)
}
</script>

<template>
    <div>
        <form @submit.prevent="createNewMessage" class="flex items-center">
            <textarea placeholder="New message"
                class="w-full px-3 py-2 bg-zinc-700 rounded-sm min-h-10.5 focus:border-slate-500 focus:ring-slate-500 focus:ring-1 outline-none"
                name="message" rows=1 v-model="message" @keydown.enter.exact="createNewMessage"></textarea>

            <input type="submit" class="bg-indigo-400 text-white px-4 py-2 ml-4 h-auto rounded-xs" value="Submit">
        </form>

        <p v-if="errorMessage != null" class="mt-2 text-red-400" data-error>{{ errorMessage }}</p>
    </div>
</template>

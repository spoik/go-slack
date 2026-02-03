<script setup lang="ts">
import { ref } from "vue"
import { createMessage, type Channel } from '@/utils/channel-service'

const props = defineProps<{ channel: Channel }>()
const message = ref<string>("")

function createNewMessage() {
    createMessage(props.channel.id, message.value)
}
</script>

<template>
    <div>
        <form @submit.prevent="createNewMessage" class="flex items-center">
            <textarea placeholder="New message" class="w-full px-3 py-2 border" name="message" rows=1
                v-model="message"></textarea>

            <input type="submit" class="bg-indigo-400 text-white px-4 py-2 ml-4 h-auto">
        </form>
    </div>
</template>

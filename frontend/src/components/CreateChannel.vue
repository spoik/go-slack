<script setup lang="ts">
import { ref } from 'vue'
import { createChannel, type Channel } from "@/utils/channel-service"
import { AxiosError } from "axios";

const emit = defineEmits<{
    channelCreated: [channel: Channel]
}>()

const showForm = ref<boolean>(false)
const errorMessage = ref<string | null>(null)
const newChannelName = ref<string>('')

async function createNewChannel() {
    errorMessage.value = null

    let newChannel: Channel

    try {
        newChannel = await createChannel(newChannelName.value)
    } catch (error: any) {
        if (error instanceof AxiosError && error.response != undefined) {
            errorMessage.value = error.response.data
        } else {
            errorMessage.value = "Failed to create channel. Please try again."
        }

        return
    }

    showForm.value = false

    emit('channelCreated', newChannel)
}
</script>

<template>
    <button v-if="!showForm" @click="showForm = true" data-test="create channel button">+
        create channel</button>

    <form v-else @submit.prevent="createNewChannel" class="flex flex-row gap-4" data-test="form">
        <label for="name" class="sr-only">New Channel Name</label>
        <input type="text" name="name" placeholder="New channel name" v-model="newChannelName"
            class="grow min-w-0 border p-1" data-test="channel name input" />
        <input type="submit" class="bg-indigo-400 px-4" data-test="submit" />
    </form>

    <p v-if="errorMessage != null" class="text-red-400" data-test="create channel error">{{ errorMessage
        }}</p>
</template>

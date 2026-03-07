<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { createChannel, type Channel } from "@/utils/channel-service"
import { AxiosError } from "axios";

const emit = defineEmits<{
    channelCreated: [channel: Channel]
}>()

const showForm = ref<boolean>(false)
const errorMessage = ref<string | null>(null)
const newChannelName = ref<string>('')
const nameInput = ref<HTMLInputElement | null>(null)

watch(showForm, async (newValue) => {
    if (newValue) {
        await nextTick()
        nameInput.value?.focus()
    } else {
        errorMessage.value = null
    }
})

async function createNewChannel() {
    errorMessage.value = null

    if (newChannelName.value.trim() == "") {
        errorMessage.value = "Please enter a channel name."
        return
    }

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
    <div class="mt-2.5">
        <button v-if="!showForm" @click="showForm = true" class="hover:cursor-pointer"
            data-test="create channel button">
            + create channel
        </button>

        <form v-else @submit.prevent="createNewChannel" class="mt-4" data-test="form">
            <p v-if="errorMessage != null" class="text-red-400 mb-1.5" data-test="create channel error">{{ errorMessage
            }}</p>

            <label for="name" class="sr-only">New Channel Name</label>
            <input type="text" name="name" placeholder="New channel name" v-model="newChannelName"
                @keydown.esc="showForm = false"
                ref="nameInput"
                class="w-full border py-1.5 px-3 rounded-sm" data-test="channel name input" />

            <div class="flex gap-4 mt-3">
                <button type="button" @click="showForm = false"
                    class="px-4 py-1 bg-red-400 grow text-gray-100 rounded-xs" data-test="cancel">Cancel</button>
                <input type="submit" class="bg-indigo-400 px-4 py-1 grow rounded-xs" data-test="submit"
                    value="Create" />
            </div>
        </form>

    </div>
</template>

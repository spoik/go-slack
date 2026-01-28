<script setup lang="ts">
import { onMounted, ref, computed, toRefs, watch } from "vue"
import { getChannels, createChannel, type Channel } from "@/utils/channel-service"

const emit = defineEmits<{
    channelSelected: [channel: Channel]
}>()

const props = defineProps<{ selectedChannel?: Channel }>()
const { selectedChannel } = toRefs(props)

const channels = ref<Channel[] | null>(null)
watch(channels, (newValue: Channel[] | null) => {
    if (channels.value == null || newValue == null) {
        return
    }
    channels.value = newValue.sort((a, b) => a.name.localeCompare(b.name))
}, { deep: true })

const error = ref<string | null>(null)
const showCreateChannel = ref<boolean>(false)
const loading = computed((): boolean => {
    return error.value == null && channels.value == null
})

const newChannelName = ref<string>('')

async function loadChannels() {
    try {
        channels.value = await getChannels()
    } catch (err) {
        error.value = "Failed to load channels. Please reload to try again."
    }
}

async function createNewChannel() {
    const newChannel = await createChannel(newChannelName.value)
    if (channels.value == null) {
        channels.value = []
    }
    channels.value.push(newChannel)
    showCreateChannel.value = false
}

onMounted(loadChannels)
</script>

<template>
    <div class="p-5 border w-1/6">
        <h2 class="mb-1">Channels</h2>
        <hr class="mb-3">

        <p v-if="loading" data-test="loading">Loading...</p>
        <p v-if="error != null" data-test="error">{{ error }}</p>

        <div v-if="!loading">
            <ul v-if="channels?.length">
                <li v-for="channel in channels" :key="channel.id">
                    <a href="#" @click="$emit('channelSelected', channel)" data-test="channel" class="channel-title"
                        :class="{ 'active': channel.id == selectedChannel?.id }">
                        {{ channel.name }}
                    </a>
                </li>
            </ul>

            <button v-if="!showCreateChannel" @click="showCreateChannel = true" data-test="create channel button">+ create channel</button>

            <form v-if="showCreateChannel" @submit.prevent="createNewChannel" class="flex flex-row gap-4" data-test="create channel form">
                <label for="name" class="sr-only">New Channel Name</label>
                <input type="text" name="name" placeholder="New channel name" v-model="newChannelName" class="grow min-w-0 border p-1" data-test="channel name input"/>
                <input type="submit" class="bg-indigo-400 px-4" data-test="new channel submit"/>
            </form>
        </div>
    </div>
</template>

<style scoped>
@reference "tailwindcss";

.channel-title {
    @apply p-2 block;
}

.channel-title:hover {
    @apply bg-gray-100;
}

.channel-title.active {
    @apply bg-gray-300;
}
</style>

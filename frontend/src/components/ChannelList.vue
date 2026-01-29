<script setup lang="ts">
import { onMounted, ref, computed, toRefs, watch } from "vue"
import { getChannels, type Channel } from "@/utils/channel-service"
import CreateChannel from "./CreateChannel.vue";

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

const loadChannelsError = ref<string | null>(null)
const loading = computed((): boolean => {
    return loadChannelsError.value == null && channels.value == null
})

async function loadChannels() {
    try {
        channels.value = await getChannels()
    } catch (error: any) {
        loadChannelsError.value = "Failed to load channels. Please reload to try again."
    }
}

function channelCreated(channel: Channel) {
    if (channels.value == null) {
        channels.value = []
    }

    channels.value.push(channel)
}

onMounted(loadChannels)
</script>

<template>
    <div class="p-5 border w-1/6">
        <h2 class="mb-1">Channels</h2>
        <hr class="mb-3">

        <p v-if="loading" data-test="loading">Loading...</p>
        <p v-if="loadChannelsError != null" data-test="error">{{ loadChannelsError }}</p>

        <div v-if="!loading">
            <ul v-if="channels?.length">
                <li v-for="channel in channels" :key="channel.id">
                    <a href="#" @click="$emit('channelSelected', channel)" data-test="channel" class="channel-title"
                        :class="{ 'active': channel.id == selectedChannel?.id }">
                        {{ channel.name }}
                    </a>
                </li>
            </ul>

            <CreateChannel @channel-created="channelCreated"/>
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

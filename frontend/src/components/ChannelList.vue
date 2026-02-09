<script setup lang="ts">
import { onMounted, ref, computed, toRefs, watch } from "vue"
import { RouterLink, useRoute } from "vue-router"
import { getChannels, type Channel } from "@/utils/channel-service"
import CreateChannel from "./CreateChannel.vue";

const emit = defineEmits<{
    channelSelected: [channel: Channel]
}>()

const route = useRoute()

const props = defineProps<{ selectedChannel?: Channel }>()
const { selectedChannel } = toRefs(props)

const channels = ref<Channel[] | null>(null)
watch(
    channels,
    (newValue: Channel[] | null) => {
        if (channels.value == null || newValue == null) {
            return
        }

        channels.value = newValue.sort((a, b) => a.name.localeCompare(b.name))
    },
    { deep: true }
)

watch(() => route.params.id, selectChannelFromRoute)

const loadChannelsError = ref<string | null>(null)
const loading = computed((): boolean => {
    return loadChannelsError.value == null && channels.value == null
})

async function loadChannels() {
    try {
        channels.value = await getChannels()
        selectChannelFromRoute()
    } catch (error: any) {
        loadChannelsError.value = "Failed to load channels. Please reload to try again."
    }
}

function selectChannelFromRoute() {
    if (route.params.id == null || channels.value == null) {
        return
    }

    const selectedChannel = channels.value.find((channel) => {
        return channel.id == route.params.id
    })

    if (selectedChannel != null) {
        emit("channelSelected", selectedChannel)
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
                    <RouterLink :to="{ name: 'channel', params: { id: channel.id } }" data-test="channel"
                        class="channel-title" :class="{ 'active': channel.id == selectedChannel?.id }">
                        {{ channel.name }}
                    </RouterLink>
                </li>
            </ul>

            <CreateChannel @channel-created="channelCreated" />
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

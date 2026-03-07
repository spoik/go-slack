<script setup lang="ts">
import { onMounted, ref, computed, toRefs, watch } from "vue"
import { RouterLink, useRoute, useRouter } from "vue-router"
import { getChannels, type Channel } from "@/utils/channel-service"
import CreateChannel from "./CreateChannel.vue";

const emit = defineEmits<{
    channelSelected: [channel: Channel]
}>()

const route = useRoute()
const router = useRouter()

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

    router.push({ name: 'channel', params: { id: channel.id } })
    emit('channelSelected', channel)
}

function linkAttrs(channel: Channel): object {
    if (channel.id == selectedChannel.value?.id) {
        return {
            class: "bg-zinc-700!",
            "data-test-active": ""
        }
    }

    return {}
}

onMounted(loadChannels)
</script>

<template>
    <div class="p-5 bg-zinc-900">
        <h2 class="mb-3.5 text-xl">Channels</h2>

        <p v-if="loading" data-test="loading">Loading...</p>
        <p v-if="loadChannelsError != null" data-test="error">{{ loadChannelsError }}</p>

        <div v-if="!loading">
            <ul v-if="channels?.length">
                <li v-for="channel in channels" :key="channel.id" class="mt-1.5">
                    <RouterLink :to="{ name: 'channel', params: { id: channel.id } }" data-test="channel"
                        class="py-2 px-3 block rounded-sm channel-title hover:bg-zinc-800"
                        v-bind="linkAttrs(channel)">
                        {{ channel.name }}
                    </RouterLink>
                </li>
            </ul>

            <CreateChannel @channel-created="channelCreated" />
        </div>
    </div>
</template>

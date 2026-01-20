<script setup lang="ts">
import { onMounted, ref, computed, toRefs } from "vue"
import { getChannels, type Channel } from "@/utils/channel-service"

const emit = defineEmits<{
    channelSelected: [channel: Channel]
}>()

const props = defineProps<{ selectedChannel?: Channel }>()
const { selectedChannel } = toRefs(props)

const channels = ref<Channel[] | null>(null)
const error = ref<string | null>(null)
const loading = computed((): boolean => {
    return error.value == null && channels.value == null
})

async function loadChannels() {
    try {
        channels.value = await getChannels()
    } catch (err) {
        error.value = "Failed to load channels. Please reload to try again."
    }
}

onMounted(loadChannels)
</script>

<template>
    <div class="p-5 border w-1/6">
        <h2 class="mb-1">Channels</h2>
        <hr class="mb-3">

        <p v-if="loading" data-test="loading">Loading...</p>
        <p v-if="error != null" data-test="error">{{ error }}</p>

        <ul v-if="channels?.length">
            <li v-for="channel in channels" :key="channel.id">
                <a href="#" @click="$emit('channelSelected', channel)" data-test="channel" class="channel-title"
                    :class="{ 'active': channel.id == selectedChannel?.id }">
                    {{ channel.name }}
                </a>
            </li>
        </ul>
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

<script setup lang="ts">
    import { onMounted, ref, computed } from "vue"
    import { getChannels, type Channel } from "@/utils/channel-service"

    const emit = defineEmits<{
        channelSelected: [channel: Channel]
    }>()

    const channels = ref<Channel[] | null>(null)
    const error = ref<string | null>(null)
    const loading = computed((): boolean => {
        return error.value == null && channels.value == null
    })

    async function loadChannels() {
        try {
            channels.value = await getChannels()
        } catch(err) {
            error.value = "Failed to load channels. Please reload to try again."
        }
    }

    onMounted(loadChannels)
</script>

<template>
    <div>
        <h1>Channels</h1>

        <p v-if="loading" data-test="loading">Loading...</p>
        <p v-if="error != null" data-test="error">{{ error }}</p>

        <ul v-if="channels?.length">
            <li
                v-for="channel in channels"
                :key="channel.id"
                @click="$emit('channelSelected', channel)"
                data-test="channel"
            >
                {{channel.name}}
            </li>
        </ul>
    </div>
</template>

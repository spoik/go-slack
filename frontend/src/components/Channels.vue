<script setup lang="ts">
    import { onMounted, ref, computed } from "vue"
    import { getChannels, type Channel } from "@/utils/channel-service"

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
            console.log("Failed to fetch channels:", err)
        }
    }

    onMounted(loadChannels)
</script>

<template>
    <h1>Channels</h1>

    <p v-if="loading" data-test="loading">Loading...</p>
    <p v-if="error != null" data-test="error">{{ error }}</p>

    <ul v-if="channels?.length">
        <li v-for="channel in channels" :key="channel.id">{{channel.name}}</li>
    </ul>
</template>

<script setup lang="ts">
    import { toRefs } from "vue"
    import { type Channel } from "@/utils/channel-service"

    const props = defineProps<{ channel?: Channel }>()
    const { channel } = toRefs(props)
</script>

<template>
    <div>
        <p v-if="channel == null" data-test="empty message">Please select a channel.</p>
        <div v-if="channel != null">
            <h1>{{ channel.name }}</h1>
        </div>
    </div>
</template>

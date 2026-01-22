<script setup lang="ts">
    import { ref } from "vue"
    import ChannelList from "@/components/ChannelList.vue"
    import CurrentChannel from "@/components/CurrentChannel.vue"
    import { type Channel } from "@/utils/channel-service.ts"

    const currentChannel = ref<Channel | undefined>()

    function channelSelected(channel: Channel) {
        currentChannel.value = channel
    }
</script>

<template>
    <div id="container" class="min-h-screen flex gap-10 p-5 lg:p-8">
        <ChannelList id="channel-list" @channel-selected="channelSelected" :selected-channel="currentChannel"/>
        <CurrentChannel id="current-channel" class="grow" :channel="currentChannel"/>
    </div>
</template>

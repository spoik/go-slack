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
    <div class="min-h-screen flex bg-zinc-800 text-zinc-200">
        <ChannelList id="channel-list" class="w-1/6" @channel-selected="channelSelected" :selected-channel="currentChannel" />
        <CurrentChannel id="current-channel" class="w-5/6" :channel="currentChannel" />
    </div>
</template>

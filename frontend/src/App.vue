<script setup lang="ts">
    import ChannelList from "./components/ChannelList.vue"
    import CurrentChannel from "./components/CurrentChannel.vue"
</script>

<template>
    <div id="container" class="min-h-screen flex gap-10 p-8">
        <ChannelList id="channel-list" class="p-8"/>
        <CurrentChannel id="current-channel" class="p-8 grow"/>
    </div>
</template>

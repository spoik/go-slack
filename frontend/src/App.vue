<script setup lang="ts">
    import Channels from "./components/Channels.vue"
    import Channel from "./components/Channel.vue"
</script>

<template>
    <div id="container" class="min-h-screen flex gap-10 p-8">
        <Channels id="channel-list" class="p-8"/>
        <Channel id="current-channel" class="p-8 grow"/>
    </div>
</template>

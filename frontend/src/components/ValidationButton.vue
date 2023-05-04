<template>
<v-btn size="large" prepend-icon="mdi-play" color="primary" :loading="loading" class="play" rounded>Start Validation</v-btn>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const loading = ref<boolean>(false)
</script>

<style scoped>
.play {
    position: fixed;
    bottom: 45px;
    right: 40px;
}
</style>
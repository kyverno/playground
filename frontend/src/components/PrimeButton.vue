<template>
    <v-btn :color="color"><slot /></v-btn>
</template>

<script setup lang="ts">
import { layoutTheme } from '@/config';
import { computed } from 'vue';

const color = computed(() => {
    if (layoutTheme.value === 'dark') {
        return 'secondary'
    }

    return 'primary'
})
</script>
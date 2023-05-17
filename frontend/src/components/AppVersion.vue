<template>
<v-card class="position background" flat :theme="layoutTheme" v-if="version">
    <v-card-text class="text-body-2 font-weight-bold pa-2 text-center px-3 py-1" style="font-size: 0.65rem!important">
        AppVersion: {{ version }}
    </v-card-text>
</v-card>
</template>

<script setup lang="ts">
import { layoutTheme } from '@/config';

const version = APP_VERSION
</script>

<style scoped>
.background {
    background-color: rgb(var(--v-theme-background))
}

.position {
  position: fixed;
  right: 0;
  bottom: 0px;
}
</style>
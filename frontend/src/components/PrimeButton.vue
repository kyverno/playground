<template>
  <v-btn :color="btnColor"><slot /></v-btn>
</template>

<script setup lang="ts">
import { btnColor } from '@/config'
</script>

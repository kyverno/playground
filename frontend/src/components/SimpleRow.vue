<template>
  <v-row :class="rowClass">
    <v-col :class="colClass">
      <slot />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
defineProps({
  rowClass: { type: String },
  colClass: { type: String }
})
</script>

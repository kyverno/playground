<template>
  <v-menu location="bottom" :close-on-content-click="false">
    <template v-slot:activator="{ props }">
      <v-btn class="ml-2" icon="mdi-menu" v-bind="props" />
    </template>
    <v-card min-width="250">
      <v-card-text>
        <ShareButton variant="outlined" block />
        <v-divider class="my-2" />
        <SaveButton variant="outlined" block />
        <v-divider class="my-2" />
        <LoadButton variant="outlined" block />
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import ShareButton from './ShareButton.vue'
import SaveButton from './SaveButton.vue'
import LoadButton from './LoadButton.vue'
</script>

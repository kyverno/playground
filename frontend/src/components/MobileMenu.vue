<template>
  <v-menu location="bottom" :close-on-content-click="false">
    <template v-slot:activator="{ props }">
      <v-btn  class="ml-2" icon="mdi-menu" v-bind="props" />
    </template>
    <v-card min-width="250">
      <v-card-text>
        <ShareButton variant="outlined" block :policy="policy" :resource="resource" :context="context" />
        <v-divider class="my-2" />
        <SaveButton variant="outlined" block :policy="policy" :resource="resource" :context="context" />
        <v-divider class="my-2" />
        <ResetButton variant="outlined" block @on-reset="emit('on-reset')" />
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import ShareButton from "@/components/ShareButton.vue";
import SaveButton from "@/components/SaveButton.vue";
import ResetButton from "@/components/ResetButton.vue";

defineProps({
  policy: { type: String, default: "" },
  resource: { type: String, default: "" },
  context: { type: String, default: "" },
});

const emit = defineEmits(['on-reset'])
</script>

<template>
  <v-toolbar color="#3783c4" theme="dark" density="compact" flat>
    <template #prepend>
      <v-btn
        small
        icon="mdi-backup-restore"
        v-if="restoreValue"
        :disabled="restoreValue === modelValue"
        @click="() => emit('update:modelValue', restoreValue)"
        :title="`Restore ${title}`"
      />
    </template>
    <v-toolbar-title :class="restoreValue ? 'ml-0' : ''">
      {{ title }}
      <v-tooltip :text="info" content-class="no-opacity-tooltip" v-if="info">
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon="mdi-information-outline" variant="text" size="small" />
        </template>
      </v-tooltip>
    </v-toolbar-title>
    <v-toolbar-items>
      <slot name="prepend-actions" />
      <UploadButton @click="(content: string) => emit('update:modelValue', content)" />
      <URLButton @click="(content: string) => emit('update:modelValue', content)" />
      <slot name="append-actions" />
    </v-toolbar-items>
  </v-toolbar>
</template>

<script setup lang="ts">
import UploadButton from "./UploadButton.vue";
import URLButton from "./URLButton.vue";

defineProps({
  restoreValue: { type: String, default: "" },
  modelValue: { type: String, default: "" },
  title: { type: String, required: true },
  info: { type: String },
});

const emit = defineEmits(["update:modelValue"]);
</script>

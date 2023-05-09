<template>
    <v-toolbar
      :title="props.title"
      color="#3783c4"
      theme="dark"
      density="compact"
      flat
    >
      <v-spacer />
      <v-toolbar-items>
        <slot name="prepend-actions" />
        <UploadButton @click="(content: string) => emit('update:modelValue', content)" />
        <URLButton @click="(content: string) => emit('update:modelValue', content)" />
      </v-toolbar-items>
    </v-toolbar>
</template>

<script setup lang="ts">
import UploadButton from "./UploadButton.vue";
import URLButton from "./URLButton.vue";

const props = defineProps({
    modelValue: { type: String, default: '' },
    title: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])
</script>
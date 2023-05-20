<template>
    <v-navigation-drawer
      :model-value="props.modelValue"
      temporary
      width="900"
      location="right"
      @update:modelValue="(event: boolean) => emit('update:modelValue', event)"
    >
        <v-card flat class="px-4 pt-4">
            <v-toolbar color="primary" density="compact" title="Kyverno Configuration"></v-toolbar>
            <v-card-text class="px-0 py-0">
                <MonacoEditor :value="inputs.config" @update:value="(event: string) => inputs.config = event" :height="500" :theme="editorTheme" language="yaml" />
            </v-card-text>
            <v-card-actions class="px-0 pt-0" style="min-height: 36px!important;">
                <v-btn variant="flat" class="rounded-0 mx-0" width="33.3%" color="warning" @click="() => inputs.config = ''">Clear Config</v-btn>
                <v-btn variant="flat" class="rounded-0 mx-0" width="33.3%" color="primary" @click="() => inputs.config = ConfigTemplate">Load Default Config</v-btn>
                <UploadButton btn-class="rounded-0 mx-0" label="Upload ConfigMap" @click="(event: string) => inputs.config = event" width="33.3%" color="secondary" variant="flat" />
            </v-card-actions>
        </v-card>

        <template v-slot:append>
        <div class="pa-2">
            <v-btn flat color="primary" block @click="() => emit('update:modelValue', false)">Close</v-btn>
        </div>
        </template>
    </v-navigation-drawer>
</template>

<script setup lang="ts">
import { watch } from 'vue';
import { inputs } from '@/store';
import { MonacoEditor } from './Panel';
import { editorTheme } from '@/config';
import { ConfigTemplate } from '@/assets/templates';
import UploadButton from './Panel/UploadButton.vue';
import { loadedConfig } from '@/composables';

const props = defineProps({
  modelValue: { type: Boolean, default: false },
});

watch(() => inputs.config, (config: string) => {
    loadedConfig.value = config
})

const emit = defineEmits(["update:modelValue"]);

</script>
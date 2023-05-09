<template>
    <input type="file" ref="input" style="display: none;" :accept="props.accept" @change="send" />
    <v-btn @click="select" prepend-icon="mdi-upload" :loading="loading" :color="color">from file</v-btn>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps({
    accept: { type: String, default: '.yaml,.yml,text/yaml,application/x-yaml' }
})

const input = ref<HTMLInputElement | null>(null)

const emit = defineEmits(['click'])

const loading = ref<boolean>(false)
const color = ref<string>('white')

const send = (event: Event) => {
    loading.value = true

    const target = (event.target as HTMLInputElement | null)
    if (!target) return;

    const files = target.files
    if (!files || !files.length) return;

    const file = files[0];
    if (!file) return
    
    const reader = new FileReader();
    reader.onload = (res) => {
        emit('click', res?.target?.result)
        loading.value = false;
    };
    reader.onerror = () => {
        color.value = "warning"
        loading.value = false;
    };
    reader.readAsText(file);
}

const select = () => {
    if (!input.value) return

    input.value.click()
}
</script>
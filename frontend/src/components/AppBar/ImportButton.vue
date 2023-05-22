<template>
    <input type="file" ref="input" style="display: none;" :accept="props.accept" @change="send" />
    <v-btn 
        small
        @click="select" 
        prepend-icon="mdi-import" 
        :loading="loading" 
        :color="color" 
        :variant="variant" 
        :block="block"
        :class="btnClass"
    >Import Profiles</v-btn>
    <v-snackbar color="error" :model-value="!!error">{{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { PropType, ref } from "vue";
import { watch } from "vue";
import { importProfiles } from "@/functions/export";

const props = defineProps({
    variant: { type: String as PropType<"outlined" | "text"> },
    block: { type: Boolean },
    btnClass: { type: String },
    accept: { type: String, default: '.yaml,.yml,text/yaml,application/yaml,application/x-yaml' }
})

const input = ref<HTMLInputElement | null>(null)
const loading = ref<boolean>(false)
const color = ref<string | undefined>(undefined)
const error = ref<string>('')

const emit = defineEmits(['finished'])

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
        importProfiles(res?.target?.result?.toString() || '')
            .then(() => {
                emit('finished')
            })
            .catch((err: Error) => {
                error.value = err.message
            }).finally(() => {
                target.value = '';
                loading.value = false;
            })
    };
    reader.onerror = () => {
        error.value = 'failed to load your import file'
        loading.value = false;
    };
    reader.readAsText(file);
}

watch(error, (err) => {
    if (!err) return
    color.value = "error"

    setTimeout(() => {
        error.value = ''
        color.value = undefined
    }, 2500)
})

const select = () => {
    if (!input.value) return

    input.value.click()
}
</script>
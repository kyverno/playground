<template>
    <input type="file" ref="input" style="display: none;" :accept="props.accept" @change="send" />
    <v-btn 
        small
        @click="select" 
        prepend-icon="mdi-upload" 
        :loading="loading" 
        :color="color" 
        :variant="variant" 
        :block="block"
        :class="btnClass"
    >Import</v-btn>
    <v-snackbar color="error" :model-value="!!error">{{ error }}</v-snackbar>
</template>

<script setup lang="ts">
import { PropType, ref } from "vue";
import { createInput } from "@/composables";
import { parse } from 'yaml'
import { ProfileExport } from '@/types';
import { init } from "@/store";
import { watch } from "vue";

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

const importPlayground = async (content: string) => {
    const state: ProfileExport | undefined = parse(content)
    if (!state || typeof state !== 'object') {
        throw new Error('invalid import file')
    }

    if (!Array.isArray(state.profiles)) {
        throw new Error('invalid import file')
    }

    if (!state.profiles.length) return;

    if (state.profiles[0]?.name === 'current-state') {
        const currentState = state.profiles.shift()

        init({
            policy: currentState?.policies,
            resource: currentState?.resources,
            context: currentState?.context,
        })
    }

    state.profiles.filter((p => !!p)).forEach((profile) => {
        if (!profile.name) {
            console.warn('invalid profile, no name found')
            return;
        }

        createInput(profile.name, {
            policy: profile.policies,
            resource: profile.resources,
            context: profile.context
        })
    })
}

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
        importPlayground(res?.target?.result?.toString() || '')
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
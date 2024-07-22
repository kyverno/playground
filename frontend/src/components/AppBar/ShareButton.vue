<template>
  <v-dialog v-model="dialog" width="600px">
    <template v-slot:activator>
      <v-btn
        @click="share"
        prepend-icon="mdi-share"
        :variant="variant"
        :block="block"
        title="create a shareable URL with the current input"
        :class="btnClass"
        id="share-button">
        Share
      </v-btn>
    </template>

    <v-card title="Share Policy">
      <v-divider />
      <v-container>
        <v-row>
          <v-col cols="3" class="my-0 align-self-center"><h3 class="text-h6 d-inline">Include</h3></v-col>
          <v-col class="my-0">
            <v-checkbox hide-details label="Context" v-model="values.context" />
          </v-col>
          <v-col class="my-0">
            <v-checkbox hide-details label="Kyverno Config" v-model="values.config" />
          </v-col>
        </v-row>
        <v-row>
          <v-col><v-text-field label="URL" v-model="url" /></v-col>
        </v-row>
      </v-container>
      <v-card-actions>
        <v-btn @click="dialog = false">Close</v-btn>
        <v-spacer />
        <v-tooltip :model-value="copied" location="top" text="Copied" :open-on-hover="false">
          <template v-slot:activator="{ props }">
            <v-btn variant="tonal" :color="btnColor" @click="copy(url)" :disabled="!isSupported" v-bind="props">Copy URL to Clipboard</v-btn>
          </template>
        </v-tooltip>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useClipboard } from '@vueuse/core'
import { type PropType } from 'vue'
import { btnColor } from '@/config'
import { generateContent, type Config } from '@/functions/share'

defineProps({
  btnClass: { type: String, default: '' },
  variant: { type: String as PropType<'outlined' | 'text'> },
  block: { type: Boolean }
})

const dialog = ref<boolean>(false)
const loading = ref<boolean>(false)
const url = ref<string>('')

const values = reactive<Config>({
  config: false,
  context: false,
  crds: true,
  clusterResources: true,
  exceptions: true,
  vapBindings: true
})

const router = useRouter()

const { copy, copied, isSupported } = useClipboard({ source: url })

const generate = (config: Config) => {
  const content = generateContent(config)

  url.value = `${window.location.origin}${window.location.pathname}${router.resolve({ name: 'home', query: { content } }).href}`
}

const share = () => {
  loading.value = true
  generate(values)
  dialog.value = true
  loading.value = false
}

watch(values, generate)

watch(dialog, (open) => {
  if (open) return

  values.config = false
  values.context = false
})
</script>

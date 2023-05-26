<template>
  <input type="file" ref="input" style="display: none" :accept="props.accept" @change="send" />
  <v-tooltip location="bottom" content-class="no-opacity-tooltip" text="Import from File" v-if="tooltip">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="props" @click="select" icon="mdi-upload" :loading="loading" :color="btnColor" :variant="variant" :width="width" :class="btnClass" />
    </template>
  </v-tooltip>
  <v-btn
    v-else
    v-bind="props"
    @click="select"
    prepend-icon="mdi-upload"
    :loading="loading"
    :color="btnColor"
    :variant="variant"
    :width="width"
    :class="btnClass">
    {{ label }}
  </v-btn>
</template>

<script setup lang="ts">
import { PropType } from 'vue'
import { ref } from 'vue'

const props = defineProps({
  accept: { type: String, default: '.yaml,.yml,text/yaml,application/x-yaml' },
  label: { type: String, default: 'file' },
  width: { type: [String, Number] },
  variant: { type: String as PropType<'flat' | 'text'> },
  color: { type: String },
  tooltip: { type: Boolean, default: true },
  btnClass: { type: String }
})

const input = ref<HTMLInputElement | null>(null)

const emit = defineEmits(['click'])

const loading = ref<boolean>(false)
const btnColor = ref<string | undefined>(props.color)

const send = (event: Event) => {
  loading.value = true

  const target = event.target as HTMLInputElement | null
  if (!target) return

  const files = target.files
  if (!files || !files.length) return

  const file = files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (res) => {
    emit('click', res?.target?.result)
    loading.value = false
  }
  reader.onerror = () => {
    btnColor.value = 'warning'
    loading.value = false
  }
  reader.readAsText(file)
}

const select = () => {
  if (!input.value) return

  input.value.click()
}
</script>

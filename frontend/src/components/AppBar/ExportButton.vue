<template>
  <v-dialog v-model="dialog" width="300px">
    <template v-slot:activator="{ props }">
      <v-btn prepend-icon="mdi-export" :variant="variant" :block="block" class="mr-2 text-left justify-start" v-bind="props">Export Profiles</v-btn>
    </template>
    <v-card title="Select Profiles">
      <v-list select-strategy="classic">
        <v-divider />
        <v-list-item @click="selectAll(all)">
          <template v-slot:prepend>
            <v-list-item-action start class="mr-0">
              <v-checkbox-btn :color="btnColor" v-bind="all"></v-checkbox-btn>
            </v-list-item-action>
          </template>

          <v-list-item-title>Select all</v-list-item-title>
        </v-list-item>
        <v-divider />

        <v-list-item>
          <v-checkbox-btn v-model="current" label="Current State" />
        </v-list-item>
        <v-list-item v-for="profile of list" :key="profile">
          <template v-slot:prepend>
            <v-list-item-action start>
              <v-checkbox-btn :value="profile" :label="profile" v-model="profiles"></v-checkbox-btn>
            </v-list-item-action>
          </template>
        </v-list-item>
      </v-list>
      <v-divider />
      <v-card-text>
        <v-btn block class="mb-3" variant="outlined" :color="btnColor" @click="persist" :loading="loading" prepend-icon="mdi-download">Export</v-btn>
        <v-btn block class="mb-2" variant="outlined" prepend-icon="mdi-close" @click="close">Cancel</v-btn>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { getPersisted } from '@/composables'
import { btnColor } from '@/config'
import { computed } from 'vue'
import { PropType, ref } from 'vue'
import { convertProfiles } from '@/functions/export'

defineProps({
  variant: { type: String as PropType<'outlined' | 'text'> },
  block: { type: Boolean },
  btnClass: { type: String }
})

const profiles = ref<string[]>([])
const dialog = ref<boolean>(false)
const current = ref<boolean>(false)
const loading = ref<boolean>(false)

const all = computed(() => {
  if (current.value && list.value.length === profiles.value.length) {
    return { modelValue: true, indeterminate: undefined }
  }

  if (current.value || profiles.value.length) {
    return { indeterminate: true, modelValue: undefined }
  }

  return { modelValue: false, indeterminate: undefined }
})

const selectAll = ({ modelValue }: { modelValue: boolean | undefined }) => {
  if (modelValue) {
    current.value = false
    profiles.value = []
    return
  }

  current.value = true
  profiles.value = [...list.value]
}

const list = getPersisted()

const close = () => {
  dialog.value = false
  profiles.value = []
  current.value = false
  loading.value = false
}

const persist = () => {
  loading.value = true

  const content = convertProfiles(current.value, profiles.value)

  const a = window.document.createElement('a')
  a.href = window.URL.createObjectURL(new Blob([content], { type: 'application/yaml' }))
  a.download = `playground.yaml`

  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)

  close()
  loading.value = false
}
</script>

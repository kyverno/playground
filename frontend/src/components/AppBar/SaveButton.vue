<template>
  <v-menu location="bottom" :close-on-content-click="false" v-model="menu">
    <template v-slot:activator="{ props }">
      <v-btn prepend-icon="mdi-content-save" :class="btnClass" v-bind="props" :block="block" :variant="variant" id="save-button">Save</v-btn>
    </template>
    <v-list class="py-0">
      <template v-for="item in list" :key="item">
        <v-list-item class="py-0 pl-0">
          <v-btn
            :color="state.name.value === item ? btnColor : undefined"
            prepend-icon="mdi-content-save"
            variant="text"
            block
            @click="persist(item)"
            class="mr-2 text-left justify-start">
            {{ item }}
          </v-btn>
        </v-list-item>
        <v-divider />
      </template>
      <v-divider />
      <v-list-item class="py-0 pl-0">
        <v-dialog v-model="dialog" width="600px" transition="fade-transition">
          <template v-slot:activator="{ props }">
            <v-btn prepend-icon="mdi-plus" variant="text" block v-bind="props" class="mr-2 text-left justify-start">New State</v-btn>
          </template>

          <v-card title="Persist your current Input">
            <v-card-text>
              <v-text-field label="Name" v-model="name" />
            </v-card-text>
            <v-card-actions>
              <v-btn @click="dialog = false">Close</v-btn>
              <v-spacer />
              <v-btn @click="add" :color="btnColor" :disabled="!name">Persist</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-list-item>
      <v-divider />
      <v-list-item class="py-0 pl-0">
        <export-button block variant="text" />
      </v-list-item>
    </v-list>
  </v-menu>
  <v-snackbar v-model="persisted" color="success">Changes persisted</v-snackbar>
</template>

<script setup lang="ts">
import { type PropType, ref, watch } from 'vue'
import { btnColor } from '@/config'
import { updateInput, getPersisted, useState } from '@/composables'
import { computed } from 'vue'
import ExportButton from './ExportButton.vue'
import { inputs } from '@/store'

defineProps({
  variant: { type: String as PropType<'outlined' | 'text'> },
  block: { type: Boolean },
  btnClass: { type: String }
})

const persisted = ref<boolean>(false)
const dialog = ref<boolean>(false)
const menu = ref<boolean>(false)
const name = ref<string>('')

const state = useState()
const persistedList = getPersisted()

const list = computed(() => {
  if (!state.name.value) return persistedList.value

  return [...new Set([state.name.value, ...persistedList.value])]
})

const persist = (name: string) => {
  const persistedInput = updateInput(name, inputs)

  state.update({
    policy: inputs.policy,
    resource: inputs.resource,
    oldResource: inputs.oldResource,
    context: inputs.context,
    config: inputs.config,
    clusterResources: inputs.clusterResources,
    exceptions: inputs.exceptions,
    name: persistedInput.name
  })

  menu.value = false
  persisted.value = true
  setTimeout(() => {
    persisted.value = false
  }, 2000)
}

const add = () => {
  persist(name.value)
  dialog.value = false
}

watch(dialog, (open) => {
  if (open) return

  name.value = ''
})
</script>

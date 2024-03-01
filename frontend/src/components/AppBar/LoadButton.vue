<template>
  <v-menu location="bottom" :close-on-content-click="false" v-model="menu">
    <template v-slot:activator="{ props }">
      <v-btn prepend-icon="mdi-folder" :class="btnClass" v-bind="props" :block="block" :variant="variant" id="load-button">Load</v-btn>
    </template>
    <v-list class="py-0">
      <v-list-item class="py-0 pl-0">
        <v-btn prepend-icon="mdi-open-in-app" variant="text" block @click="loadDefault" class="mr-2 text-left justify-start">Default</v-btn>
      </v-list-item>
      <v-divider />
      <template v-for="item in list" :key="item">
        <v-list-item class="py-0 pl-0">
          <v-btn prepend-icon="mdi-open-in-app" variant="text" block @click="load(item)" class="mr-2 text-left justify-start">{{ item }}</v-btn>
          <template #append>
            <v-btn small class="my-1" variant="flat" @click="remove(item)" icon="mdi-close" title="remove entry" />
          </template>
        </v-list-item>
        <v-divider />
      </template>
      <v-list-item class="py-0 pl-0">
        <import-button block variant="text" btn-class="mr-2 text-left justify-start" />
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { useState } from '@/composables'
import { createInput, removeInput, getPersisted } from '@/composables'
import { setDefaults, init } from '@/store'
import { ref } from 'vue'
import ImportButton from './ImportButton.vue'
import { PropType } from 'vue'

defineProps({
  btnClass: { type: String, default: '' },
  variant: { type: String as PropType<'text' | 'flat' | 'outlined'> },
  block: { type: Boolean, default: false }
})

const menu = ref<boolean>(false)

const list = getPersisted()

const { name: state } = useState()

const loadDefault = () => {
  setDefaults()

  menu.value = false
  state.value = ''
}

const load = (name: string) => {
  const input = createInput(name)

  init({
    policy: input.policy.value,
    resource: input.resource.value,
    oldResource: input.oldResource.value,
    context: input.context.value,
    config: input.config.value,
    exceptions: input.exceptions.value,
    clusterResources: input.clusterResources.value,
    name: input.name,
    vapBindings: input.vapBindings.value,
    imageData: input.imageData.value
  })

  menu.value = false
}

const remove = (name: string) => removeInput(name)
</script>

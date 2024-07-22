<template>
  <v-list class="border py-0" v-model="selections">
    <v-virtual-scroll :items="foundings" :height="441" width="568" :max-height="foundings.length * 49" v-if="foundings.length">
      <template v-slot:default="{ item }">
        <v-list-item @click="() => select(item)">
          <template v-slot:prepend>
            <v-list-item-action start>
              <v-checkbox-btn :model-value="selections" :value="item"></v-checkbox-btn>
            </v-list-item-action>
          </template>

          <v-list-item-title>
            <template v-if="item.namespace">{{ item.namespace }}/{{ item.name }}</template>
            <template v-else>{{ item.name }}</template>
          </v-list-item-title>
        </v-list-item>
        <v-divider />
      </template>
    </v-virtual-scroll>
    <v-list-item v-else>
      <v-list-item-title>No resources found</v-list-item-title>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { type PropType } from 'vue'

type Resource = { namespace?: string; name: string }

const props = defineProps({
  modelValue: { type: Array as PropType<Resource[]>, default: () => [] },
  foundings: { type: Array as PropType<Resource[]>, default: () => [] }
})

const emit = defineEmits(['update:modelValue'])

const selections = computed({
  get() {
    return props.modelValue
  },
  set(val) {
    emit('update:modelValue', val)
  }
})

const select = (res: Resource) => {
  const index = selections.value.findIndex((i) => `${i.namespace}/${i.name}` === `${res.namespace}/${res.name}`)
  if (index === -1) {
    selections.value.push(res)
    return
  }

  selections.value.splice(index, 1)
}
</script>

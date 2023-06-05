<template>
  <v-menu location="top" :close-on-content-click="false" v-model="menu">
    <template #activator="{ props }">
      <v-btn v-bind="props" prepend-icon="mdi-note-text">Templates</v-btn>
    </template>
    <v-list class="py-0 mb-2">
      <template v-if="list.length > 10">
        <v-list-item>
          <v-text-field variant="outlined" density="compact" hide-details placeholder="Search" class="pb-2" v-model="search" />
        </v-list-item>
        <v-divider />
      </template>
      <template v-if="baseTemplate">
        <v-list-item class="py-0 pl-0">
          <v-btn prepend-icon="mdi-open-in-app" variant="flat" block @click="selectDefault" class="mr-2 text-left justify-start">Default</v-btn>
        </v-list-item>
      </template>
      <v-virtual-scroll :items="templates" :height="height" min-width="210">
        <template v-slot:default="{ item }">
          <v-divider />
          <v-list-item class="py-0 pl-0">
            <v-btn prepend-icon="mdi-open-in-app" variant="flat" block @click="loadTemplate(item)" class="mr-2 text-left justify-start">{{ item }}</v-btn>
            <template v-slot:append>
              <v-btn icon="mdi-close" size="small" variant="flat" @click="() => removeTemplate(item)" />
            </template>
          </v-list-item>
        </template>
      </v-virtual-scroll>
      <template v-if="selected">
        <v-divider />
        <v-btn
          :color="btnColor"
          prepend-icon="mdi-content-save"
          class="mr-2 text-left justify-start"
          :disabled="!content"
          block
          height="48"
          variant="text"
          @click="() => add(selected)">
          Update {{ selected }}
        </v-btn>
      </template>
      <v-divider v-if="list.length" />
      <v-dialog v-model="dialog" width="600px" transition="fade-transition">
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" prepend-icon="mdi-content-save" class="mr-2 text-left justify-start" :disabled="!content" block height="48" variant="text">
            Save as Template
          </v-btn>
        </template>

        <v-card title="Persist current content as template">
          <v-card-text>
            <v-text-field label="Name" v-model="name" />
          </v-card-text>
          <v-card-actions>
            <v-btn @click="dialog = false">Close</v-btn>
            <v-spacer />
            <v-btn @click="() => add(name)" :color="btnColor" :disabled="!name">Persist</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { getTemplate, getTemplates, createTemplate } from '@/composables/templates'
import { btnColor } from '@/config'

const props = defineProps({
  baseTemplate: { type: String, required: false },
  panel: { type: String, required: true },
  content: { type: String, required: true }
})

const emit = defineEmits(['select'])
const menu = ref(false)
const search = ref('')
const name = ref('')
const selected = ref('')
const dialog = ref<boolean>(false)

const { list, remove } = getTemplates(props.panel)

const templates = computed(() => {
  if (!search.value) return list.value

  return list.value.filter((i) => i.toLowerCase().search(search.value.toLowerCase()) !== -1)
})

const height = computed(() => {
  const base = list.value.length * 50
  if (base > 350) return 350

  return base
})

const selectDefault = () => {
  emit('select', props.baseTemplate)
  menu.value = false

  setTimeout(() => {
    selected.value = ''
  }, 500)
}

const loadTemplate = (name: string) => {
  const template = getTemplate(props.panel, name)

  emit('select', template.value)
  menu.value = false

  setTimeout(() => {
    selected.value = name
  }, 500)
}

const removeTemplate = (name: string) => {
  const template = getTemplate(props.panel, name)
  remove(name)

  template.value = null
  if (selected.value === name) {
    selected.value = ''
  }
}

const add = (name: string) => {
  createTemplate(props.panel, { name, content: props.content })

  dialog.value = false
  menu.value = false
}

watch(dialog, (open: boolean) => {
  if (open) return

  setTimeout(() => {
    name.value = ''
  }, 500)
})

watch(
  () => props.content,
  (content) => {
    if (content) return

    setTimeout(() => {
      selected.value = ''
    }, 500)
  }
)
</script>

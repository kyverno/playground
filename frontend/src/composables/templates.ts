import { type RemovableRef, useLocalStorage } from '@vueuse/core'
import { type Ref, ref, watch } from 'vue'

export type Template = {
  name: string
  content: string
}

const normalizeName = (name: string) => name.replaceAll(';;', ';').trim()
const convertNames = (names: string) => names.split(';;').filter((l) => !!l)

export const getTemplates = (
  panel: string,
): { list: Ref<string[]>; add: (template: string) => void; remove: (template: string) => void } => {
  const persisted = useLocalStorage<string>(`${panel}:templates:list`, '')

  const list = ref<string[]>(convertNames(persisted.value || ''))

  watch(persisted, (content: string) => {
    list.value = convertNames(content || '')
  })

  const add = (name: string) => {
    persisted.value = [...new Set([...list.value, normalizeName(name)])].join(';;')
  }

  const remove = (name: string) => {
    persisted.value = list.value.filter((l) => l !== name).join(';;')
  }

  return { list, add, remove }
}

export const getTemplate = (panel: string, name: string): RemovableRef<string> => {
  return useLocalStorage<string>(`${panel}:persist:policy:${name}`, null)
}

export const createTemplate = (panel: string, value: Template): Ref<string> => {
  const template = useLocalStorage<string>(`${panel}:persist:policy:${value.name}`, value.content)
  if (template.value !== value.content) {
    template.value = value.content
  }

  const { add } = getTemplates(panel)

  add(value.name)

  return template
}

import { useLocalStorage, usePreferredDark } from '@vueuse/core'
import { watch, computed } from 'vue'
import { type Policy } from './functions/github'
import examples from '../public/tutorials/tutorials.json'

export type Config = {
  editorThemes: { name: string; theme: string }[]
  layoutThemes: string[]
  onboarding: { text: string }
}

type Example = {
  name: string
  color?: string
  url: string
  subgroups?: { name: string; policies: Policy[]; url?: string }[]
  policies?: Policy[]
}

const isDark = usePreferredDark()
export const layoutTheme = useLocalStorage<'light' | 'dark'>(
  'config:layoutTheme',
  isDark.value ? 'dark' : 'light'
)
watch(isDark, (dark: boolean) => {
  layoutTheme.value = dark ? 'dark' : 'light'
})

export const btnColor = computed(() => {
  if (layoutTheme.value === 'dark') return 'secondary'

  return 'primary'
})

export const editorTheme = useLocalStorage('config:editorTheme', 'vs-dark')
export const hideNoMatch = useLocalStorage('config:hideNoMatch', false)
export const showOnboarding = useLocalStorage('onboarding:open', true)

export const options = {
  panels: {
    policyInfo: 'Kyverno Policy',
    resourceInfo: 'Kubernetes resources which get applied to policies',
    contextInfo: 'Context information like admission context, variables, and Kubernetes version',
    exceptionsInfo: 'Configure Kyverno PolicyException Resources',
    vapBindingInfo: 'Configure ValidatingAdmissionPolicyBinding Resources',
    crdInfo: 'Define unknown CRDs you want to use as resource',
    configInfo: 'Configure the Kyverno Engine',
    clusterResourcesInfo:
      'Already existing resources to simulate clone operations or context substitution',
    imageDataInfo: 'Simulate loading of not accessable ImageData'
  },
  onboarding: {
    text: 'Notice: This tool only works with public image registries. No data is gathered, stored, or shared.'
  },
  layoutThemes: ['light', 'dark'],
  editorThemes: [
    { name: 'VS Dark', theme: 'vs-dark' },
    { name: 'VS Light', theme: 'vs' },
    { name: 'HC Black', theme: 'hc-black' },
    { name: 'HC Light', theme: 'hc-light' }
  ],
  examples: examples as Example[]
}

export const useConfig = () => ({
  editorTheme,
  layoutTheme,
  showOnboarding,
  options,
  hideNoMatch
})

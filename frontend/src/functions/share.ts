import { init, inputs } from '@/store'
import * as lzstring from 'lz-string'

export type Config = {
  config: boolean
  crds: boolean
  clusterResources: boolean
  exceptions: boolean
  context: boolean
}

export const generateContent = (config: Config): string => {
  return lzstring.compressToBase64(
    JSON.stringify({
      policy: inputs.policy,
      resource: inputs.resource,
      oldResource: inputs.oldResource,
      ...(config.context ? { context: inputs.context } : {}),
      ...(config.config ? { config: inputs.config } : {}),
      ...(config.crds ? { customResourceDefinitions: inputs.customResourceDefinitions } : {}),
      ...(config.clusterResources ? { clusterResources: inputs.clusterResources } : {}),
      ...(config.exceptions ? { exceptions: inputs.exceptions } : {})
    })
  )
}

export const parseContent = (decoded: string): void => {
  const content = JSON.parse(lzstring.decompressFromBase64(decoded)) as {
    policy: string
    resource: string
    context?: string
    config?: string
    oldResource?: string
    customResourceDefinitions?: string
    clusterResources?: string
    exceptions?: string
  }

  init(content)
}

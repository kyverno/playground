import { init, inputs } from '@/store'
import * as lzstring from 'lz-string'

export const generateContent = (): string => {
  return lzstring.compressToBase64(
    JSON.stringify({
      policy: inputs.policy,
      resource: inputs.resource,
      context: inputs.context,
      config: inputs.config
    })
  )
}

export const parseContent = (decoded: string): void => {
  const content = JSON.parse(lzstring.decompressFromBase64(decoded)) as {
    policy: string
    resource: string
    context: string
    config: string
  }

  init(content)
}

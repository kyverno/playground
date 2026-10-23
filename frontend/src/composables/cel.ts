import { resolveAPI } from '@/utils'
import type { CelEvaluateResponse } from '@/types'

export type CelEvaluateRequest = {
  expression: string
  resource?: string
  oldResource?: string
  namespaceResource?: string
  context?: {
    operation?: string
    username?: string
    groups?: string[]
  }
}

export const useCelAPI = () => {
  const api = resolveAPI()

  const evaluate = (request: CelEvaluateRequest): Promise<CelEvaluateResponse> => {
    return fetch(`${api}/cel`, {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    }).then(async (resp) => {
      if (resp.status > 200) {
        throw new Error(await resp.text())
      }

      return resp.json()
    })
  }

  return { evaluate }
}

export type RuleType = 'Vaidation'
export type RuleStatus = 'fail' | 'pass' | 'warn' | 'error' | 'skip' | 'no match'
export type ErrorReason = 'POLICY_VALIDATION' | 'ERROR'

type RequireOnlyOne<T, Keys extends keyof T = keyof T> = Pick<T, Exclude<keyof T, Keys>> &
  {
    [K in Keys]-?: Required<Pick<T, K>> & Partial<Record<Exclude<Keys, K>, undefined>>
  }[Keys]

export type ErrorResponse = {
  reason: ErrorReason
  error: string
  violations: {
    detail: string
    field: string
    policyName: string
    policyNamespace: string
    type: string
  }[]
}

export type Metadata = {
  name: string
  namespace?: string
  labels?: { [key: string]: string }
  annotations?: { [key: string]: string }
}

export type Policy = {
  apiVersion: string
  kind: string
  metadata: Metadata
}

export type Resource = {
  apiVersion: string
  kind: string
  metadata: Metadata
}

export type Rule = {
  name: string
  ruleType: RuleType
  message: string
  status: RuleStatus
  generatedResource: string
}

export type PolicyResponse = {
  rules: Rule[] | null
}

export type Result = RequireOnlyOne<
  {
    policy?: Policy
    validatingAdmissionPolicy?: Policy
  },
  'policy' | 'validatingAdmissionPolicy'
>

export type Validation = Result & {
  resource: Resource
  policyResponse: PolicyResponse
}

export type Mutation = Result & {
  resource: Resource
  policyResponse: PolicyResponse
  originalResource: string
  patchedResource: string
}

export type Generation = Result & {
  resource: Resource
  policyResponse: PolicyResponse
}

export type EngineResponse = {
  policies: Policy[]
  resources: Resource[]
  validation?: Validation[]
  mutation?: Mutation[]
  imageVerification?: Mutation[]
  generation?: Generation[]
}

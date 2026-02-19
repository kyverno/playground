export type RuleType = 'Vaidation'
export type RuleStatus = 'fail' | 'pass' | 'warn' | 'error' | 'skip' | 'no match'
export type ErrorReason = 'POLICY_VALIDATION' | 'ERROR'

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
  name: string
  namespace?: string
  labels?: { [key: string]: string }
  annotations?: { [key: string]: string }
  mode: string
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
  responseStatus?: {
    code: number
    message: string
  }
}

export type PolicyResponse = {
  rules: Rule[] | null
}

export type Validation = {
  resource: Resource
  policyResponse: PolicyResponse
  policy: Policy
}

export type Mutation = {
  resource: Resource
  policyResponse: PolicyResponse
  originalResource: string
  patchedResource: string
  policy: Policy
}

export type Generation = {
  resource: Resource
  policyResponse: PolicyResponse
  policy: Policy
}

export type EngineResponse = {
  policies?: Policy[]
  resources: Resource[]
  validation?: Validation[]
  deletion?: Validation[]
  mutation?: Mutation[]
  imageVerification?: Mutation[]
  generation?: Generation[]
}

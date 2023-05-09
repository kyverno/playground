export type RuleType = 'Vaidation';
export type RuleStatus = 'fail' | 'pass' | 'warn' | 'error' | 'skip' | 'no match';

export type Metadata = {
    name: string;
    namespace?: string;
    labels?: { [key: string]: string }
    annotations?: { [key: string]: string }
}

export type Policy = {
    apiVersion: string;
    kind: string;
    metadata: Metadata;
}

export type Resource = {
    apiVersion: string;
    kind: string;
    metadata: Metadata;
}

export type Rule = {
    name: string;
    ruleType: RuleType;
    message: string;
    status: RuleStatus;
}

export type PolicyResponse = {
    rules: Rule[] | null
}

export type Validation = {
    resource: Resource;
    policy: Policy;
    policyResponse: PolicyResponse;
}

export type Mutation = {
    resource: Resource;
    policy: Policy;
    policyResponse: PolicyResponse;
    originalResource: string;
    patchedResource: string;
}

export type EngineResponse = {
    Policies: Policy[];
    Resources: Resource[];
    Validation?: Validation[];
    Mutation?: Mutation[];
}

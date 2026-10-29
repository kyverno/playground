package playground_test

import (
	"context"
	"testing"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/playground"
)

func TestRun(t *testing.T) {
	req := &playground.EngineRequest{
		Policies:  "apiVersion: policies.kyverno.io/v1\nkind: ValidatingPolicy\nmetadata:\n  name: check-deployment-labels\n  annotations:\n    policies.kyverno.io/title: Check Deployment Labels\n    policies.kyverno.io/category: Other\n    policies.kyverno.io/severity: medium\nspec:\n  evaluation:\n    background:\n      enabled: true\n    admission:\n      enabled: false\n  matchConstraints:\n    resourceRules:\n    - apiGroups:   [apps]\n      apiVersions: [v1]\n      operations:  [CREATE, UPDATE]\n      resources:   [deployments]\n  variables:\n    - name: environment\n      expression: >-\n        has(object.metadata.labels) && 'env' in object.metadata.labels && object.metadata.labels['env'] == 'prod'\n  validations:\n    - expression: >-\n        variables.environment == true\n      message: >-\n        Deployment labels must be env=prod\n  auditAnnotations:\n  - key: env\n    valueExpression: \"'env:' + object.metadata.labels['env']\"",
		Resources: "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: good-deployment\n  labels:\n    app: nginx\n    env: prod\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:latest\n---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: bad-deployment\n  labels:\n    app: nginx\n    env: testing\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:latest",
	}
	resp, err := playground.Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{
		BuiltInCrds: []string{"argocd", "cert-manager", "prometheus-operator", "tekton-pipeline", "wgpolicyk8s"},
		LocalCrds:   nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

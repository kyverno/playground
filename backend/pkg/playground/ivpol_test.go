package playground_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/playground"
)

func TestRunImageValidatingPolicy(t *testing.T) {
	// The match condition skips the pod, so no registry is contacted, but the
	// image validating policy engine is still built for the request.
	req := &playground.EngineRequest{
		Policies: `apiVersion: policies.kyverno.io/v1
kind: ImageValidatingPolicy
metadata:
  name: ivpol-sample
spec:
  validationActions: [Deny]
  matchConstraints:
    resourceRules:
    - apiGroups: [""]
      apiVersions: [v1]
      operations: [CREATE]
      resources: [pods]
  matchConditions:
  - name: check-prod-label
    expression: has(object.metadata.labels) && has(object.metadata.labels.prod) && object.metadata.labels.prod == 'true'
  matchImageReferences:
  - glob: ghcr.io/*
  attestors:
  - name: cosign
    cosign:
      keyless:
        identities:
        - subject: https://example.com/*
          issuer: https://example.com
  validations:
  - expression: images.containers.map(image, verifyImageSignatures(image, [attestors.cosign])).all(e, e > 0)
    message: failed to verify image signature
`,
		Resources: `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
  labels:
    prod: "false"
spec:
  containers:
  - name: nginx
    image: ghcr.io/kyverno/test-verify-image:signed
`,
	}
	resp, err := playground.Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{})
	require.NoError(t, err)
	require.NotNil(t, resp)
}

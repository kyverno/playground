package utils

import (
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/resource/convert"
	"github.com/kyverno/playground/backend/pkg/resource/loader"
)

func ToPolicyInterface(untyped unstructured.Unstructured) (kyvernov1.PolicyInterface, error) {
	kind := untyped.GetKind()
	if kind == "Policy" {
		policy, err := convert.To[kyvernov1.Policy](untyped)
		if err != nil {
			return nil, err
		}
		return policy, nil
	} else if kind == "ClusterPolicy" {
		policy, err := convert.To[kyvernov1.ClusterPolicy](untyped)
		if err != nil {
			return nil, err
		}
		return policy, nil
	}
	return nil, fmt.Errorf("invalid kind: %s", kind)
}

func LoadPolicies(l loader.Loader, content []byte) ([]kyvernov1.PolicyInterface, error) {
	untyped, err := loader.LoadResources(l, content)
	if err != nil {
		return nil, err
	}
	var policies []kyvernov1.PolicyInterface
	for _, policy := range untyped {
		policy, err := ToPolicyInterface(policy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}
	return policies, nil
}

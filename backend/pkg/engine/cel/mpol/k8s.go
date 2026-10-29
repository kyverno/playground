package mpol

import (
	"context"
	"fmt"
	"slices"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/admissionpolicy"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
)

func K8sProcess(ctx context.Context, dClient dclient.Interface, tcm patch.TypeConverterManager, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, resource, oldResource unstructured.Unstructured, mpols []v1beta1.MutatingPolicyLike, exceptions []*v1beta1.PolicyException) ([]models.Response, error) {
	provider, err := NewProvider(compiler.NewCompiler(), mpols, exceptions, contextProvider)
	if err != nil {
		return nil, err
	}

	eng := engine.NewEngine(provider, utils.NSResolver(dClient), matching.NewMatcher(), tcm, contextProvider)
	request := utils.NewCELRequest(restMapper, contextProvider, params, resource, oldResource)
	results := make([]models.Response, 0)

	scope := engine.Or(engine.ClusteredPolicy(), engine.NamespacedPolicy(resource.GetNamespace()))
	engineResponse, err := eng.Handle(ctx, request, engine.And(scope, engine.NoTargetMatchConstraintPolicy()))
	if err != nil {
		return nil, err
	}

	for _, res := range engineResponse.Policies {
		var patched unstructured.Unstructured
		if engineResponse.PatchedResource != nil {
			patched = *engineResponse.PatchedResource
		}

		response := engineapi.EngineResponse{
			Resource:        *engineResponse.Resource,
			PatchedResource: patched,
			PolicyResponse: engineapi.PolicyResponse{
				Rules: res.Rules,
			},
		}

		results = append(results, models.ConvertResponse(response.WithPolicy(engineapi.NewMutatingPolicyFromLike(res.Policy))))
	}

	policyKeys := eng.MatchedMutateExistingPolicies(ctx, request)
	slices.Sort(policyKeys)
	for _, key := range slices.Compact(policyKeys) {
		policy, err := eng.GetCompiledPolicy(key)
		if err != nil {
			return nil, err
		}
		if !scope(policy.Policy) || engine.NoTargetMatchConstraintPolicy()(policy.Policy) {
			continue
		}

		trigger := admission.NewAttributesRecord(
			resource.DeepCopy(), oldResource.DeepCopy(), resource.GroupVersionKind(), resource.GetNamespace(), resource.GetName(),
			schema.GroupVersionResource(request.Request.Resource), request.Request.SubResource,
			admission.Operation(request.Request.Operation), nil, false, admissionpolicy.NewUser(request.Request.UserInfo),
		)
		targets, err := resolveTargets(ctx, dClient, restMapper, policy, trigger, &request.Request)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve targets for policy %s: %w", key, err)
		}

		// Evaluate only this compiled policy, so each result belongs to one policy
		// and cannot pick up another policy's final patched resource.
		targetEngine := engine.NewEngine(
			targetProvider{Provider: provider, policy: policy}, utils.NSResolver(dClient),
			targetMatcher{client: dClient, constraints: policy.Policy.GetTargetMatchConstraints().MatchResources, ctx: ctx},
			tcm, contextProvider,
		)
		for _, target := range targets {
			original := target.object.DeepCopy()
			attr := admission.NewAttributesRecord(
				target.object.DeepCopy(), nil, target.object.GroupVersionKind(), target.object.GetNamespace(), target.object.GetName(),
				target.resource, target.subresource, admission.Operation(request.Request.Operation), nil, false,
				admissionpolicy.NewUser(request.Request.UserInfo),
			)
			response, err := targetEngine.Evaluate(ctx, attr, request.Request, func(p v1beta1.MutatingPolicyLike) bool {
				return engine.PolicyKey(p) == key && scope(p)
			})
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate target %s for policy %s: %w", target.key(), key, err)
			}
			for _, evaluated := range response.Policies {
				if len(evaluated.Rules) == 0 {
					continue
				}
				patched := original
				if response.PatchedResource != nil {
					patched = response.PatchedResource
				}
				result := engineapi.EngineResponse{
					Resource:        *original,
					PatchedResource: *patched,
					PolicyResponse:  engineapi.PolicyResponse{Rules: evaluated.Rules},
				}
				results = append(results, models.ConvertResponse(result.WithPolicy(engineapi.NewMutatingPolicyFromLike(evaluated.Policy))))
			}
		}
	}
	return results, nil
}

package engine

import (
	"context"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/go-logr/logr"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov1beta1 "github.com/kyverno/kyverno/api/kyverno/v1beta1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/store"
	"github.com/kyverno/kyverno/pkg/background/generate"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	"github.com/kyverno/kyverno/pkg/config"
	kyvernoengine "github.com/kyverno/kyverno/pkg/engine"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	"github.com/kyverno/kyverno/pkg/engine/policycontext"
	"github.com/kyverno/kyverno/pkg/logging"
	"github.com/kyverno/kyverno/pkg/registryclient"
	jsonutils "github.com/kyverno/kyverno/pkg/utils/json"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var log = logging.WithName("playground")

type Processor struct {
	params        *Parameters
	engine        engineapi.Engine
	genController *generate.GenerateController
	config        config.Configuration
	jmesPath      jmespath.Interface
}

func (p *Processor) Run(ctx context.Context, policies []kyvernov1.PolicyInterface, resources []unstructured.Unstructured) (*Results, error) {
	if err := validateParams(p.params, policies); err != nil {
		return nil, err
	}

	response := &Results{}

	for _, resource := range resources {
		// mutate
		for _, policy := range policies {
			result, err := p.mutate(ctx, policy, resource)
			if err != nil {
				return nil, err
			}

			response.Mutation = append(response.Mutation, result)
		}

		// verify images
		for _, policy := range policies {
			result, err := p.verifyImages(ctx, policy, resource)
			if err != nil {
				return nil, err
			}

			response.ImageVerification = append(response.ImageVerification, result)
		}

		// validate
		for _, policy := range policies {
			result, err := p.validate(ctx, policy, resource)
			if err != nil {
				return nil, err
			}

			response.Validation = append(response.Validation, result)
		}

		// generation
		for _, policy := range policies {
			result, err := p.generate(ctx, policy, resource)
			if err != nil {
				return nil, err
			}

			response.Generation = append(response.Generation, result)
		}
	}

	return response, nil
}

func (p *Processor) mutate(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, err
	}

	response := p.engine.Mutate(ctx, policyContext)
	resource = response.PatchedResource

	return ConvertResponse(response), nil
}

func (p *Processor) verifyImages(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, err
	}

	response, verifiedImageData := p.engine.VerifyAndPatchImages(ctx, policyContext)
	// TODO: we apply patches manually because the engine doesn't
	patches := response.GetPatches()
	if !verifiedImageData.IsEmpty() {
		annotationPatches, err := verifiedImageData.Patches(len(resource.GetAnnotations()) != 0, logr.Discard())
		if err != nil {
			return Response{}, err
		}

		// add annotation patches first
		patches = append(annotationPatches, patches...)
	}

	if len(patches) != 0 {
		patch := jsonutils.JoinPatches(patches...)
		decoded, err := jsonpatch.DecodePatch(patch)
		if err != nil {
			return Response{}, err
		}
		options := &jsonpatch.ApplyOptions{SupportNegativeIndices: true, AllowMissingPathOnRemove: true, EnsurePathExistsOnAdd: true}
		resourceBytes, err := resource.MarshalJSON()
		if err != nil {
			return Response{}, err
		}
		patchedResourceBytes, err := decoded.ApplyWithOptions(resourceBytes, options)
		if err != nil {
			return Response{}, err
		}
		if err := response.PatchedResource.UnmarshalJSON(patchedResourceBytes); err != nil {
			return Response{}, err
		}
	}

	resource = response.PatchedResource
	return ConvertResponse(response), nil
}

func (p *Processor) validate(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, err
	}

	response := p.engine.Validate(ctx, policyContext)

	return ConvertResponse(response), nil
}

func (p *Processor) generate(ctx context.Context, policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, resource)
	if err != nil {
		return Response{}, err
	}

	response := p.engine.Generate(ctx, policyContext)
	if len(response.PolicyResponse.Rules) == 0 {
		return ConvertResponse(response), nil
	}

	gr := toGenerateRequest(policy, resource)

	var newRuleResponse []engineapi.RuleResponse
	for _, rule := range response.PolicyResponse.Rules {
		genRes, err := p.genController.ApplyGeneratePolicy(log.V(2), policyContext, gr, []string{rule.Name()})
		if err != nil {
			return Response{}, err
		}
		unstrGenResource, err := p.genController.GetUnstrResource(genRes[0])
		if err != nil {
			return Response{}, err
		}
		newRuleResponse = append(newRuleResponse, *rule.WithGeneratedResource(*unstrGenResource))
	}
	response.PolicyResponse.Rules = newRuleResponse

	return ConvertResponse(response), nil
}

func (p *Processor) newPolicyContext(policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) (*policycontext.PolicyContext, error) {
	context, err := policycontext.NewPolicyContext(
		p.jmesPath,
		resource,
		p.params.Context.Operation,
		&kyvernov1beta1.RequestInfo{
			AdmissionUserInfo: authenticationv1.UserInfo{
				Username: p.params.Context.Username,
				Groups:   p.params.Context.Groups,
			},
			Roles:        p.params.Context.Roles,
			ClusterRoles: p.params.Context.ClusterRoles,
		},
		p.config,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range p.params.Variables {
		err = context.JSONContext().AddVariable(k, v)
		if err != nil {
			return nil, err
		}
	}

	context = context.
		WithPolicy(policy).
		WithNamespaceLabels(p.params.Context.NamespaceLabels)

	return context, nil
}

func validateParams(params *Parameters, policies []kyvernov1.PolicyInterface) error {
	if params == nil {
		return nil
	}

	for _, policy := range policies {
		for _, rule := range policy.GetSpec().Rules {
			for _, variable := range rule.Context {
				if variable.APICall == nil && variable.ConfigMap == nil {
					continue
				}

				if _, ok := params.Variables[variable.Name]; !ok {
					return fmt.Errorf("Variable %s is not defined in the context", variable.Name)
				}
			}
		}
	}

	return nil
}

func toGenerateRequest(policy kyvernov1.PolicyInterface, resource unstructured.Unstructured) kyvernov1beta1.UpdateRequest {
	return kyvernov1beta1.UpdateRequest{
		Spec: kyvernov1beta1.UpdateRequestSpec{
			Type:   kyvernov1beta1.Generate,
			Policy: policy.GetName(),
			Resource: kyvernov1.ResourceSpec{
				Kind:       resource.GetKind(),
				Namespace:  resource.GetNamespace(),
				Name:       resource.GetName(),
				APIVersion: resource.GetAPIVersion(),
			},
		},
	}
}

func NewEngine(cfg config.Configuration, jp jmespath.Interface) (engineapi.Engine, error) {
	rclient, err := registryclient.New(registryclient.WithLocalKeychain())
	if err != nil {
		return nil, err
	}

	store.SetMock(true)
	store.SetRegistryAccess(true)

	return kyvernoengine.NewEngine(
		cfg,
		config.NewDefaultMetricsConfiguration(),
		jp,
		nil,
		rclient,
		store.ContextLoaderFactory(nil),
		nil,
	), nil
}

func NewProcessor(params *Parameters) (*Processor, error) {
	cfg := config.NewDefaultConfiguration(false)
	jp := jmespath.New(cfg)

	engine, err := NewEngine(cfg, jp)
	if err != nil {
		return nil, err
	}

	contr := generate.NewGenerateControllerWithOnlyClient(dclient.NewEmptyFakeClient(), engine)

	return &Processor{
		params:        params,
		engine:        engine,
		genController: contr,
		config:        cfg,
		jmesPath:      jp,
	}, nil
}

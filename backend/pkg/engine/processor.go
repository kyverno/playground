package engine

import (
	"context"
	"fmt"
	"strings"

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
	"github.com/kyverno/kyverno/pkg/engine/mutate/patch"
	"github.com/kyverno/kyverno/pkg/engine/policycontext"
	"github.com/kyverno/kyverno/pkg/registryclient"
	jsonutils "github.com/kyverno/kyverno/pkg/utils/json"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type Processor struct {
	params        *Parameters
	engine        engineapi.Engine
	genController *generate.GenerateController
	config        config.Configuration
	jmesPath      jmespath.Interface
	cluster       bool
}

func (p *Processor) Run(
	ctx context.Context,
	policies []kyvernov1.PolicyInterface,
	resources []unstructured.Unstructured,
	oldResources []unstructured.Unstructured,
) (*Results, error) {
	if !p.cluster {
		if err := validateParams(p.params, policies); err != nil {
			return nil, err
		}
	}

	response := &Results{}

	for i := range resources {
		oldResource := unstructured.Unstructured{}
		newResource := unstructured.Unstructured{}
		if p.params.Context.Operation == kyvernov1.Delete {
			oldResource = resources[i]
		} else if p.params.Context.Operation == kyvernov1.Update {
			// TODO: bounds check
			oldResource = oldResources[i]
			newResource = resources[i]
		} else {
			newResource = resources[i]
		}

		// mutate
		for _, policy := range policies {
			result, res, err := p.mutate(ctx, policy, oldResource, newResource)
			if err != nil {
				return nil, err
			}
			newResource = res
			response.Mutation = append(response.Mutation, result)
		}

		// verify images
		for _, policy := range policies {
			result, res, err := p.verifyImages(ctx, policy, oldResource, newResource)
			if err != nil {
				return nil, err
			}
			newResource = res
			response.ImageVerification = append(response.ImageVerification, result)
		}

		// validate
		for _, policy := range policies {
			result, err := p.validate(ctx, policy, oldResource, newResource)
			if err != nil {
				return nil, err
			}
			response.Validation = append(response.Validation, result)
		}

		// generation
		for _, policy := range policies {
			result, err := p.generate(ctx, policy, oldResource, newResource)
			if err != nil {
				return nil, err
			}
			response.Generation = append(response.Generation, result)
		}
	}

	return response, nil
}

func (p *Processor) mutate(ctx context.Context, policy kyvernov1.PolicyInterface, old, new unstructured.Unstructured) (Response, unstructured.Unstructured, error) {
	policyContext, err := p.newPolicyContext(policy, old, new)
	if err != nil {
		return Response{}, new, err
	}

	response := p.engine.Mutate(ctx, policyContext)

	return ConvertResponse(response), response.PatchedResource, nil
}

func (p *Processor) verifyImages(ctx context.Context, policy kyvernov1.PolicyInterface, old, new unstructured.Unstructured) (Response, unstructured.Unstructured, error) {
	policyContext, err := p.newPolicyContext(policy, old, new)
	if err != nil {
		return Response{}, new, err
	}

	response, verifiedImageData := p.engine.VerifyAndPatchImages(ctx, policyContext)
	// TODO: we apply patches manually because the engine doesn't
	patches := response.GetPatches()
	if !verifiedImageData.IsEmpty() {
		annotationPatches, err := verifiedImageData.Patches(len(new.GetAnnotations()) != 0, logr.Discard())
		if err != nil {
			return Response{}, new, err
		}

		// add annotation patches first
		patches = append(annotationPatches, patches...)
	}

	if len(patches) != 0 {
		patch := jsonutils.JoinPatches(patch.ConvertPatches(patches...)...)
		decoded, err := jsonpatch.DecodePatch(patch)
		if err != nil {
			return Response{}, new, err
		}
		options := &jsonpatch.ApplyOptions{SupportNegativeIndices: true, AllowMissingPathOnRemove: true, EnsurePathExistsOnAdd: true}
		resourceBytes, err := new.MarshalJSON()
		if err != nil {
			return Response{}, new, err
		}
		patchedResourceBytes, err := decoded.ApplyWithOptions(resourceBytes, options)
		if err != nil {
			return Response{}, new, err
		}
		if err := response.PatchedResource.UnmarshalJSON(patchedResourceBytes); err != nil {
			return Response{}, new, err
		}
	}

	return ConvertResponse(response), response.PatchedResource, nil
}

func (p *Processor) validate(ctx context.Context, policy kyvernov1.PolicyInterface, old, new unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, old, new)
	if err != nil {
		return Response{}, err
	}

	response := p.engine.Validate(ctx, policyContext)

	return ConvertResponse(response), nil
}

func (p *Processor) generate(ctx context.Context, policy kyvernov1.PolicyInterface, old, new unstructured.Unstructured) (Response, error) {
	policyContext, err := p.newPolicyContext(policy, old, new)
	if err != nil {
		return Response{}, err
	}

	response := p.engine.Generate(ctx, policyContext)
	if len(response.PolicyResponse.Rules) == 0 {
		return ConvertResponse(response), nil
	}

	gr := toGenerateRequest(policy, new)

	var newRuleResponse []engineapi.RuleResponse
	for _, rule := range response.PolicyResponse.Rules {
		genRes, err := p.genController.ApplyGeneratePolicy(logr.Discard(), policyContext, gr, []string{rule.Name()})
		if err != nil {
			return Response{}, err
		}
		if len(genRes) == 0 {
			continue
		}
		unstrGenResource, err := p.genController.GetUnstrResource(genRes[0])
		if err != nil {
			return Response{}, err
		}

		// cleanup metadata
		if meta, ok := unstrGenResource.Object["metadata"]; ok {
			delete(meta.(map[string]any), "managedFields")
		}

		newRuleResponse = append(newRuleResponse, *rule.WithGeneratedResource(*unstrGenResource))
	}
	response.PolicyResponse.Rules = newRuleResponse

	return ConvertResponse(response), nil
}

func (p *Processor) newPolicyContext(policy kyvernov1.PolicyInterface, old, new unstructured.Unstructured) (*policycontext.PolicyContext, error) {
	resource := old
	if resource.Object == nil {
		resource = new
	}
	userInfo := authenticationv1.UserInfo{
		UID:      "user-123",
		Username: p.params.Context.Username,
		Groups:   p.params.Context.Groups,
		Extra:    nil,
	}
	var oldBytes, newBytes []byte
	if old.Object != nil {
		bytes, _ := old.MarshalJSON()
		oldBytes = bytes
	}
	if new.Object != nil {
		bytes, _ := new.MarshalJSON()
		newBytes = bytes
	}
	gvk := resource.GroupVersionKind()
	gvr := gvk.GroupVersion().WithResource(strings.ToLower(gvk.Kind + "s"))
	context, err := policycontext.NewPolicyContextFromAdmissionRequest(
		p.jmesPath,
		admissionv1.AdmissionRequest{
			UID:                "abc-123",
			Kind:               metav1.GroupVersionKind(gvk),
			Resource:           metav1.GroupVersionResource(gvr),
			SubResource:        "",
			RequestKind:        nil,
			RequestResource:    nil,
			RequestSubResource: "",
			Name:               resource.GetName(),
			Namespace:          resource.GetNamespace(),
			Operation:          admissionv1.Operation(p.params.Context.Operation),
			UserInfo:           userInfo,
			Object: runtime.RawExtension{
				Raw: newBytes,
			},
			OldObject: runtime.RawExtension{
				Raw: oldBytes,
			},
			DryRun:  nil,
			Options: runtime.RawExtension{},
		},
		kyvernov1beta1.RequestInfo{
			AdmissionUserInfo: userInfo,
			Roles:             p.params.Context.Roles,
			ClusterRoles:      p.params.Context.ClusterRoles,
		},
		gvk,
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

func newEngine(
	cfg config.Configuration,
	jp jmespath.Interface,
	client dclient.Interface,
	cmResolver engineapi.ConfigmapResolver,
	exceptionSelector engineapi.PolicyExceptionSelector,
) (engineapi.Engine, error) {
	rclient, err := registryclient.New(registryclient.WithLocalKeychain())
	if err != nil {
		return nil, err
	}

	store.SetMock(true)
	store.SetRegistryAccess(true)

	factory := store.ContextLoaderFactory(nil)
	if cmResolver != nil {
		factory = ContextLoaderFactory(cmResolver)
	}

	return kyvernoengine.NewEngine(
		cfg,
		config.NewDefaultMetricsConfiguration(),
		jp,
		client,
		rclient,
		factory,
		exceptionSelector,
	), nil
}

func NewProcessor(
	params *Parameters,
	kyvernoConfig *corev1.ConfigMap,
	dClient dclient.Interface,
	cmResolver engineapi.ConfigmapResolver,
	exceptionSelector engineapi.PolicyExceptionSelector,
) (*Processor, error) {
	cfg := config.NewDefaultConfiguration(false)
	if kyvernoConfig != nil {
		cfg.Load(kyvernoConfig)
	}

	jp := jmespath.New(cfg)
	cluster := false

	engine, err := newEngine(cfg, jp, dClient, cmResolver, exceptionSelector)
	if err != nil {
		return nil, err
	}

	if dClient == nil {
		dClient = dclient.NewEmptyFakeClient()
	} else {
		cluster = true
	}

	contr := generate.NewGenerateController(dClient, nil, nil, engine, nil, nil, nil, nil, cfg, nil, logr.Discard(), jp)

	return &Processor{
		params:        params,
		engine:        engine,
		genController: contr,
		config:        cfg,
		jmesPath:      jp,
		cluster:       cluster,
	}, nil
}

func ContextLoaderFactory(cmResolver engineapi.ConfigmapResolver) engineapi.ContextLoaderFactory {
	return func(policy kyvernov1.PolicyInterface, rule kyvernov1.Rule) engineapi.ContextLoader {
		inner := engineapi.DefaultContextLoaderFactory(cmResolver)

		return inner(policy, rule)
	}
}

package utils

import (
	"github.com/kyverno/kyverno/pkg/cel/engine"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	"github.com/kyverno/playground/backend/pkg/engine/models"
)

func NewCELRequest(restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, resource, oldResource unstructured.Unstructured) engine.EngineRequest {
	gvk := resource.GroupVersionKind()
	if gvk.Kind == "" {
		return engine.RequestFromJSON(contextProvider, &resource)
	}

	mapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return engine.EngineRequest{}
	}

	if oldResource.GetName() == "" || oldResource.GetKind() == "" {
		oldResource = resource
	}

	return engine.RequestFromAdmission(
		contextProvider,
		admissionv1.AdmissionRequest{
			UID:                "abc-123",
			Kind:               metav1.GroupVersionKind(gvk),
			Resource:           metav1.GroupVersionResource(mapping.Resource),
			SubResource:        "",
			RequestKind:        ptr.To(metav1.GroupVersionKind(gvk)),
			RequestResource:    ptr.To(metav1.GroupVersionResource(mapping.Resource)),
			RequestSubResource: "",
			Name:               resource.GetName(),
			Namespace:          resource.GetNamespace(),
			Operation:          admissionv1.Operation(params.Context.Operation),
			Object:             runtime.RawExtension{Object: &resource},
			OldObject:          runtime.RawExtension{Object: &oldResource},
			UserInfo: authenticationv1.UserInfo{
				UID:      "user-123",
				Username: params.Context.Username,
				Groups:   params.Context.Groups,
				Extra:    nil,
			},
		},
	)
}

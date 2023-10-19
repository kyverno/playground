package cluster

import (
	"context"

	"github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/pkg/client/clientset/versioned"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type policyExceptionSelector struct {
	additional    []*v2beta1.PolicyException
	kyvernoClient versioned.Interface
	namespace     string
}

func (c policyExceptionSelector) List(selector labels.Selector) ([]*v2beta1.PolicyException, error) {
	var exceptions []*v2beta1.PolicyException
	if c.kyvernoClient != nil {
		list, err := c.kyvernoClient.KyvernoV2alpha1().PolicyExceptions(c.namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		if err == nil {
			for i := range list.Items {
				pe := v2beta1.PolicyException(list.Items[i])
				exceptions = append(exceptions, &pe)
			}
		} else if !kerrors.IsNotFound(err) {
			return nil, err
		}
	}
	for _, exception := range c.additional {
		if c.namespace == "" || exception.GetNamespace() == c.namespace {
			exceptions = append(exceptions, exception)
		}
	}
	return exceptions, nil
}

func NewPolicyExceptionSelector(namespace string, client versioned.Interface, exceptions ...*v2beta1.PolicyException) engineapi.PolicyExceptionSelector {
	return policyExceptionSelector{
		additional:    exceptions,
		kyvernoClient: client,
		namespace:     namespace,
	}
}

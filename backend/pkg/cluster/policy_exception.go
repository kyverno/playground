package cluster

import (
	"context"

	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	"github.com/kyverno/kyverno/pkg/client/clientset/versioned"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type policyExceptionSelector struct {
	additional    []*kyvernov2alpha1.PolicyException
	kyvernoClient versioned.Interface
}

func (c policyExceptionSelector) List(selector labels.Selector) ([]*kyvernov2alpha1.PolicyException, error) {
	var exceptions []*kyvernov2alpha1.PolicyException
	if c.kyvernoClient != nil {
		list, err := c.kyvernoClient.KyvernoV2alpha1().PolicyExceptions(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
			LabelSelector: selector.String(),
		})
		if err != nil {
			return nil, err
		}
		for i := range list.Items {
			exceptions = append(exceptions, &list.Items[i])
		}
	}
	exceptions = append(exceptions, c.additional...)
	return exceptions, nil
}

func NewPolicyExceptionSelector(client versioned.Interface, exceptions []*kyvernov2alpha1.PolicyException) engineapi.PolicyExceptionSelector {
	return policyExceptionSelector{
		additional:    exceptions,
		kyvernoClient: client,
	}
}

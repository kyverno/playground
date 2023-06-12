package models

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
)

type Context struct {
	Username        string                       `json:"username"`
	Groups          []string                     `json:"groups"`
	Roles           []string                     `json:"roles"`
	ClusterRoles    []string                     `json:"clusterRoles"`
	Operation       kyvernov1.AdmissionOperation `json:"operation"`
	NamespaceLabels map[string]string            `json:"namespaceLabels"`
	DryRun          bool                         `json:"dryRun"`
}

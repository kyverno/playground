############
# DEFAULTS #
############

KIND_IMAGE           ?= kindest/node:v1.26.3
KYVERNO_VERSION      ?= 3.0.0-alpha.2

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
HELM                               := $(TOOLS_DIR)/helm
HELM_VERSION                       := v3.10.1
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.17.0
TOOLS                              := $(KIND) $(HELM)

$(HELM):
	@echo Install helm... >&2
	@GOBIN=$(TOOLS_DIR) go install helm.sh/helm/v3/cmd/helm@$(HELM_VERSION)

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

.PHONY: install-tools
install-tools: $(TOOLS) ## Install tools

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

.PHONY: codegen-schema-openapi
codegen-schema-openapi: $(KIND) $(HELM) ## Generate openapi schema
	@echo Generate openapi schema... >&2
	@$(KIND) create cluster --name schema --image $(KIND_IMAGE)
	@$(HELM) upgrade --install --wait --timeout 15m --atomic \
  		--version $(KYVERNO_VERSION) \
  		--namespace kyverno --create-namespace \
  		--repo https://kyverno.github.io/kyverno kyverno kyverno
	@kubectl get --raw /openapi/v3/apis/kyverno.io/v1 > ./schemas/openapiv3/schema.json
	@$(KIND) delete cluster --name schema

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

############
# DEFAULTS #
############

KIND_IMAGE           ?= kindest/node:v1.33.1
KIND_NAME            ?= kind
KYVERNO_VERSION      ?= main
KOCACHE              ?= /tmp/ko-cache
USE_CONFIG           ?= standard,no-ingress,in-cluster,all-read-rbac
KUBECONFIG           ?= ""

#############
# VARIABLES #
#############

GIT_SHA             := $(shell git rev-parse HEAD)
TIMESTAMP           := $(shell date '+%Y-%m-%d_%I:%M:%S%p')
GOOS                ?= $(shell go env GOOS)
GOARCH              ?= $(shell go env GOARCH)
REGISTRY            ?= ghcr.io
REPO                ?= kyverno
BACKEND_DIR         := backend
BACKEND_BIN         := $(BACKEND_DIR)/backend
LD_FLAGS            := "-s -w"
LOCAL_PLATFORM      := linux/$(GOARCH)
PLATFORMS           := linux/arm64,linux/amd64
KO_PLATFORMS        := all
KO_TAGS             := $(GIT_SHA)
PLAYGROUND_IMAGE    := playground
REPO_PLAYGROUND     := $(REGISTRY)/$(REPO)/$(PLAYGROUND_IMAGE)
KO_REGISTRY         := ko.local
COMMA               := ,

ifndef VERSION
APP_VERSION         := $(GIT_SHA)
else
APP_VERSION         := $(VERSION)
endif

#########
# TOOLS #
#########

TOOLS_DIR                     := $(PWD)/.tools
CURRENT_DIR                   := $(PWD)
HELM                          := $(TOOLS_DIR)/helm
HELM_VERSION                  := v3.14.0
KIND                          := $(TOOLS_DIR)/kind
KIND_VERSION                  := v0.29.0
KO                            := $(TOOLS_DIR)/ko
KO_VERSION     				  := v0.17.1
HELM_DOCS                     := $(TOOLS_DIR)/helm-docs
HELM_DOCS_VERSION             := v1.14.2
GCI                           := $(TOOLS_DIR)/gci
GCI_VERSION                   := v0.13.4
GOFUMPT                       := $(TOOLS_DIR)/gofumpt
GOFUMPT_VERSION               := v0.6.0
TOOLS                         := $(KIND) $(HELM) $(KO) $(HELM_DOCS) $(GCI) $(GOFUMPT)

$(HELM):
	@echo Install helm... >&2
	@GOBIN=$(TOOLS_DIR) go install helm.sh/helm/v3/cmd/helm@$(HELM_VERSION)

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

$(HELM_DOCS):
	@echo Install helm-docs... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/norwoodj/helm-docs/cmd/helm-docs@$(HELM_DOCS_VERSION)

$(GCI):
	@echo Install gci... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/daixiang0/gci@$(GCI_VERSION)

$(GOFUMPT):
	@echo Install gofumpt... >&2
	@GOBIN=$(TOOLS_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

.PHONY: gci
gci: $(GCI)
	@echo "Running gci"
	@$(GCI) write -s standard -s default -s "prefix(github.com/kyverno/playground/backend)" ./backend

.PHONY: gofumpt
gofumpt: $(GOFUMPT)
	@echo "Running gofumpt"
	@$(GOFUMPT) -w ./backend

.PHONY: fmt
fmt: gci gofumpt

.PHONY: install-tools
install-tools: $(TOOLS) ## Install tools

.PHONY: clean-tools
clean-tools: ## Remove installed tools
	@echo Clean tools... >&2
	@rm -rf $(TOOLS_DIR)

###########
# CODEGEN #
###########

.PHONY: codegen-helm-docs
codegen-helm-docs: ## Generate helm docs
	@echo Generate helm docs... >&2
	@docker run -v ${PWD}/charts:/work -w /work jnorwood/helm-docs:v1.11.0 -s file

.PHONY: codegen-schema-openapi
codegen-schema-openapi: $(KIND) $(HELM) ## Generate openapi schemas (v2 and v3)
	@echo Generate openapi schema... >&2
	@rm -rf ./schemas/openapi
	@mkdir -p ./schemas/openapi/v2
	@mkdir -p ./schemas/openapi/v3/apis/kyverno.io
	@mkdir -p ./schemas/openapi/v3/apis/policies.kyverno.io
	@mkdir -p ./schemas/openapi/v3/apis/admissionregistration.k8s.io
	@$(KIND) create cluster --name schema --image $(KIND_IMAGE) --config ./scripts/config/kind.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/kyverno/kyverno.io_cleanuppolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/kyverno/kyverno.io_clustercleanuppolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/kyverno/kyverno.io_globalcontextentries.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/kyverno/kyverno.io_clusterpolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/kyverno/kyverno.io_policies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/kyverno/kyverno.io_policyexceptions.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/kyverno/kyverno.io_updaterequests.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policies.kyverno.io/policies.kyverno.io_validatingpolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policies.kyverno.io/policies.kyverno.io_policyexceptions.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policies.kyverno.io/policies.kyverno.io_imagevalidatingpolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policies.kyverno.io/policies.kyverno.io_deletingpolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policies.kyverno.io/policies.kyverno.io_generatingpolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policies.kyverno.io/policies.kyverno.io_mutatingpolicies.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policyreport/wgpolicyk8s.io_clusterpolicyreports.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/policyreport/wgpolicyk8s.io_policyreports.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/reports/reports.kyverno.io_clusterephemeralreports.yaml
	@kubectl create -f https://raw.githubusercontent.com/kyverno/kyverno/$(KYVERNO_VERSION)/config/crds/reports/reports.kyverno.io_ephemeralreports.yaml
	@sleep 15
	@kubectl get --raw /openapi/v3/apis/kyverno.io/v1 > ./schemas/openapi/v3/apis/kyverno.io/v1.json
	@kubectl get --raw /openapi/v3/apis/kyverno.io/v2 > ./schemas/openapi/v3/apis/kyverno.io/v2.json
	@kubectl get --raw /openapi/v3/apis/kyverno.io/v2beta1 > ./schemas/openapi/v3/apis/kyverno.io/v2beta1.json
	@kubectl get --raw /openapi/v3/apis/policies.kyverno.io/v1alpha1 > ./schemas/openapi/v3/apis/policies.kyverno.io/v1alpha1.json
	@kubectl get --raw /openapi/v3/apis/admissionregistration.k8s.io/v1 > ./schemas/openapi/v3/apis/admissionregistration.k8s.io/v1.json
	@$(KIND) delete cluster --name schema

.PHONY: codegen-schema-json
codegen-schema-json: codegen-schema-openapi ## Generate json schemas
	@rm -rf ./schemas/json
	@mkdir -p ./schemas/json
	@chmod 777 ./schemas/json
	@docker run --rm --name openapi2jsonschema --mount type=bind,source="$(PWD)"/schemas/openapi/v3,target=/v3 --mount type=bind,source="$(PWD)"/schemas/json,target=/json ghcr.io/fjogeleit/openapi2jsonschema:master /v3/apis/kyverno.io/v1.json --kubernetes --stand-alone --expanded -o /json/v3
	@docker run --rm --name openapi2jsonschema --mount type=bind,source="$(PWD)"/schemas/openapi/v3,target=/v3 --mount type=bind,source="$(PWD)"/schemas/json,target=/json ghcr.io/fjogeleit/openapi2jsonschema:master /v3/apis/kyverno.io/v2beta1.json --kubernetes --stand-alone --expanded -o /json/v3
	@docker run --rm --name openapi2jsonschema --mount type=bind,source="$(PWD)"/schemas/openapi/v3,target=/v3 --mount type=bind,source="$(PWD)"/schemas/json,target=/json ghcr.io/fjogeleit/openapi2jsonschema:master /v3/apis/kyverno.io/v2.json --kubernetes --stand-alone --expanded -o /json/v3
	@docker run --rm --name openapi2jsonschema --mount type=bind,source="$(PWD)"/schemas/openapi/v3,target=/v3 --mount type=bind,source="$(PWD)"/schemas/json,target=/json ghcr.io/fjogeleit/openapi2jsonschema:master /v3/apis/policies.kyverno.io/v1alpha1.json --kubernetes --stand-alone --expanded -o /json/v3
	@docker run --rm --name openapi2jsonschema --mount type=bind,source="$(PWD)"/schemas/openapi/v3,target=/v3 --mount type=bind,source="$(PWD)"/schemas/json,target=/json ghcr.io/fjogeleit/openapi2jsonschema:master /v3/apis/admissionregistration.k8s.io/v1.json --kubernetes --stand-alone --expanded -o /json/v3

.PHONY: codegen-all
codegen-all: codegen-helm-docs codegen-schema-json ## Generate all codegen

.PHONY: verify-schemas
verify-schemas: codegen-schema-json ## Check openapi and json schemas are up to date
	@echo Checking openapi schemas are up to date... >&2
	@git --no-pager diff -- schemas
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen-schema-openapi".' >&2
	@echo 'To correct this, locally run "make codegen-schema-openapi", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- schemas

.PHONY: verify-helm-docs
verify-helm-docs: codegen-helm-docs ## Check Helm charts are up to date
	@echo Checking helm charts are up to date... >&2
	@git --no-pager diff -- charts
	@echo 'If this test fails, it is because the git diff is non-empty after running "make codegen-helm-docs".' >&2
	@echo 'To correct this, locally run "make codegen-helm-docs", commit the changes, and re-run tests.' >&2
	@git diff --quiet --exit-code -- charts

.PHONY: verify-codegen
verify-codegen: verify-helm-docs verify-schemas ## Verify all generated code and docs are up to date

#########
# BUILD #
#########

.PHONY: build-clean
build-clean: ## Clean built files
	@echo Cleaning built files... >&2
	@rm -rf frontend/dist
	@rm -rf backend/backend
	@rm -rf backend/pkg/server/ui/dist
	@rm -rf backend/data/schemas

.PHONY: build-frontend
build-frontend: ## Build frontend
	@echo Building frontend... >&2
	@cp schemas/json/v3/clusterpolicy-kyverno.io-v1.json frontend/src/schemas
	@cp schemas/json/v3/clusterpolicy-kyverno.io-v2beta1.json frontend/src/schemas
	@cp schemas/json/v3/policy-kyverno.io-v1.json frontend/src/schemas
	@cp schemas/json/v3/policy-kyverno.io-v2beta1.json frontend/src/schemas
	@cp schemas/json/v3/policyexception-kyverno.io-v2.json frontend/src/schemas
	@cp schemas/json/v3/policyexception-kyverno.io-v2beta1.json frontend/src/schemas
	@cp schemas/json/v3/validatingadmissionpolicy-admissionregistration-v1.json frontend/src/schemas
	@cp schemas/json/v3/policyexception-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/policyexceptionlist-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/imagevalidatingpolicy-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/imagevalidatingpolicylist-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/deletingpolicy-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/deletingpolicylist-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/mutatingpolicy-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/mutatingpolicylist-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/generatingpolicy-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/generatingpolicylist-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/validatingpolicy-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/validatingpolicylist-policies.kyverno.io-v1alpha1.json frontend/src/schemas
	@cp schemas/json/v3/validatingadmissionpolicybinding-admissionregistration-v1.json frontend/src/schemas
	@cd frontend && npm install && APP_VERSION=$(APP_VERSION) npm run build

.PHONY: build-backend-assets
build-backend-assets: build-frontend ## Build backend assets
	@echo Building backend assets... >&2
	@rm -rf backend/pkg/server/ui/dist && cp -r frontend/dist backend/pkg/server/ui/dist
	@rm -rf backend/data/schemas
	@mkdir -p backend/data/schemas/apis/kyverno.io && cp -r schemas/openapi/v3/apis/kyverno.io/* backend/data/schemas/apis/kyverno.io
	@mkdir -p backend/data/schemas/apis/policies.kyverno.io && cp -r schemas/openapi/v3/apis/policies.kyverno.io/* backend/data/schemas/apis/policies.kyverno.io

.PHONY: build-backend
build-backend: build-backend-assets ## Build backend
	@echo Building backend... >&2
	@cd backend && go mod tidy && go build .

.PHONY: build-all
build-all: build-frontend build-backend ## Build frontend and backend

.PHONY: ko-build
ko-build: $(KO) build-backend-assets ## Build playground image (with ko)
	@echo Build image with ko... >&2
	@cd backend && LDFLAGS=$(LD_FLAGS) KOCACHE=$(KOCACHE) KO_DOCKER_REPO=$(KO_REGISTRY) \
		$(KO) build . --preserve-import-paths --tags=$(KO_TAGS) --platform=$(LOCAL_PLATFORM)

.PHONY: ko-publish
ko-publish: $(KO) ## Build and publish playground image (with ko)
	@echo Publishing image with ko... >&2
	@cd backend && LDFLAGS=$(LD_FLAGS) KOCACHE=$(KOCACHE) KO_DOCKER_REPO=$(REPO_PLAYGROUND) \
		$(KO) build . --bare --tags=$(KO_TAGS) --platform=$(KO_PLATFORMS)

########
# TEST #
########

.PHONY: test-backend
test-backend: ## Test backend
	@echo Testing backend... >&2
	@cd backend && go test ./... -race -coverprofile=coverage.out -covermode=atomic

#######
# RUN #
#######

.PHONY: run
run: build-backend-assets ## Run locally (with connected cluster)
	@echo Run backend... >&2
	@cd backend && go run . \
		--gin-mode=release \
		--gin-log \
		--gin-max-body-size=2097152 \
		--cluster \
		--engine-builtin-crds=argocd \
		--engine-builtin-crds=cert-manager \
		--engine-builtin-crds=prometheus-operator \
		--engine-builtin-crds=tekton-pipeline

.PHONY: run-standalone
run-standalone: build-backend-assets ## Run locally (without connected cluster)
	@echo Run backend... >&2
	@cd backend && go run . \
		--gin-mode=release \
		--gin-log \
		--gin-max-body-size=2097152 \
		--engine-builtin-crds=argocd \
		--engine-builtin-crds=cert-manager \
		--engine-builtin-crds=prometheus-operator \
		--engine-builtin-crds=tekton-pipeline

########
# KIND #
########

.PHONY: kind-create-cluster
kind-create-cluster: $(KIND) ## Create kind cluster
	@echo Create kind cluster... >&2
	@$(KIND) create cluster --name $(KIND_NAME) --image $(KIND_IMAGE) --config ./scripts/config/kind.yaml
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	@sleep 15
	@kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s

.PHONY: kind-delete-cluster
kind-delete-cluster: $(KIND) ## Delete kind cluster
	@echo Delete kind cluster... >&2
	@$(KIND) delete cluster --name $(KIND_NAME)

.PHONY: kind-load
kind-load: $(KIND) ko-build ## Build playground image and load it in kind cluster
	@echo Load playground image... >&2
	@$(KIND) load docker-image --name $(KIND_NAME) ko.local/github.com/kyverno/playground/backend:$(GIT_SHA)

.PHONY: kind-install
kind-install: $(HELM) kind-load ## Build image, load it in kind cluster and deploy playground helm chart
	@echo Install playground chart... >&2
	@$(HELM) upgrade --install kyverno-playground --namespace kyverno-playground --create-namespace --wait ./charts/kyverno-playground \
		--set image.registry=$(KO_REGISTRY) \
		--set image.repository=github.com/kyverno/playground/backend \
		--set image.tag=$(GIT_SHA) \
		$(foreach CONFIG,$(subst $(COMMA), ,$(USE_CONFIG)),--values ./scripts/config/$(CONFIG)/kyverno-playground.yaml)

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

############
# DEFAULTS #
############

KIND_IMAGE           ?= kindest/node:v1.26.3
KIND_NAME            ?= kind
KYVERNO_VERSION      ?= 3.0.0-alpha.2
KOCACHE              ?= /tmp/ko-cache

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
PLAYGROUND_IMAGE    := playground
REPO_PLAYGROUND     := $(REGISTRY)/$(REPO)/$(PLAYGROUND_IMAGE)

KO_REGISTRY         := ko.local
ifndef VERSION
KO_TAGS             := $(GIT_SHA)
else ifeq ($(VERSION),main)
KO_TAGS             := $(GIT_SHA),latest
else
KO_TAGS             := $(GIT_SHA),$(subst /,-,$(VERSION))
endif

#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
HELM                               := $(TOOLS_DIR)/helm
HELM_VERSION                       := v3.10.1
KIND                               := $(TOOLS_DIR)/kind
KIND_VERSION                       := v0.18.0
KO                                 := $(TOOLS_DIR)/ko
KO_VERSION                         := main #e93dbee8540f28c45ec9a2b8aec5ef8e43123966
TOOLS                              := $(KIND) $(HELM) $(KO)

$(HELM):
	@echo Install helm... >&2
	@GOBIN=$(TOOLS_DIR) go install helm.sh/helm/v3/cmd/helm@$(HELM_VERSION)

$(KIND):
	@echo Install kind... >&2
	@GOBIN=$(TOOLS_DIR) go install sigs.k8s.io/kind@$(KIND_VERSION)

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

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
codegen-schema-openapi: $(KIND) $(HELM) ## Generate openapi schemas (v2 and v3)
	@echo Generate openapi schema... >&2
	@rm -rf ./schemas
	@mkdir -p ./schemas/openapi/v2
	@mkdir -p ./schemas/openapi/v3
	@$(KIND) create cluster --name schema --image $(KIND_IMAGE)
	@$(HELM) upgrade --install --wait --timeout 15m --atomic \
  		--version $(KYVERNO_VERSION) \
  		--namespace kyverno --create-namespace \
  		--repo https://kyverno.github.io/kyverno kyverno kyverno
	@kubectl get --raw /openapi/v2 > ./schemas/openapi/v2/schema.json
	@kubectl get --raw /openapi/v3/apis/kyverno.io/v1 > ./schemas/openapi/v3/apis/kyverno.io/v1.json
	@$(KIND) delete cluster --name schema

#########
# BUILD #
#########

.PHONY: build-clean
build-clean: ## Clean built files
	@echo Cleaning built files... >&2
	@rm -rf frontend/dist
	@rm -rf backend/backend
	@rm -rf backend/data/dist
	@rm -rf backend/data/schemas

.PHONY: build-frontend
build-frontend: ## Build frontend
	@echo Building frontend... >&2
	@cd frontend && npm install && npm run build

.PHONY: build-backend-assets
build-backend-assets: build-frontend ## Build backend assets
	@echo Building backend assets... >&2
	@rm -rf backend/data/dist && cp -r frontend/dist backend/data/dist
	@rm -rf backend/data/schemas && cp -r schemas/openapi/v3 backend/data/schemas

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
ko-publish: $(KO) build-backend-assets ## Build and publish playground image (with ko)
	@echo Publishing image with ko... >&2
	@cd backend && LDFLAGS=$(LD_FLAGS) KOCACHE=$(KOCACHE) KO_DOCKER_REPO=$(REPO_PLAYGROUND) \
		$(KO) build . --bare --tags=$(KO_TAGS) --platform=$(KO_PLATFORMS)

.PHONY: docker-build
docker-build: ## Build playground image (with docker)
	@docker buildx build --progress plane --platform $(PLATFORMS) --tag "$(REGISTRY)/$(REPO)/playground:latest" . --build-arg LD_FLAGS=$(LD_FLAGS)

#######
# RUN #
#######

.PHONY: run
run: build-backend-assets ## Run locally
	@echo Run backend... >&2
	@cd backend && go run .

########
# KIND #
########

.PHONY: kind-create-cluster
kind-create-cluster: $(KIND) ## Create kind cluster
	@echo Create kind cluster... >&2
	@$(KIND) create cluster --name $(KIND_NAME) --image $(KIND_IMAGE) --config ./scripts/config/kind.yaml

.PHONY: kind-delete-cluster
kind-delete-cluster: $(KIND) ## Delete kind cluster
	@echo Delete kind cluster... >&2
	@$(KIND) delete cluster --name $(KIND_NAME)

########
# HELP #
########

.PHONY: help
help: ## Shows the available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

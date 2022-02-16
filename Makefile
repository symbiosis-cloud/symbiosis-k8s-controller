GO_VERSION := 1.17.5

DOCKER_REPO ?= registry.symbiosis.host/symbiosiscloud/symbiosis-k8s-controller
PKG ?= symbiosis-cloud/symbiosis-k8s-controller/cmd/symbiosis-k8s-controller
VERSION ?= v0.0.5

.PHONY: compile
compile:
	@echo "Building project"
	@docker run --rm -e GOOS=${OS} -e GOARCH=amd64 -v ${PWD}/:/app -w /app golang:${GO_VERSION}-alpine sh -c 'apk add git && go build -mod=vendor -o cmd/symbiosis-k8s-controller/${NAME} ${PKG}'

.PHONY: build
build:
	@echo "Building docker image $(DOCKER_REPO):$(VERSION)"
	@docker build -t $(DOCKER_REPO):$(VERSION) --platform linux/amd64 cmd/symbiosis-k8s-controller -f cmd/symbiosis-k8s-controller/Dockerfile

.PHONY: push
push:
	@echo "Push docker image $(DOCKER_REPO):$(VERSION)"
	@docker push $(DOCKER_REPO):$(VERSION)

.PHONY: all
all: compile build push

# Test all packages by default
TEST ?= ./...

ifeq ($(shell go env GOOS),darwin)
SEDOPTS = -i ''
else
SEDOPTS = -i
endif

.PHONY: default
default: build

.PHONY: tools
tools:
	cd $(shell go env GOROOT) && go install github.com/hashicorp/go-bindata/go-bindata@latest && go install gotest.tools/gotestsum@latest

.PHONY: proto
proto:
	cd proto/ && buf generate

.PHONY: build
build:
	go install

# The name of Terraform custom provider.
CUSTOM_PROVIDER_NAME ?= terraform-provider-acme
# The url of Terraform provider.
CUSTOM_PROVIDER_URL ?= example.local/myklst/acme

.PHONY: install-local-custom-provider
install-local-custom-provider:
	export PROVIDER_LOCAL_PATH='$(CUSTOM_PROVIDER_URL)'
	go install .
	GO_INSTALL_PATH="$$(go env GOPATH)/bin"; \
	HOME_DIR="$$(ls -d ~)"; \
	mkdir -p  $$HOME_DIR/.terraform.d/plugins/$(CUSTOM_PROVIDER_URL)/0.1.0/linux_amd64/; \
	cp $$GO_INSTALL_PATH/$(CUSTOM_PROVIDER_NAME) $$HOME_DIR/.terraform.d/plugins/$(CUSTOM_PROVIDER_URL)/0.1.0/linux_amd64/$(CUSTOM_PROVIDER_NAME)
	unset PROVIDER_LOCAL_PATH

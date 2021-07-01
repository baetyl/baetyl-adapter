MODULES?=modbus opcua

OUTPUT_MODS:=$(MODULES:%=cmd/%)
IMAGE_MODS:=$(MODULES:%=image/cmd/%)
GO_TEST_FLAGS?=-race -short -covermode=atomic -coverprofile=coverage.out
GO_TEST_PKGS?=$(shell go list ./...)

.PHONY: all $(OUTPUT_MODS)
all: $(OUTPUT_MODS)

$(OUTPUT_MODS):
	@${MAKE} -C $@

.PHONY: image $(IMAGE_MODS)
image: $(IMAGE_MODS)

$(IMAGE_MODS):
	@${MAKE} -C $(patsubst image/%,%,$@) image

.PHONY: test
test: fmt
	@go test ${GO_TEST_FLAGS} ${GO_TEST_PKGS}
	@go tool cover -func=coverage.out | grep total

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	@rm -rf output

GO_TEST_FLAGS?=-race -short -covermode=atomic -coverprofile=coverage.out
GO_TEST_PKGS?=$(shell go list ./...)

.PHONY: all
all: $(SRC_FILES)
	make -C cmd/modbus all

.PHONY: image
image:
	make -C cmd/modbus image

.PHONY: test
test: fmt
	@go test ${GO_TEST_FLAGS} ${GO_TEST_PKGS}
	@go tool cover -func=coverage.out | grep total

.PHONY: fmt
fmt:
	go fmt ./...

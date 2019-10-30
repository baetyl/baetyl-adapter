MODULES?=modbus
PLATFORM_ALL:=linux/amd64 linux/arm64 linux/386 linux/arm/v7

GO_TEST_FLAGS?=
GO_TEST_PKGS?=$(shell go list ./...)

ifndef PLATFORMS
	GO_OS:=$(shell go env GOOS)
	GO_ARCH:=$(shell go env GOARCH)
	GO_ARM:=$(shell go env GOARM)
	PLATFORMS:=$(if $(GO_ARM),$(GO_OS)/$(GO_ARCH)/$(GO_ARM),$(GO_OS)/$(GO_ARCH))
	ifeq ($(GO_OS),darwin)
		PLATFORMS+=linux/amd64
	endif
else ifeq ($(PLATFORMS),all)
	override PLATFORMS:=$(PLATFORM_ALL)
endif

OUTPUT:=output

OUTPUT_MODS:=$(MODULES:%=baetyl-adapter-%)
TEST_MODS:=$(MODULES:%=test/baetyl-adapter-%) #a little tricky to add prefix 'test/' in order to distinguish from OUTPUT_MODS
IMAGE_MODS:=$(MODULES:%=image/baetyl-adapter-%) # a little tricky to add prefix 'image/' in order to distinguish from OUTPUT_MODS

.PHONY: all $(OUTPUT_MODS)
all: $(OUTPUT_MODS)

$(OUTPUT_MODS):
	@make -C $@

.PHONY: image $(IMAGE_MODS)
image: $(IMAGE_MODS)

$(IMAGE_MODS):
	@make -C $(notdir $@) image

$(TEST_MODS):
	@make -C $(notdir $@) test

.PHONY: test
test: $(TEST_MODS)

.PHONY: rebuild
rebuild: clean all

.PHONY: clean
clean:
	@-rm -rf $(OUTPUT)

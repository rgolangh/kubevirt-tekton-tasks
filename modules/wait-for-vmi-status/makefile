all: clean build

CONFIG_FILE ?= ../../configs/wait-for-vmi-status.yaml
include ../../scripts/makefile-snippets/makefile-release.mk
include ../../scripts/makefile-snippets/makefile-common.mk

include ../../scripts/makefile-snippets/makefile-clean.mk # clean
include ../../scripts/makefile-snippets/makefile-golang.mk # lint lint-fix vendor
include ../../scripts/makefile-snippets/makefile-golang-test.mk # test cover junit-report

release-dev:
	$(CONTAINER_ENGINE) build -f build/$(IMAGE_NAME)/Dockerfile -t $(IMAGE) .

release-dev-with-push: release-dev
	$(CONTAINER_ENGINE) push $(IMAGE) $(ARGS)

release: lint test
	$(CONTAINER_ENGINE) build --no-cache --pull -f build/$(IMAGE_NAME)/Dockerfile -t $(MAIN_IMAGE) .

release-with-push: release
	$(CONTAINER_ENGINE) push $(MAIN_IMAGE) $(ARGS)


.PHONY: \
	all \
	release-dev \
	release-dev-with-push \
	release \
	release-with-push

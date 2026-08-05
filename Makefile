.PHONY: release release-patch release-minor release-major main

BUMP ?= patch

# Compute next tag from latest semver tag
LATEST_TAG := $(shell git tag --sort=-v:refname | head -1)
MAJOR := $(shell echo $(LATEST_TAG) | sed 's/v//' | cut -d. -f1)
MINOR := $(shell echo $(LATEST_TAG) | sed 's/v//' | cut -d. -f2)
PATCH := $(shell echo $(LATEST_TAG) | sed 's/v//' | cut -d. -f3)

ifeq ($(BUMP),major)
  NEXT_TAG := v$(shell echo $$(($(MAJOR)+1))).0.0
else ifeq ($(BUMP),minor)
  NEXT_TAG := v$(MAJOR).$(shell echo $$(($(MINOR)+1))).0
else
  NEXT_TAG := v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))
endif

release:
	@if [ -n "$$(git status --porcelain)" ]; then echo "error: uncommitted changes"; exit 1; fi
	@echo "tagging $(NEXT_TAG) (was $(LATEST_TAG))"
	git tag $(NEXT_TAG)
	git push origin HEAD $(NEXT_TAG)
	@echo "released $(NEXT_TAG)"

release-patch: ; @$(MAKE) release BUMP=patch
release-minor: ; @$(MAKE) release BUMP=minor
release-major: ; @$(MAKE) release BUMP=major

main:
	git push origin dev:main

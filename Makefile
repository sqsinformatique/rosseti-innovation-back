SHELL=/bin/bash -o pipefail

.PHONY: init-stable
init-stable:
		GO111MODULE=on go mod vendor

.PHONY: build-stable
build-stable:
		@TAG=$$(git describe --tags 2>/dev/null; true); \
		if [[ $$TAG == *"-g"* ]] || [[ $$TAG == "" ]]; then \
	        GITVER=$$(git rev-parse --abbrev-ref HEAD 2>/dev/null; true); \
	    else \
	        GITVER=$$TAG; \
	    fi; \
		BUILDTIME=$$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null); \
		HASH=$$(git rev-parse HEAD 2>/dev/null); \
		HASH=$${HASH:0:8}; \
		if [[ $$GITVER == "HEAD" ]]; then \
			GITVER=$$(git name-rev $$HASH 2>/dev/null; true); \
			GITVER=$${GITVER:24}; \
		fi; \
		PACKAGE=$$(go list -mod vendor 2>/dev/null); \
		echo Build $$PACKAGE $$GITVER $$BUILDTIME $$HASH; \
		GO111MODULE=on go install -mod vendor -ldflags "\
				-X $$PACKAGE/internal/cmd.version=$$GITVER \
				-X $$PACKAGE/internal/cmd.builded=$$BUILDTIME \
				-X $$PACKAGE/internal/cmd.hash=$$HASH" \

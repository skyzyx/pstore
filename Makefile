all:
	@cat Makefile | grep : | grep -v PHONY | grep -v @ | sed 's/:/ /' | awk '{print $$1}' | sort

#-------------------------------------------------------------------------------

.PHONY: build
build:
	go build -v -ldflags="-s -w" -o pstore main.go

.PHONY: optimize
optimize:
	upx --brute pstore

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: lint
lint:
	gometalinter.v2 ./...

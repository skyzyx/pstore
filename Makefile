all:
	@cat Makefile | grep : | grep -v PHONY | grep -v @ | sed 's/:/ /' | awk '{print $$1}' | sort

#-------------------------------------------------------------------------------

.PHONY: build
build:
	go build -v -ldflags="-s -w" -o pstore main.go

.PHONY: optimize
optimize:
	upx --brute pstore

.PHONY: fat
fat:
	@ # https://hackernoon.com/a-story-of-a-fat-go-binary-20edc6549b97
	eval $$(go build -a -work 2>&1) && find $$WORK -type f -name "*.a" | xargs -I% du -hxs "%" | sort -rh | sed -e s:$${WORK}/::g

.PHONY: package
package:
	for platform in windows linux darwin; do \
		CGO_ENABLED=0 GOOS=$$platform GOARCH=amd64 go build -a -v -ldflags="-s -w" -o pstore-$$platform-amd64 main.go && \
		upx --brute pstore-$$platform-amd64; \
	done;

	mv pstore-windows-amd64 pstore-windows-amd64.exe && \
		zip pstore-windows-amd64.zip pstore-windows-amd64.exe && \
		rm -f pstore-windows-amd64.exe

	for platform in linux darwin; do \
		tar cvf pstore-$$platform-amd64.tar pstore-$$platform-amd64 && bzip2 -9 pstore-$$platform-amd64.tar && rm -f pstore-$$platform-amd64; \
	done;

.PHONY: lint
lint:
	gometalinter.v2 ./...

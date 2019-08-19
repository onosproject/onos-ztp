export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_ZTP_VERSION := latest
ONOS_ZTP_DEBUG_VERSION := debug
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-ztp ./cmd/onos-ztp
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/onos-ztp-debug ./cmd/onos-ztp

test: # @HELP run the unit tests and source code validation
test: build deps linters license_check
	go test github.com/onosproject/onos-ztp/pkg/...
	go test github.com/onosproject/onos-ztp/cmd/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./build/bin/coveralls-coverage

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/)"

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-ztp \
		-w /go/src/github.com/onosproject/onos-ztp \
		--entrypoint pkg/northbound/proto/compile-protos.sh \
		onosproject/protoc-go:stable

update-deps: # @HELP pull updated dependencies
	go get github.com/onosproject/onos-topo
	go get github.com/onosproject/onos-config

onos-ztp-base-docker: # @HELP build onos-ztp base Docker image
onos-ztp-base-docker: update-deps
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-ztp-base:${ONOS_ZTP_VERSION}
	@rm -rf vendor

onos-ztp-docker: onos-ztp-base-docker # @HELP build onos-ztp Docker image
	docker build . -f build/onos-ztp/Dockerfile \
		--build-arg ONOS_ZTP_BASE_VERSION=${ONOS_ZTP_VERSION} \
		-t onosproject/onos-ztp:${ONOS_ZTP_VERSION}

onos-ztp-debug-docker: onos-ztp-base-docker # @HELP build onos-ztp Docker debug image
	docker build . -f build/onos-ztp-debug/Dockerfile \
		--build-arg ONOS_ZTP_BASE_VERSION=${ONOS_ZTP_VERSION} \
		-t onosproject/onos-ztp:${ONOS_ZTP_DEBUG_VERSION}

images: # @HELP build all Docker images
images: build onos-ztp-docker onos-ztp-debug-docker

kind: # @HELP build Docker images and add them to the currently configured kind cluster
kind: images
	@if [ `kind get clusters` = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/onos-ztp:${ONOS_ZTP_DEBUG_VERSION}

all: build images


clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-ztp/onos-ztp ./cmd/dummy/dummy

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '

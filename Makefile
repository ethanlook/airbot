BUILD_CHANNEL?=local
OS=$(shell uname)
VERSION=v1.12.0
GIT_REVISION = $(shell git rev-parse HEAD | tr -d '\n')
TAG_VERSION?=$(shell git tag --points-at | sort -Vr | head -n1)
CGO_LDFLAGS=""
GO_BUILD_LDFLAGS = -ldflags "-X 'main.Version=${TAG_VERSION}' -X 'main.GitRevision=${GIT_REVISION}'"
TOOL_BIN = bin/gotools/$(shell uname -s)-$(shell uname -m)

.PHONY: default
default: build-module

.PHONY: tool-install
tool-install:
	GOBIN=`pwd`/$(TOOL_BIN) go install \
		github.com/edaniels/golinters/cmd/combined \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/AlekSi/gocov-xml \
		github.com/axw/gocov/gocov \
		gotest.tools/gotestsum \
		github.com/rhysd/actionlint/cmd/actionlint

.PHONY: lint
lint: tool-install
	go mod tidy
	export pkgs="`go list -f '{{.Dir}}' ./... | grep -v /proto/`" && echo "$$pkgs" | xargs go vet -vettool=$(TOOL_BIN)/combined
	GOGC=50 $(TOOL_BIN)/golangci-lint run -v --fix --config=./golangci.yaml

.PHONY: test
test:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...

bin/buf bin/protoc-gen-go bin/protoc-gen-grpc-gateway bin/protoc-gen-go-grpc:
	GOBIN=$(shell pwd)/bin go install \
		github.com/bufbuild/buf/cmd/buf \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: build-proto
build-proto: proto/v1/api.proto bin/buf bin/protoc-gen-go bin/protoc-gen-grpc-gateway bin/protoc-gen-go-grpc
	PATH="$(shell pwd)/bin" buf generate

.PHONY: build
build: 
	mkdir -p bin && rm bin/airbot; go build $(GO_BUILD_LDFLAGS) -o bin/airbot ./module/main.go

.PHONY: buildarm
buildarm: 
	mkdir -p bin && rm bin/airbot; GOARCH="arm64" GOOS="linux" go build $(GO_BUILD_LDFLAGS) -o bin/airbot ./module/main.go

.PHONY: run
run:
	go run $(GO_BUILD_LDFLAGS) ./module/main.go

.PHONY: buildstatic
buildstatic: 
	mkdir -p bin && CGO_ENABLED=0 CGO_LDFLAGS=${CGO_LDFLAGS} go build $(GO_BUILD_LDFLAGS) -o bin/module ./module/main.go

package: build

.PHONY: clean
clean: 
	rm -rf bin

.PHONY: package
package: build
	tar -czf bin/airbot.tar.gz bin/airbot routes

.PHONY: mock-build-arm
mock-build-arm: 
	docker build -t mockrobot . -f ./mockrobot/Dockerfile.aarch64

.PHONY: mock-build-x86
mock-build-x86: 
	docker build -t mockrobot . -f ./mockrobot/Dockerfile.x86_64

.PHONY: mock-run
mock-run: 
	./mockrobot/runImage.sh

APP=SmartLyu-go-core
VERSION=`cat VERSION`
GITBRANCH=`git symbolic-ref --short -q HEAD`
GITREVISION=`git log -n1 --format=%H`
BUILDUSER=luyu151111@gamil.com
BUILDDATE=`date "+%Y-%m-%d %H:%M:%S"`
PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

## mod config
mod-config:
	@go env -w GO111MODULE='on'
	@go env -w GOPROXY='https://goproxy.cn,direct'

mod-tidy: mod-config
	@go mod tidy

mod-vendor: mod-config
	@go mod download
	@go mod vendor

## swag build
swag-build:
	@swag init -parseDependency --parseInternal  --parseDepth=1

## list: list packages and go files
list:
	@echo $(PACKAGES)
	@echo $(VETPACKAGES)
	@echo $(GOFILES)

## vet: static check
vet: mod-vendor
	@go vet $(VETPACKAGES)

## fmt: format code
fmt:
	@gofmt -s -w $(GOFILES)

## fmt-check: format result check
fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
	  echo "Please run 'make fmt' and commit the result:"; \
	  echo "$${diff}"; \
	  exit 1; \
	fi;

## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean
	@if [ -d target ] ; then rm -rf target ; fi

## bindata: package static resources
bindata:
	@go-bindata -o asset/asset.go -pkg=asset views/ docs/swagger.*

## test: runs go test with default values
test: bindata
	@go test ./... -v

## build: build the application to registry
build: bindata clean
	@go build -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' -X 'github.com/prometheus/common/version.Branch=$(GITBRANCH)' -X 'github.com/prometheus/common/version.Revision=$(GITREVISION)'" -o target/$(APP)

## run: runs go run main.go
run:
	@go run -ldflags="-X 'github.com/prometheus/common/version.Version=$(VERSION)' -X 'github.com/prometheus/common/version.BuildUser=$(BUILDUSER)' -X 'github.com/prometheus/common/version.BuildDate=$(BUILDDATE)' " main.go

## compile: build the application of different operating systems
compile:
	@# 32-Bit Systems
	@# FreeBDS
	@GOOS=freebsd GOARCH=386 go build -o target/$(APP)-freebsd-386 main.go
	@# MacOS
	@GOOS=darwin GOARCH=386 go build -o target/$(APP)-darwin-386 main.go
	@# Linux
	@GOOS=linux GOARCH=386 go build -o target/$(APP)-linux-386 main.go
	@# Windows
	@GOOS=windows GOARCH=386 go build -o target/$(APP)-windows-386 main.go
	@# 64-Bit
	@# FreeBDS
	@GOOS=freebsd GOARCH=amd64 go build -o target/$(APP)-freebsd-amd64 main.go
	@# MacOS
	@GOOS=darwin GOARCH=amd64 go build -o target/$(APP)-darwin-amd64 main.go
	@# Linux
	@GOOS=linux GOARCH=amd64 go build -o target/$(APP)-linux-amd64 main.go
	@# Windows
	@GOOS=windows GOARCH=amd64 go build -o target/$(APP)-windows-amd64 main.go

help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' Makefile | column -t -s ':' |  sed -e 's/^/ /'

## before push code: vet、fmt、fmt-check
before-push-code: list mod-tidy vet fmt fmt-check

## all: execut test、build、docker-build、docker-push targets
all: vet fmt fmt-check test build

.PHONY: before-push-code all

# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
       BUILD_VERBOSE = $(V)
endif

ifeq ($(BUILD_VERBOSE),1)
       Q =
else
       Q = @
endif

.PHONY: all
all: build

.PHONY: mod
mod:
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: vet
vet: mod
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go vet -mod=vendor ./...

.PHONY: format
format: vet
	gofmt -s -l -w main.go

.PHONY: build
build: format
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a .

.PHONY: clean
clean:
	rm -rf ./image-labels-test


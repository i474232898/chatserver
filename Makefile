APP := $(shell basename $(shell git remote get-url origin) | sed 's/\.git$$//')
REGISTRY := i474232898

VERSION := $(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS := linux
TARGETARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

.PHONY: build format lint clean

format:
	gofmt -s -w ./

test:
	go test ./...

lint:
	golangci-lint run

build: format
	CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH)	go build -o bin/chatserver \
		-ldflags "-X=github.com/$(REGISTRY)/$(APP)/cmd.appVersion=$(VERSION)" \
		cmd/chatserver/main.go

image:
	docker build -t $(REGISTRY)/$(APP):$(VERSION)-$(TARGETOS)-$(TARGETARCH) .

push:
	docker push $(REGISTRY)/$(APP):$(VERSION)-$(TARGETOS)-$(TARGETARCH)

clean:
	rm -rf bin/chatserver

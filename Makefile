GOFILES := $(shell find cmd internal -name '*.go' -type f | sort)
IMAGE ?= k8s-top-exporter
TAG ?= latest

.PHONY: fmt test docker

fmt:
	test -z "$$(gofmt -l $(GOFILES))"

test:
	go test ./...

docker:
	docker build -t $(IMAGE):$(TAG) .

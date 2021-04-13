.PHONY: all
  
all: build test
IMG="sysdig/dockerfile-benchmarker"
VERSION=$(shell cat version)

test:
	@echo "+ $@"
	go test ./...
build:
	@echo "+ $@"
	./scripts/build-binary
build-image:
	@echo "+ $@"
	docker build -f container/Dockerfile -t ${IMG}:${VERSION} .
push-image:
	@echo "+ $@"
	docker push ${IMG}:${VERSION}
	docker tag ${IMG}:${VERSION} ${IMG}:latest
	docker push ${IMG}:latest
	

PROJECT  := expired
VERSION  := 0.0.1
REGISTRY := superorbital
TAG      := ${REGISTRY}/${PROJECT}:${VERSION}

build: bin/expired

bin/expired:
	go build -o $@ ./cmd/main.go

container-build:
	docker build . -t ${TAG}

container-push:
	docker push ${TAG}

project=graphite-ng

version=0.1.0

PACKAGES=${project} carbon-es
GOM_GROUPS=test,carbon

export PATH := $(abspath ./_vendor/bin):$(PATH)
GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

GOM=$(if $(TRAVIS),$(HOME)/gopath/bin/gom,gom)

all: ${project}

clean:
	rm -rf _vendor
	rm -f executor-*.go

run:
	$(GOPATH)/bin/${project}

$(project): deps
	rm -f executor-*.go
	$(GOM) build -i -o $(GOPATH)/bin/${project} -a

deps:
	$(GOM) -groups=$(GOM_GROUPS) install

lint: deps
	$(GOM) exec go fmt ./...
	$(GOM) exec go vet -x ./...
	$(GOM) exec golint .
	$(foreach p, $(PACKAGES), $(GOM) exec golint ./$(p)/.; )

test: deps lint
	(test -f coverage.out && "$(TRAVIS)" == "true") || \
		$(GOM) exec go test -covermode=count -coverprofile=coverage.out .
	(test -f ${project}/${project}.coverprofile && "$(TRAVIS)" == "true") || \
		$(GOM) exec ginkgo -cover=true $(shell pwd)/.

deb:
	dpkg-buildpackage -uc -b -d -tc

build-container:
	docker build -t ${project} $(shell pwd)/.

run-container:
	docker run -i -p 8080:8080 -v $(shell pwd)/:/usr/share/go/src/github.com/${project}/${project} -t ${project} /bin/bash

.PHONY: all clean run deps test
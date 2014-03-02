VERSION := 0.1

ifndef GOROOT
$(error GOROOT is not set)
endif
ifndef GOPATH
$(error GOPATH is not set)
endif

$(info GOROOT=$(GOROOT))
$(info GOPATH=$(GOPATH))

GO := $(GOROOT)/bin/go

format :
	$(info Formatting sources)
	$(GO) fmt eu.dominiek/dtop

build-all: build-arm5 build-x64

build-arm5 :
	$(info Building for ARM5)
	GOARCH=arm GOARM=5 $(GO) build -o bin/arm5/dtop eu.dominiek/dtop

build-x64 :
	$(info Building for x64)
	$(GO) build -o bin/x64/dtop eu.dominiek/dtop

dist-all: build-all dist-arm5 dist-x64

dist-arm5 :
	$(info Packaging ARM5 distribution)
	mkdir -p dist
	tar cf dist/dtop-$(VERSION).arm5.tar.gz README.md scripts/install.sh static bin/arm5/dtop

dist-x64 :
	$(info Packaging x64 distribution)
	mkdir -p dist
	tar cf dist/dtop-$(VERSION).x64.tar.gz README.md scripts/install.sh static bin/arm5/dtop
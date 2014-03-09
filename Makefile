VERSION := 0.1

ifndef GOROOT
$(warning GOROOT is not set, setting default (/opt/go).)
GOROOT := /opt/go/
endif
ifndef GOPATH
$(warning GOPATH is not set, setting current folder.)
GOPATH := $(shell pwd)
endif

$(info GOROOT=$(GOROOT))
$(info GOPATH=$(GOPATH))

GO := GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GOROOT)/bin/go

format :
	$(info Formatting sources)
	$(GO) fmt eu.dominiek/dtop

clean :
	rm dist -rf
	rm bin -rf

build-all : build-arm5 build-x64

build-arm5 :
	$(info Building for ARM5)
	GOARCH=arm \
	GOARM=5 \
	cd $(GOROOT)/src \
	source make.bash \
	cd - \
	$(GO) build -o bin/arm5/dtop eu.dominiek/dtop

build-x64 :
	$(info Building for x64)
	$(GO) build -o bin/x64/dtop eu.dominiek/dtop

dist-all : build-all dist-arm5 dist-x64

dist-arm5 :
	$(info Packaging ARM5 distribution)
	mkdir -p dist
	tar cf dist/dtop-$(VERSION).bin.arm5.tar.gz README.md scripts/install.sh static bin/arm5/dtop

dist-x64 :
	$(info Packaging x64 distribution)
	mkdir -p dist
	tar cf dist/dtop-$(VERSION).bin.x64.tar.gz README.md scripts/install.sh static bin/x64/dtop

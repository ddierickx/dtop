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

GO := $(GOROOT)/bin/go

build-all : build-linux-x64 build-linux-386 build-linux-arm5 build-darwin-386 build-darwin-x64
dist-all : build-all dist-linux-deb-x64

build-linux-arm5 :
	$(info Building for linux arm5)
	GOROOT=$(GOROOT) ; GOPATH=$(GOPATH) ; GOARM=5 ; scripts/compile linux arm bin/linux-arm5/dtop

build-linux-386 :
	$(info Building for linux x86)
	GOROOT=$(GOROOT) ; GOPATH=$(GOPATH) ; scripts/compile linux 386 bin/linux-x86/dtop

build-linux-x64 :
	$(info Building for linux x64)
	GOROOT=$(GOROOT) ; GOPATH=$(GOPATH) ; scripts/compile linux amd64 bin/linux-x64/dtop

build-darwin-386 :
	$(info Building for darwin 386)
	GOROOT=$(GOROOT) ; GOPATH=$(GOPATH) ; scripts/compile darwin 386 bin/darwin-386/dtop

build-darwin-x64 :
	$(info Building for darwin x64)
	GOROOT=$(GOROOT) ; GOPATH=$(GOPATH) ; scripts/compile darwin amd64 bin/darwin-x64/dtop

dist-linux-deb-x64 :
	$(info Packaging linux x64 deb distribution)
	mkdir -p dist
	#tar cf dist/dtop-$(VERSION).linux.x64.tar.gz README.md scripts/install.sh static bin/x64/dtop

format :
	$(info Formatting sources)
	$(GO) fmt eu.dominiek/dtop

clean :
	rm dist -rf
	rm bin -rf
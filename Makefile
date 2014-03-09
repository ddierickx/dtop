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

build-all : build-linux-amd64
dist-all : build-all dist-linux-deb-amd64 dist-linux-rpm-amd64

build-linux-amd64 :
	$(info Building for linux amd64)
	scripts/compile $(GOROOT) $(GOPATH) linux amd64 bin/linux-amd64/dtop

dist-linux-deb-amd64 :
	$(info Packaging linux amd64 deb distribution)
	cp bin/linux-amd64/dtop /usr/bin/ -f
	mkdir -p /usr/local/share/dtop
	cp static /usr/local/share/dtop -rf
	mkdir -p dist
	fpm --provides dtop -s dir -t deb -n dtop -v $(VERSION) -p dist/dtop_VERSION-linux-ARCH.deb /usr/bin/dtop /usr/local/share/dtop/static

dist-linux-rpm-amd64 :
	$(info Packaging linux amd64 rpm distribution)
	cp bin/linux-amd64/dtop /usr/bin/ -f
	mkdir -p /usr/local/share/dtop
	cp static /usr/local/share/dtop -rf
	mkdir -p dist
	fpm --provides dtop -s dir -t rpm -n dtop -v $(VERSION) -p dist/dtop_VERSION-linux-ARCH.rpm /usr/bin/dtop /usr/local/share/dtop/static


format :
	$(info Formatting sources)
	$(GO) fmt eu.dominiek/dtop

clean :
	rm dist -rf
	rm bin -rf
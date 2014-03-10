DIST_NAME := dtop
DIST_VERSION := 0.1
DIST_DESCRIPTION := A monitoring tools that brings htop like functionality and more to the webbrowser.
DIST_VENDOR := Dominique Dierickx
DIST_MAINTAINER := d.dierickx@gmail.com
DIST_URL := https://github.com/ddierickx/dtop

DIST_CATEGORY := monitoring

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

build-all : build-linux-amd64 build-linux-386
dist-all : build-all dist-linux-deb-amd64 dist-linux-rpm-amd64 dist-linux-deb-386 dist-linux-rpm-386

build-linux-386 :
	$(info Building for linux 386)
	scripts/compile $(GOROOT) $(GOPATH) linux 386 bin/linux-386/dtop

build-linux-amd64 :
	$(info Building for linux amd64)
	scripts/compile $(GOROOT) $(GOPATH) linux amd64 bin/linux-amd64/dtop

dist-linux-deb-amd64 :
	$(info Packaging linux amd64 deb distribution)
	cp bin/linux-amd64/dtop /usr/bin/ -f
	mkdir -p /usr/local/share/dtop
	cp static /usr/local/share/dtop -rf
	cp scripts/debian/dtopd -f /etc/init.d/
	mkdir -p dist
	fpm -s dir \
		-t deb \
		--provides "$(DIST_NAME)" \
		--name "$(DIST_NAME)" \
		--description "$(DIST_DESCRIPTION)" \
		--version "$(DIST_VERSION)" \
		--vendor "$(DIST_VENDOR)" \
		--maintainer "$(DIST_MAINTAINER)" \
		--url "$(DIST_URL)" \
		--category "$(DIST_CATEGORY)" \
		--vendor "$(DIST_VENDOR)" \
		--architecture "x86_64" \
		--package "dist/dtop_VERSION-linux-ARCH.deb" \
		--after-install "scripts/debian/run" \
		"/usr/bin/dtop" \
		"/usr/local/share/dtop/static" \
		"/etc/init.d/dtopd"

dist-linux-rpm-amd64 :
	$(info Packaging linux amd64 rpm distribution)
	cp bin/linux-amd64/dtop /usr/bin/ -f
	mkdir -p /usr/local/share/dtop
	cp static /usr/local/share/dtop -rf
	cp scripts/rhel/dtopd -f /etc/init.d/
	mkdir -p dist
	fpm -s dir \
		-t rpm \
		--depends "redhat-lsb-core" \
		--provides "$(DIST_NAME)" \
		--name "$(DIST_NAME)" \
		--description "$(DIST_DESCRIPTION)" \
		--version "$(DIST_VERSION)" \
		--vendor "$(DIST_VENDOR)" \
		--maintainer "$(DIST_MAINTAINER)" \
		--url "$(DIST_URL)" \
		--category "$(DIST_CATEGORY)" \
		--vendor "$(DIST_VENDOR)" \
		--architecture "amd64" \
		--package "dist/dtop_VERSION-linux-amd64.rpm" \
		--after-install "scripts/rhel/run" \
		"/usr/bin/dtop" \
		"/usr/local/share/dtop/static" \
		"/etc/init.d/dtopd"

dist-linux-rpm-386 :
	$(info Packaging linux 386 rpm distribution)
	cp bin/linux-386/dtop /usr/bin/ -f
	mkdir -p /usr/local/share/dtop
	cp static /usr/local/share/dtop -rf
	cp scripts/rhel/dtopd -f /etc/init.d/
	mkdir -p dist
	fpm -s dir \
		-t rpm \
		--depends "redhat-lsb-core" \
		--provides "$(DIST_NAME)" \
		--name "$(DIST_NAME)" \
		--description "$(DIST_DESCRIPTION)" \
		--version "$(DIST_VERSION)" \
		--vendor "$(DIST_VENDOR)" \
		--maintainer "$(DIST_MAINTAINER)" \
		--url "$(DIST_URL)" \
		--category "$(DIST_CATEGORY)" \
		--vendor "$(DIST_VENDOR)" \
		--package "dist/dtop_VERSION-linux-i386.rpm" \
		--after-install "scripts/rhel/run" \
		"/usr/bin/dtop" \
		"/usr/local/share/dtop/static" \
		"/etc/init.d/dtopd"

dist-linux-deb-386 :
	$(info Packaging linux 386 deb distribution)
	cp bin/linux-386/dtop /usr/bin/ -f
	mkdir -p /usr/local/share/dtop
	cp static /usr/local/share/dtop -rf
	cp scripts/debian/dtopd -f /etc/init.d/
	mkdir -p dist
	fpm -s dir \
		-t deb \
		--provides "$(DIST_NAME)" \
		--name "$(DIST_NAME)" \
		--description "$(DIST_DESCRIPTION)" \
		--version "$(DIST_VERSION)" \
		--vendor "$(DIST_VENDOR)" \
		--maintainer "$(DIST_MAINTAINER)" \
		--url "$(DIST_URL)" \
		--category "$(DIST_CATEGORY)" \
		--vendor "$(DIST_VENDOR)" \
		--architecture "i386" \
		--package "dist/dtop_VERSION-linux-ARCH.deb" \
		--after-install "scripts/debian/run" \
		"/usr/bin/dtop" \
		"/usr/local/share/dtop/static" \
		"/etc/init.d/dtopd"

format :
	$(info Formatting sources)
	$(GO) fmt eu.dominiek/dtop

clean :
	rm dist -rf
	rm bin -rf
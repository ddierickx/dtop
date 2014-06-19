DIST_NAME := dtop
DIST_VERSION := 0.3
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
dist-all : build-all dist-linux-deb-amd64 dist-linux-rpm-amd64 dist-linux-deb-386 # dist-linux-rpm-386

build-linux-386 :
	$(info Building for linux 386)
	scripts/compile $(GOROOT) $(GOPATH) linux 386 bin/linux-386/dtop

build-linux-amd64 :
	$(info Building for linux amd64)
	scripts/compile $(GOROOT) $(GOPATH) linux amd64 bin/linux-amd64/dtop

dist-linux-deb-amd64 :
	$(info Packaging linux amd64 deb distribution)
	cp bin/linux-amd64/dtop /usr/bin/ -f
	mkdir -p /var/dtop
	mkdir -p /etc/dtop
	mkdir -p dist
	
	cp static /var/dtop -rf
	cp scripts/debian/dtopd -f /etc/init.d/
	cp conf/distribution.json /etc/dtop/default.json

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
		--before-remove "scripts/debian/stop" \
		--after-install "scripts/debian/run" \
		"/usr/bin/dtop" \
		"/var/dtop" \
		"/etc/init.d/dtopd" \
		"/etc/dtop/default.json"

dist-linux-rpm-amd64 :
	$(info Packaging linux amd64 rpm distribution)
	cp bin/linux-amd64/dtop /usr/bin/ -f
	mkdir -p /var/dtop
	mkdir -p /etc/dtop
	mkdir -p dist

	cp static /var/dtop -rf
	cp scripts/rhel/dtopd -f /etc/init.d/
	cp conf/distribution.json /etc/dtop/default.json

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
		--before-remove "scripts/rhel/stop" \
		--after-install "scripts/rhel/run" \
		"/usr/bin/dtop" \
		"/var/dtop" \
		"/etc/init.d/dtopd" \
		"/etc/dtop/default.json"

# disabled because i can't find the correct 'architecture' param...
dist-linux-rpm-386 :
	$(info Packaging linux 386 rpm distribution)
	cp bin/linux-386/dtop /usr/bin/ -f
	mkdir -p /var/dtop
	mkdir -p /etc/dtop
	mkdir -p dist
	
	cp static /var/dtop -rf
	cp scripts/rhel/dtopd -f /etc/init.d/
	cp conf/distribution.json /etc/dtop/default.json

	fpm -s dir \
		--verbose \
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
		--before-remove "scripts/rhel/stop" \
		--after-install "scripts/rhel/run" \
		"/usr/bin/dtop" \
		"/var/dtop" \
		"/etc/init.d/dtopd" \
		"/etc/dtop/default.json"
	setarch i386 rpmbuild dist/dtop_$(DIST_VERSION)-linux-i386.rpm

dist-linux-deb-386 :
	$(info Packaging linux 386 deb distribution)
	cp bin/linux-386/dtop /usr/bin/ -f
	mkdir -p /var/dtop
	mkdir -p /etc/dtop
	mkdir -p dist
	
	cp static /var/dtop -rf
	cp scripts/debian/dtopd -f /etc/init.d/
	cp conf/distribution.json /etc/dtop/default.json

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
		--before-remove "scripts/debian/stop" \
		--after-install "scripts/debian/run" \
		"/usr/bin/dtop" \
		"/var/dtop" \
		"/etc/dtop/default.json" \
		"/etc/init.d/dtopd"

run : test
	$(info Compiling and running dtop)
	$(GO) build eu.dominiek/dtop
	./dtop -c conf/default.json

test : format
	$(info Running tests)
	$(GO) test fmt eu.dominiek/dtop

format :
	$(info Formatting sources)
	$(GO) fmt eu.dominiek/dtop

clean :
	rm dist -rf
	rm bin -rf

# dtop #

A web-based implementation of the excellent [htop project](http://htop.sourceforge.net).

## Introduction ##

dtop is a tool that tries to deliver a large part of htop's realtime functionality to the webbrowser. It is currently very much a work in progress so please feel free to add features and improve!

![Image](/doc/screenshot1.png?raw=true)
![Image](/doc/screenshot2.png?raw=true)

## Features ##

*   using Golang so no runtime requirements needed after compilation
*   only uses resources when the UI is accessed (and not a lot)
*   works on big and small, from Raspberry Pi to multi-processor systems
*	authentication support
*	cpu usage per core
*   basic info: hostname, distro, kernel
*	service status
*	memory- and swap-usage (overall and per process)
*	uptime
*	load avg
*	disk free/usage
*	users
*	basic process search functionality
*   mobile / tablet compatible ui

## Installation ##

**Note:** packaging in rpm and deb is currently in the works so there may be some issues since I only tested them on CentOS and Ubuntu. Nonetheless, please give them a try and report any issues you might experience here on github.

* [dtop_0.2-linux-amd64.deb](http://www.do-jo.be/dtop/dtop_0.2-linux-amd64.deb)
* [dtop_0.2-linux-amd64.rpm](http://www.do-jo.be/dtop/dtop_0.2-linux-amd64.rpm)

Note: i386 packages are disabled because of an rpm building issue.

To install the deb packages:

	sudo dpkg -i dtop_0.2-linux-amd64.deb

For the rpm's (this requires redhat-lsb-core but should be installed automatically):

	sudo yum localinstall dtop_0.2-linux-amd64.rpm

The configuration files are stored in `/etc/dtop/`

The static files are stored in `/var/dtop/static`

The binary file (`dtop`) is in `/usr/bin/dtop`

You can easily create these packages yourself by following the steps in Development => Distribution.

If you only require the binary, you can use (Golang compiler required):

```
go build -o drun ./cmd/drun.go
```

### Run adhoc

You can start dtop as follows:

	./bin/linux-amd64/dtop -c conf/default.json

Then point your webbrowser at http://localhost:12345 and you should see the dashboard.

The default configuration file (default.json) looks as follows:

	{
	    "Name": "my server",
	    "Description": "my description",
	    "StaticFolder": "/usr/local/share/dtop/static",
	    "Port": 12345,
	    "Services": [
	        { "Name": "ufw" },
	        { "Name": "ssh" },
	        { "Name": "cups" },
	        { "Name": "tomcat6" },
	        { "Name": "puppet" },
	        { "Name": "nfs-kernel-server" },
	        { "Name": "postgresql" }
	    ]
	 }

If you'd like to require authentication, you need to adapt the configuration file and add the users property. Here's an example:

	{
	    "Name": "my server",
	    "Description": "my description",
	    "StaticFolder": "/var/dtop/static",
	    "Port": 12345,
	    "Users": [
	    	{ "Username": "hodor", "Password": "hodorhodor" }
	    ],
	    "Services": [
	        { "Name": "ufw" },
	        { "Name": "ssh" },
	        { "Name": "cups" },
	        { "Name": "tomcat6" },
	        { "Name": "puppet" },
	        { "Name": "nfs-kernel-server" },
	        { "Name": "postgresql" }
	    ]
	 }

When configured like this, a login dialog will appear before the dashboard is shown.

![Image](/doc/screenshot3.png?raw=true)

Please note that the password is not hashed yet...

### Run as a daemon

Assuming you have the correct packages for your OS (see development => distribution). The deamon, dtopd, will be installed by default along with the package.

The service can then be controlled as usual:

	sudo service dtopd status|start|stop|restart

## Development ##

dtop is developed in [Go](http://golang.org) so you need the go compiler.

Clone the repo and cd into it.

	git clone https://github.com/ddierickx/dtop
	cd dtop

A Makefile is available with several helpful commands: test, format, build-all, dist-all, ... Note that they require that GOPATH and GOROOT are set.

### Distribution ###

To ease the creation of packages for different OSes a Vagrantfile is available to create these packages. The VM can be found in the scripts/vm-distribution folder. If you have Vagrant and Virtualbox installed you can simply execute:

	vagrant up

This should create RPM and DEB packages for the i386 and amd64 in the dist folder.

## TODO ##

*	Javascript refactoring / cleaning
*	OSX support
*   re-add swap usage
*	user defined sorting iso. cpu percentage
*   a nice favicon/logo
*	hash passwords
*	processtree
*	process kill feature
*	htop like keyboard shortcuts
*	implement Pri, Ni, Virt, Res, Shr, S, Time in process list.

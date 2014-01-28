# dtop #

A web-based implementation of the excellent [htop project](http://htop.sourceforge.net).

## introduction ##

dtop is a tool that tries to deliver a large part of htop's realtime functionality to the webbrowser. It is currently very much a work in progress so please feel free to add features and improve!

## usage ##

dtop is developed in [Go](http://golang.org) so you need the go compiler.

Clone the repo and cd into it.

> git clone https://github.com/ddierickx/dtop

> cd dtop

### run adhoc

Set your GOPATH.

> export GOPATH=`pwd`

Compile the package.

> go build eu.dominiek/dtop

Run dtop.

> ./dtop

Then point your webbrowser at http://localhost:12345 and you should see:

![Image](/screenshot.png?raw=true)

### run as a daemon

> cd scripts

Now run the installation script which will install dtop to `/opt/dtop` and register the daemon within init.d.

> sudo ./install-dtopd.sh

You can start, stop, status and restart the service as usual:

> sudo service dtopd status

## features ##

*	cpu
*	memory- and swap-usage (overall and per process)
*	uptime
*	load avg
*	users
*	basic search functionality

## todo ##

*	user defined sorting iso. cpu percentage
*	processtree
*	basic (?) authentication
*	process kill feature
*	htop like keyboard shortcuts
*	implement Pri, Ni, Virt, Res, Shr, S, Time in process list.
*	find solution for location of the 'static' folder; because the CWD of dtop should be its residing folder atm.

# OutOfMyWay Time Tracker

*** (12/11/19) In active development, it may not work exactly as described, but it's getting close ***

A minimalist time tracker.  The primary purposes of this tool are:

1. Help a user track time and tasks without getting in the way of flow
2. Provide a simple, extendable reporting interface to help transfer tasks to an external system

Any contributions to the tool will not impact these purposes.

Secondary goals are:

* Support Linux, Windows, and MacOS
* Do not require network connectivity to fulfill primary functionality
* Stay as lightweight as practical
* (future) Provide an extensible API to the reporting functions
* (future) Sychronize your timesheet to a backend to allow seamless sharing between devices

This time-tracking tool was inspired by the [Ultimate Time Tracker](https://github.com/larose/utt).

Transfering quick tasks from this tool into an external tracking system (ie: Workday) is still largely manual. Since this is usually a requirement for enterprise users, streamlining this process will be a core part of reaching version 1.0.

# Prerequisites

## For running

* The latest release of Omw
* If you want to use the progressive web app, a recent web browser (we target Chrome for now but others may work)

### Getting Started

The program has a command-line and HTTP-accessible API builtin.

To use the command-line interface, run the program `omw` without any arguments to get help.

To use the network API:

1. Run `omw server` and note the URL returned
2. Visit the Omw PWA URL and install the Chrome extension **coming soon*

## For developing

* Go 1.11+

### Building

[Mage](https://magefile.org) provides build, install, and packing functions.

To run all build steps and install:

`go run mage.go`

To run all build stages:

`go run mage.go build`

To run an individual build stage:

```
go run mage.go buildgo
go run mage.go buildpkg
```

# Architecture

Omw is a simple, stateless, time tracker application, in that there is never a running clock in the background.  It only adds a task with the current timestamp to a text file log, and then compares adjacent timestamps to generate reports.  The timesheet is written line-by-line and stored in the default home directory as returned by the `go-homedir` package under `.local/share/omw/omw.log`.

The binary provides a command-line interface and a Go Gorilla Mux HTTP server providing a REST-ish API.  An flock() package provides an interface to operating system file locking.

# References

* [Ultimate Time Tracker](https://github.com/larose/utt)

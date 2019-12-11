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
* A recent version of Chrome

### Getting Started

The program has a command-line and browser interface builtin.

To use the command-line interface, run the program `omw` without any arguments to get help.

To use the browser interface:

1. Run `omw server` and note the URL returned
2. Browse to the URL
3. Enter a command.  If you need help, toggle the slider or enter the command `?`. A successful command should quickly execute its function and then minimize the window.
4. The calendar is hidden by default but should automatically fetch recent events in the background.  To display them, check the `report` checkbox, or enter the `r` command.
5. Use the `r` and `l` commands as a keyboard-driven way to move the calendar focus. **(not yet implemented)** 

Optionally:
**not yet released**

1. If you want to add support for global hotkeys to Omw, install the Chrome app omw-hotkeys
2. After setting your preferred default action key combo, press it, and the Omw browser interface should appear

## For developing

* Go 1.11+
* NodeJS v11.14.0+

### Building

[Mage](https://magefile.org) provides build, install, and packing functions.

To run all build steps and install:

`go run mage.go`

To run all build stages:

`go run mage.go build`

To run an individual build stage:

```
go run mage.go buildui
go run mage.go buildgo
go run mage.go buildpkg
```

# Architecture

Omw is a simple, stateless, time tracker application, in that there is never a running clock in the background.  It only adds a task with the current timestamp to a text file log, and then compares adjacent timestamps to generate reports.  The timesheet is written line-by-line and stored in the default home directory as returned by `go-homedir` under `.local/share/omw/omw.log`.

The binary provides a command-line interface as well as an embedded web application, using the `statik` package, accessible via a Go HTTP server providing a REST-ish API.  An flock() package provides an interface to operating system file locking.

I chose to leverage Chrome to provide cross-platform global (always available) keyboard shortcuts, which proved difficult to do elegantly across Windows, Linux, and MacOS, using only Go dependencies. The global shortcuts are the critical component of the tool to keep you "in flow".  Chrome provides the ability to register global keyboard shortcuts with the Chrome Extensions [Commands API](https://developer.chrome.com/extensions/commands). 

A [LitElement](https://lit-element.polymer-project.org/) web application provides the browser interface.

# Building

Planning to move this into [Mage](https://github.com/magefile/mage) to handle the npm/polymer build commands, and investigating [xgo](https://github.com/karalabe/xgo) to handle the CGO cross-compilation necessary for the Robotgo Hook library.  Until then, running the `go build` step on the desired operating system is probably easier.

# References

* [PWA Starter Kit](https://github.com/Polymer/pwa-starter-kit)
* [Ultimate Time Tracker](https://github.com/larose/utt)
* [go-homedir](https://github.com/mitchellh/go-homedir)
* [statik](https://github.com/rakyll/statik)

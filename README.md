# OutOfMyWay Time Tracker

*** (5/26/19) In active development, it may not work exactly as described ***

A minimalist time tracker.  The primary purposes of this tool are:

1. Help a user track time and tasks without getting in the way of flow
2. Provide a simple, extendable reporting interface to help transfer tasks to an external system

Any contributions to the tool will not impact these purposes.

Secondary goals are:

* Support Linux, Windows, and MacOS
* Do not require network connectivity to fulfill primary functionality
* Stay as lightweight as practical
* (future) Provide an extensible API to the reporting functions

This time-tracking tool was inspired by the [Ultimate Time Tracker](https://github.com/larose/utt).  The original version (and this version at the moment) calls `utt` as an external dependency.

Transfering quick tasks from this tool into an external tracking system (ie: Workday) is still largely manual. Since this is usually a requirement for enterprise users, streamlining this process will be a core part of reaching version 1.0.

# Prerequisites

## For running

* A recent version of Chrome

Lorca uses the locally installed version of Chrome with remote debugging protocol to provide the UI.

* Python 2 or 3
* [Python UTT](https://github.com/larose/utt/)

The Python dependencies will be removed before 1.0.

### Getting Started

1. Run the program, it should display the main app window in a few seconds.
2. Enter a command.  If you need help, toggle the slider or enter the command `?`. A successful command should quickly execute its function and then minimize the window.
3. Hit the global shortcut (`<right shift> + <left shift>`) on your keyboard to bring up the window again throughout the day.
4. Use the `r` and `l` commands to generate reports that you can use to capture time in an external tool.

## For developing

* NodeJS v11.14.0+
* Polymer-cli (installed with `npm install`)
* Go 1.11+
* [Robotgo dependencies](https://github.com/go-vgo/robotgo#requirements)
    * Ubuntu 19.04 seems to also need `apt install libxkbcommon-x11-dev`

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

This tool is being rewritten using Go, [Lorca](https://github.com/zserge/lorca), and [LitElement](https://lit-element.polymer-project.org/).

The keyboard shortcut, a critical part of the tool to keep you "in flow", is provided by [RobotGo Hook](https://github.com/robotn/gohook/), which also relies on [cgo](https://golang.org/cmd/cgo/).  The global keyboard shortcut is the primary reason for developing the application with Lorca and Go, instead of as a pure Progressive Web App.  Chrome provides the ability to register global keyboard shortcuts with the Chrome Extensions [Commands API](https://developer.chrome.com/extensions/commands), but I find this design nice to work with for the time being, and it makes working with the local filesystem a piece of cake.  Google announce support for Desktop PWAs starting in Chrome 73 and that [shortcuts would be added soon](https://developers.google.com/web/progressive-web-apps/desktop#whats_next).  Because the UI is written with PWA in mind, it should be relatively easy to take advantage of this in the future if it seems like a good choice.

For now, Go and Lorca provide the "server", interfacing with the host operating system and controlling the instance of Chrome using its remote debugging protocol.  [LitElement](https://lit-element.polymer-project.org/) provides the user interface.  Lorca provides a way for Go to execute Javascript inside the instance of Chrome, and also expose functions in Javascript that Chrome can execute to access functions on the OS.

# Building

Planning to move this into [Mage](https://github.com/magefile/mage) to handle the npm/polymer build commands, and investigating [xgo](https://github.com/karalabe/xgo) to handle the CGO cross-compilation necessary for the Robotgo Hook library.  Until then, running the `go build` step on the desired operating system is probably easier.

# History

The original version was written with Python and [Remi](https://github.com/dddomodossola/remi/tree/master/remi).lElectron would've worked, but Lorca gives us nice Go cross-platform capabilities and is much less bloated.

# References

* [PWA Starter Kit](https://github.com/Polymer/pwa-starter-kit)
* [WiredJS Web Components](https://wiredjs.com)

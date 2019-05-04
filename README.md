# OutOfMyWay Time Tracker

*** (5/4/19) In active development, it may not work ***

A minimalist time tracker.  The primary purposes of this tool are two-fold:

1. Help a user track time and tasks without getting in the way of flow
2. Provide a simple reporting interface to help transfer tasks to an external system

Any contributions to the tool will not impact these purposes.

Secondary goals are:

* Support Linux, Windows, and MacOS
* Do not require network connectivity to fulfill primary functionality
* Stay as lightweight as practical
* (future) Provide an extensible API to the reporting functions

The time-tracking process was inspired by [UTT](https://github.com/larose/utt), and originally called it as an external dependency.

Transfering quick tasks from this tool into an external tracking system (ie: Workday) is still largely manual. Since this is usually a requirement for enterprise users, streamlining this process will be a core part of reaching version 1.0.

# Prerequisites

## For running

* A recent version of Chrome

Lorca uses the locally installed version of Chrome with remote debugging protocol to provide the UI.

## For developing

* Go 1.11+
* [Robotgo dependencies](https://github.com/go-vgo/robotgo#requirements)

# Getting Started

1. Run the program
2. Hit the global shortcut (currently `<CTRL>-<SHIFT>-<;>`) on your keyboard to bring up the window
3. If you need help, toggle the slider

# Architecture

This tool is being rewritten using Go, [Lorca](https://github.com/zserge/lorca), and [LitElement](https://lit-element.polymer-project.org/).

The keyboard shortcut, a critical part of the tool to keep you "in flow", is provided by [RobotGo](https://github.com/go-vgo/robotgo), which also relies on [cgo](https://golang.org/cmd/cgo/).  The global keyboard shortcut is the primary reason for developing the application with Lorca and Go, instead of as a pure Progressive Web App.  Chrome provides the ability to register global keyboard shortcuts with the Chrome Extensions [Commands API](https://developer.chrome.com/extensions/commands), but I find this design nice to work with for the time being, and it makes working with the local filesystem a piece of cake.  Because the UI is written with PWA in mind, it should be relatively easy to change this decision in the future.

Go and Lorca provide the "server", interfacing with the host operating system and controlling the instance of Chrome using its remote debugging protocol.  LitElement provides the user interface.  Lorca provides a way for Go to execute Javascript inside the instance of Chrome, and also expose functions in Javascript that Chrome can execute to access functions on the OS.

# History

The original version was written with Python and [Remi](https://github.com/dddomodossola/remi/tree/master/remi).lElectron would've worked, but Lorca gives us nice Go cross-platform capabilities and is much less bloated.

# References

* [PWA Starter Kit](https://github.com/Polymer/pwa-starter-kit)
* [WiredJS Web Components](https://wiredjs.com)

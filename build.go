// +build mage

package main

import (
	"context"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build target is any exported function with zero args with no return or an error return.
// If a target has an error return and returns an non-nil error, mage will print
// that error to stdout and return with an exit code of 1.
func Install() error {
	return nil
}

// The first sentence in the comment will be the short help text shown with mage -l.
// The rest of the comment is long help text that will be shown with mage -h <target>
func Target() {
	// by default, the log stdlib package will be set to discard output.
	// Running with mage -v will set the output to stdout.
}

// A var named Default indicates which target is the default.  If there is no
// default, running mage will list the targets available.
var Default = Install

// BuildUI builds the web app in the www/ directory
func BuildUI(ctx context.Context) (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	println("Changing to www/")
	err = os.Chdir("www")
	if err != nil {
		return err
	}
	println("Running npm install")
	err = sh.RunCmd("npm", "install")()
	if err != nil {
		return err
	}
	println("Running npm run build")
	sh.RunCmd("npm", "run", "build")()
	if err != nil {
		return err
    }
	err = os.Chdir(cwd)
	if err != nil {
		return err
    }
    mg.CtxDeps(ctx, Target)
    
	return nil
}

// BuildGo runs `go generate` and `go build`
func BuildGo(ctx context.Context) error {
	println("Running go generate")
	err := sh.RunCmd("go", "generate")()
	if err != nil {
		return err
	}
	println("Running go build")
	err = sh.RunCmd("go", "build")()
	if err != nil {
		return err
	}
	mg.CtxDeps(ctx, Target)

	return nil
}

// BuildPkg runs packaging scripts for Linux, Mac, and Windows
func BuildPkg(ctx context.Context) error {
	return nil
}

// Build will:
// * Run BuildUI()
// * Run BuildGo()
// * Run BuildPkg()
func Build(ctx context.Context) error {
	println("Running BuildUI()")
	err := BuildUI(ctx)
	if err != nil {
		return err
	}
	println("Running BuildGo()")
	err = BuildGo(ctx)
	if err != nil {
		return err
	}
	println("Running BuildPkg()")
	err = BuildPkg(ctx)
	if err != nil {
		return err
	}
	mg.CtxDeps(ctx, Target)
	println("Success! Enjoy OuOfMyWay Time Tracker.")

	return nil
}

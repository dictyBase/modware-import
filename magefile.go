// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Build builds the binary for modware-import project
func Build() error {
	fmt.Println("building ......")
	return sh.Run("go", "build", "-o", "importer", "cmd/import/main.go")
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return nil
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
}

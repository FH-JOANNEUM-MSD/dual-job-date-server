package main

import "testing"

func TestMainPackageBuilds(t *testing.T) {
	mainFn := main
	if mainFn == nil {
		t.Fatal("expected main function to be defined")
	}
}

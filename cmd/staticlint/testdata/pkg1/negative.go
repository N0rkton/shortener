package main

import "os"

func foo() {
	os.Exit(2) // want
}

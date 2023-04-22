package main

import "os"

func main() {
	os.Exit(2) // want "Exit is forbidden"
}

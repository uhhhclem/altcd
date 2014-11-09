// Package main runs the cd web server; it expects to be run from the src directory.
package main

import (
	"altcd"
)

func main() {
	altcd.RunServer("./client/")
}

//Package to implement the tszero module
package main

import (
	"flag"
	"fmt"
)

const (
	tarFmt = "tar"
	zipFmt = "zip"
)

// This should be set to "tar" or "zip"
var format string

// If true, the help flag was specified.
var help bool

// The name of the archive file to process or empty string to indicate processing the stdin
var fileName string

// Set up the command line parsing
func init() {
	flag.StringVar(&format, "format", "", "The value of format must be tar or zip.")
	flag.BoolVar(&help, "help", false, "Specify this to see the help message.")
	flag.Parse()
}

// Print the help message.
func printHelp() {
	//TODO finish this
	fmt.Println("This is the help message")
}

// Handle a tar archive
func doTar() {

}

// Handle a zip archive
func doZip() {

}

//main entry point into
func main() {
	if help {
		printHelp()
	} else {
		switch format {
		case tarFmt:
			doTar()
		case zipFmt:
			doZip()
		default:
			printHelp()
		}
	}
	fmt.Println("tszero: ", help)
}

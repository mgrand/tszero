//Package to implement the tszero module
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	tarFmt = "tar"
	zipFmt = "zip"
)

// This should be set to "tar" or "zip"
var format string

// If true, the help flag was specified.
var help bool

// if true, the v (verbose) flag was specified
var verbose bool

// The name of the archive file to process or empty string to indicate processing the stdin
var fileName string

// Set up the command line parsing
func init() {
	log.SetOutput(os.Stderr)
	flag.StringVar(&format, "format", "", "The value of format must be tar or zip.")
	flag.BoolVar(&help, "help", false, "Specify this to see the help message.")
	flag.BoolVar(&verbose, "v", false, "Verbose")
	flag.Parse()
	fileName = flag.Arg(0)
}

// Print the help message.
func printHelp() {
	//TODO finish this
	fmt.Println("This is the help message")
}

// Handle a tar archive
func doTar(writer io.Reader) {

}

// Handle a zip archive
func doZip(writer io.Reader) {

}

// get the reader that we will use to read the archive.
func getReader() io.Reader {
	if len(fileName) == 0 {
		logMaybe("Reading from stdin")
		return os.Stdin
	}
	logMaybe("File name is " + fileName)
	file, err := os.Open("file.go") // For read access.
	if err == nil {
		return bufio.NewReader(file)
	}
	log.Fatal(err)
	//goland:noinspection GoUnreachableCode
	return nil // will never get here, but the compiler does not realize that os.Exit doesn't return
}

//Top level archive processing
func processArchive(processor func(writer io.Reader)) {
	reader := getReader()
	processor(reader)
}

//main entry point into
func main() {
	logMaybe("tszero starting")
	if help {
		printHelp()
	} else {
		switch format {
		case tarFmt:
			processArchive(doTar)
		case zipFmt:
			processArchive(doZip)
		default:
			// TODO auto-detect type of archive file and make the format flag optional.
			printHelp()
		}
	}
	logMaybe("tszero finished")
}

func logMaybe(msg string) {
	if verbose {
		log.Println(msg)
	}
}

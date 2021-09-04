//Package to implement the tszero module
package main

import (
	"archive/tar"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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

// Set the timestamp fields of a header to zero
func zeroHeaderTimeFields(header *tar.Header) {
	header.ChangeTime = time.Time{}
	header.AccessTime = time.Time{}
	header.ModTime = time.Time{}
}

// Handle a tar archive
func doTar(reader io.Reader) {
	tarReader := tar.NewReader(reader)
	tarWriter := tar.NewWriter(os.Stdout)
	//var error error = nil
	for {
		//var header *tar.Header = nil
		header, error := tarReader.Next()
		if error != nil {
			//TODO deal with specific errors.
			break
		}
		zeroHeaderTimeFields(header)
		tarWriter.WriteHeader(header)
		// TODO Finish this
	}
}

// Handle a zip archive
//goland:noinspection GoUnusedParameter
func doZip(reader io.Reader) {
	//TODO Finish this
}

// get the reader that we will use to read the archive.
func withReader(consumer func(io.Reader)) {
	if len(fileName) == 0 {
		logMaybe("Reading from stdin")
		consumer(os.Stdin)
		return
	}
	logMaybe("File name is " + fileName)
	file, err := os.Open("file.go") // For read access.
	if err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		consumer(bufio.NewReader(file))
		return
	}
	log.Fatal(err)
}

//main entry point into
func main() {
	logMaybe("tszero starting")
	if help {
		printHelp()
	} else {
		switch format {
		case tarFmt:
			withReader(doTar)
		case zipFmt:
			withReader(doZip)
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

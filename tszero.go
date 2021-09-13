//Package to implement the tszero module
package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"
)

const (
	tarFmt = "tar"
	zipFmt = "zip"
)

type config struct {
	// The program name
	programName string

	// This should be set to "tar" or "zip"
	format string

	// If true, the help flag was specified.
	help bool

	// if true, the v (verbose) flag was specified
	verbose bool

	// The name of the archive file to process or empty string to indicate processing the stdin
	fileName string

	// The command line arguments that are not flags
	args []string
}

// Set by main to the configuration
var conf config

const bufferSize = 1024 * 8

// parses the command-line arguments provided to the program and initialize flags set from the command line.
//
// os.Args[0] is provided as 'programName' and os.args[1:] as 'args'.
// Returns the Config in case parsing succeeded, or an error.
func initFlags(programName string, args []string) (cnf *config, err error) {
	flags := flag.NewFlagSet(programName, flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	var myConf config
	myConf.programName = programName
	flags.StringVar(&myConf.format, "format", "", "The value of format must be tar or zip.")
	flags.BoolVar(&myConf.help, "help", false, "Specify this to see the help message.")
	flags.BoolVar(&myConf.verbose, "v", false, "Verbose")

	err = flags.Parse(args)
	if err != nil {
		return nil, err
	}
	conf.args = flags.Args()
	return &myConf, nil
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
func doTar(reader io.Reader, out io.Writer) {
	tarReader := tar.NewReader(reader)
	tarWriter := tar.NewWriter(out)
	for {
		header, done := readTarHeader(tarReader)
		if done {
			return
		}
		logMaybe("processing ", header.Name)
		zeroHeaderTimeFields(header)
		writeTarHeader(tarWriter, header)
		copyContent(tarReader, tarWriter, header)
	}
}

func writeTarHeader(tarWriter *tar.Writer, header *tar.Header) {
	var h = *header
	h.Format = tar.FormatGNU
	writeErr := tarWriter.WriteHeader(&h)
	if writeErr != nil {
		log.Fatal("Error writing header", h, '\n', writeErr)
	}
}

func readTarHeader(tarReader *tar.Reader) (*tar.Header, bool) {
	header, readErr := tarReader.Next()
	if readErr != nil {
		if readErr == io.EOF {
			return nil, true
		}
		log.Fatal("Error reading next header: ", readErr)
	}
	return header, false
}

func copyContent(tarReader *tar.Reader, tarWriter *tar.Writer, header *tar.Header) {
	buffer := make([]byte, bufferSize)
	for {
		count, done := readTarContent(tarReader, buffer)
		if done {
			return
		}
		writeTarContent(tarWriter, buffer[:count], header)
	}
}

func writeTarContent(tarWriter *tar.Writer, buffer []byte, header *tar.Header) {
	switch header.Typeflag {
	case tar.TypeLink, tar.TypeSymlink, tar.TypeChar, tar.TypeBlock, tar.TypeDir, tar.TypeFifo:
		return // These types contain no data
	default:
		_, err := tarWriter.Write(buffer)
		if err != nil {
			log.Fatalf("Error writing contents of %+v: %s\nLength of write buffer is %d", header, err, len(buffer))
		}
	}
}

// Read the current tar content into the buffer.
// Returns true if all the bytes have been read or false if there are more bytes to be read.
func readTarContent(tarReader *tar.Reader, buffer []byte) (int, bool) {
	count, err := tarReader.Read(buffer)
	if err == io.EOF && count == 0 {
		return 0, true
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return count, false
}

// Handle a zip archive
//goland:noinspection GoUnusedParameter
func doZip(fileName string, out io.Writer) {
	zipReader, err := zip.OpenReader(fileName)
	if err != nil {
		log.Fatalf("Failed to open file %s: %s", flag.Arg(0), err)
	}
	defer func(zipReader *zip.ReadCloser) {
		err := zipReader.Close()
		if err != nil {
			log.Printf("Error closing input %s", fileName)
		}
	}(zipReader)
	zipWriter := zip.NewWriter(out)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Fatal("Error closing output")
		}
	}(zipWriter)
	for _, thisFile := range zipReader.File {
		logMaybe("Copying ", thisFile.Name)
		zeroZipHeaderTimestamps(thisFile)
		fileWriter := createHeader(zipWriter, thisFile.FileHeader)
		fileReader := getReader(thisFile)
		byteCount := copyFile(fileWriter, fileReader, thisFile)
		logMaybe("Copied ", fmt.Sprint(byteCount), " bytes.")
	}
}

func copyFile(fileWriter io.Writer, fileReader io.ReadCloser, thisFile *zip.File) int64 {
	byteCount, copyErr := io.Copy(fileWriter, fileReader)
	if copyErr != nil {
		log.Fatalf("Error (%s) copying file from source: %s", copyErr, thisFile.FileHeader.Name)
	}
	return byteCount
}

func getReader(thisFile *zip.File) io.ReadCloser {
	fileReader, err := thisFile.Open()
	if err != nil {
		log.Fatalf("Error (%s) opening file in zip for reading: %s", err, thisFile.FileHeader.Name)
	}
	return fileReader
}

func createHeader(zipWriter *zip.Writer, fh zip.FileHeader) io.Writer {
	fileWriter, err := zipWriter.CreateHeader(&fh)
	if err != nil {
		log.Fatalf("Error (%s) creating header in output: %+v", err, fh)
	}
	return fileWriter
}

//goland:noinspection GoDeprecation
func zeroZipHeaderTimestamps(thisFile *zip.File) {
	thisFile.FileHeader.Modified = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}

// provide the reader that we will use to read the archive.
func withFileReader(consumer func(io.Reader)) {
	file, err := os.Open(conf.args[0]) // For read access.
	if err == nil {
		consume(consumer, file)
		return
	}
	log.Fatal(err, " File name: ", conf.fileName)
}

// Call the given consumer function, passing it the given file object wrapped in a buffered reader.
// Ensure that the file is closed before returning.
func consume(consumer func(io.Reader), file *os.File) {
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	consumer(bufio.NewReader(file))
}

//main entry point into
func main() {
	log.SetOutput(os.Stderr)
	defer stacktrace()
	conf, err := initFlags(os.Args[0], os.Args[1:])
	if err != nil {
		log.Fatalf("Error parsing command line: %s", err)
	}
	logMaybe("tszero starting")
	if conf.help || len(conf.args) != 1 {
		printHelp()
	} else {
		switch conf.format {
		case tarFmt:
			withFileReader(func(reader io.Reader) { doTar(reader, os.Stdout) })
		case zipFmt:
			if len(conf.args) < 1 {
				log.Fatal("Processing a zip file requires at least one file name.")
			}
			doZip(conf.args[0], os.Stdout)
		default:
			// TODO auto-detect type of archive file and make the format flag optional.
			printHelp()
		}
	}
	logMaybe("tszero finished")
}

func logMaybe(msg ...string) {
	if conf.verbose {
		log.Println(msg)
	}
}

func stacktrace() {
	if r := recover(); r != nil {
		fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
	}
}

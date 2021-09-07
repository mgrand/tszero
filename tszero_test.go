package main

import (
	"archive/tar"
	"bytes"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

const tar1 = "testdata/test.tar"

func Test_consume(t *testing.T) {
	var called = false
	file, err := os.Open(tar1)
	if err != nil {
		t.Fatalf("Error opening %s", tar1)
	}
	consumer := func(reader io.Reader) { called = true }

	consume(consumer, file)
	if !called {
		t.Error("consume did not call the consumer function")
	}
	c := file.Close()
	if c == nil {
		t.Error("consum did not close its file")
	}
}

func Test_doTar(t *testing.T) {
	fileReader := openTestTarFile(t)
	buffer := bytes.NewBuffer(make([]byte, 200000))
	verbose = true
	doTar(fileReader, buffer)
	tarFileReader := tar.NewReader(openTestTarFile(t))
	tarBufferReader := tar.NewReader(buffer)
	for {
		fileHeader, fileHeaderErr := tarFileReader.Next()
		if fileHeaderErr != nil {
			if fileHeaderErr == io.EOF {
				break
			}
			log.Fatal("Error reading next file header: ", fileHeaderErr)
		}
		bufferHeader, bufferHeaderErr := tarBufferReader.Next()
		if bufferHeaderErr != nil {
			if bufferHeaderErr == io.EOF {
				log.Fatal("Output had EoF before input")
			}
			log.Fatal("Error reading next file header: ", bufferHeaderErr)
		}
		if !nonTimestampHeaderFieldsMatch(fileHeader, bufferHeader) {
			t.Fatalf("Headers do not match: %+v\nvs: %+v", fileHeader, bufferHeader)
		}
		if !timestampsAreZero(bufferHeader) {
			t.Fatalf("Timestamps are not zero: %+v", bufferHeader)
		}
		var readSize int = 2048
		var fileBuffer = make([]byte, readSize)
		var bufferBuffer = make([]byte, readSize)
		for {
			fileCount, err1 := tarFileReader.Read(fileBuffer)
			if fileCount == 0 && err1 == io.EOF {
				break
			}
			if err1 != nil {
				log.Fatal(err1)
			}

			bufferCount, err2 := tarBufferReader.Read(bufferBuffer)
			if bufferCount != fileCount {
				log.Fatalf("Length of content for %s is different for input and output", fileHeader.Name)
			}
			if err2 != nil {
				log.Fatal(err2)
			}
			if bytes.Compare(fileBuffer, bufferBuffer) != 0 {
				log.Fatalf("Content for %s is different", fileHeader.Name)
			}
		}
	}
}

func timestampsAreZero(header *tar.Header) bool {
	return header.AccessTime == time.Time{} &&
		header.ModTime == time.Time{} &&
		header.ChangeTime == time.Time{}
}

func nonTimestampHeaderFieldsMatch(h1 *tar.Header, h2 *tar.Header) bool {
	return h1.Name == h2.Name &&
		h1.Format == h2.Format &&
		h1.Size == h2.Size &&
		h1.Devmajor == h2.Devmajor &&
		h1.Devminor == h2.Devminor &&
		h1.Gid == h2.Gid &&
		h1.Gname == h2.Gname &&
		h1.Linkname == h2.Linkname &&
		h1.Mode == h2.Mode &&
		reflect.DeepEqual(h1.PAXRecords, h2.PAXRecords) &&
		h1.Typeflag == h2.Typeflag &&
		h1.Uid == h2.Uid &&
		h1.Uname == h2.Uname
}

func openTestTarFile(t *testing.T) *os.File {
	fileReader, err := os.Open(tar1)
	if err != nil {
		t.Fatalf("Failed to open %s", tar1)
	}
	return fileReader
}

func Test_doZip(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_logMaybe(t *testing.T) {
	type args struct {
		msg []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_printHelp(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_withFileReader(t *testing.T) {
	type args struct {
		consumer func(io.Reader)
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_withReader(t *testing.T) {
	type args struct {
		consumer func(io.Reader)
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

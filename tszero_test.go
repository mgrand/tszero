package main

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"testing"
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
	fileReader, err := os.Open(tar1)
	if err != nil {
		t.Fatalf("Failed to open %s", tar1)
	}
	buffer := bytes.NewBuffer(make([]byte, 200000))
	doTar(fileReader, buffer)

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

func Test_zeroHeaderTimeFields(t *testing.T) {
	type args struct {
		header *tar.Header
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

package main

import (
	"archive/tar"
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_consume(t *testing.T) {
	type args struct {
		consumer func(io.Reader)
		file     *os.File
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

func Test_copyContent(t *testing.T) {
	type args struct {
		tarReader *tar.Reader
		tarWriter *tar.Writer
		header    *tar.Header
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

func Test_doTar(t *testing.T) {
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

func Test_readTarContent(t *testing.T) {
	type args struct {
		tarReader *tar.Reader
		buffer    []byte
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readTarContent(tt.args.tarReader, tt.args.buffer)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readTarContent() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readTarContent() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_readTarHeader(t *testing.T) {
	type args struct {
		tarReader *tar.Reader
	}
	tests := []struct {
		name  string
		args  args
		want  *tar.Header
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readTarHeader(tt.args.tarReader)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readTarHeader() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readTarHeader() got1 = %v, want %v", got1, tt.want1)
			}
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

func Test_writeTarContent(t *testing.T) {
	type args struct {
		err       error
		tarWriter *tar.Writer
		buffer    []byte
		header    *tar.Header
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

func Test_writeTarHeader(t *testing.T) {
	type args struct {
		tarWriter *tar.Writer
		header    *tar.Header
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

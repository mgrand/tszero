package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"crypto/sha512"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

const tar1 = "testdata/test.tar"
const tar2 = "testdata/test2.tar"
const zip1 = "testdata/test.zip"
const zip2 = "testdata/test2.zip"

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
	fileReader, fileLength := openTestTarFile(t, tar1)
	tmpFile := createTempFile("tar")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Logf("Attempt to remove %s failed with error: %s", name, err)
		}
	}(tmpFile.Name())
	conf = &config{programName: "Test_doTar", format: "tar", help: false, verbose: true, fileName: fileReader.Name(), args: []string{fileReader.Name()}}
	doTar(fileReader, tmpFile)
	rewindTempFile(tmpFile)
	fileReader2, fileLength2 := openTestTarFile(t, tar1)
	if fileLength != fileLength2 {
		log.Fatalf("Tar file length changed; was %d and now %d", fileLength, fileLength2)
	}
	tarFileReader := tar.NewReader(fileReader2)
	tarBufferReader := tar.NewReader(tmpFile)
	for headerCount := 0; ; {
		fileHeader, fileHeaderErr := tarFileReader.Next()
		bufferHeader, bufferHeaderErr := tarBufferReader.Next()
		log.Printf("file err: %+v; header: %+v\n", fileHeaderErr, fileHeader)
		if fileHeaderErr != nil {
			if fileHeaderErr == io.EOF {
				log.Printf("Ending test for file next returning EOF after reading %d headers", headerCount)
				if bufferHeaderErr != io.EOF {
					log.Fatal("File is at EOF, but buffer is not.")
				}
				break
			}
			log.Fatal("Error reading next file header: ", fileHeaderErr)
		}
		if bufferHeaderErr != nil {
			if bufferHeaderErr == io.EOF {
				log.Fatalf("Output had EoF before input; Last input header was %+v", fileHeader)
			}
			log.Fatal("Error reading next file header: ", bufferHeaderErr)
		}
		checkHeaders(t, fileHeader, bufferHeader)
		headerCount += 1
		compareContent(uint64(fileHeader.Size), fileHeader.Name, tarFileReader.Read, tarBufferReader.Read)
	}
}

func compareContent(size uint64, name string, readFn1 func([]byte) (int, error), readFn2 func([]byte) (int, error)) {
	if size == 0 {
		log.Printf("No content for %s", name)
		return
	}
	var readSize int = 2048
	var fileBuffer = make([]byte, readSize)
	var tmpBuffer = make([]byte, readSize)
	log.Printf("Comparing content for %s", name)
	for {
		fileCount, err1 := readFn1(fileBuffer)
		if fileCount == 0 && err1 == io.EOF {
			log.Println("End of content")
			return
		}
		if err1 != nil && err1 != io.EOF {
			log.Fatalf("Error reading from original file: %s", err1)
		}

		tmpCount, err2 := readFn2(tmpBuffer)
		if tmpCount != fileCount {
			log.Fatalf("Length of content for %s is different for input and output", name)
		}
		if err2 != nil && err2 != io.EOF {
			log.Fatalf("Error reading from temp file: %s", err1)
		}
		if bytes.Compare(fileBuffer, tmpBuffer) != 0 {
			log.Fatalf("Content for %s is different", name)
		}
	}
}

func rewindTempFile(tmpFile *os.File) {
	newOffset, rewindErr := tmpFile.Seek(0, 0)
	if newOffset != 0 || rewindErr != nil {
		log.Fatalf("Rewind of temp file failed. After seek, offset was %d. Error: %s", newOffset, rewindErr)
	}
}

func createTempFile(prefix string) *os.File {
	tmpFile, tmpErr := ioutil.TempFile(".", prefix)
	if tmpErr != nil {
		log.Fatalf("Error creating temp file: %s", tmpErr)
	}
	return tmpFile
}

func checkHeaders(t *testing.T, fileHeader *tar.Header, bufferHeader *tar.Header) {
	if !nonTimestampHeaderFieldsMatch(fileHeader, bufferHeader) {
		t.Fatalf("Headers do not match: %+v\nvs: %+v", fileHeader, bufferHeader)
	}
	if !timestampsAreZero(bufferHeader) {
		t.Fatalf("Timestamps are not zero: %+v", bufferHeader)
	}
}

func timestampsAreZero(header *tar.Header) bool {
	return header.AccessTime.IsZero() &&
		header.ModTime == time.Unix(0, 0) &&
		header.ChangeTime.IsZero()
}

func nonTimestampHeaderFieldsMatch(h1 *tar.Header, h2 *tar.Header) bool {
	return h1.Name == h2.Name &&
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

func openTestTarFile(t *testing.T, fileName string) (*os.File, int64) {
	fileInfo, errStat := os.Stat(fileName)
	if errStat != nil {
		t.Fatalf("failed to stat %s; error: %s", tar1, errStat)
	}
	fileReader, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("Failed to open %s", fileName)
	}
	return fileReader, fileInfo.Size()
}

func Test_doZip(t *testing.T) {
	tmpFile := createTempFile("zip")
	tmpName := tmpFile.Name()
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Logf("Error removing temp output file %s", name)
		}
	}(tmpName)
	conf = &config{programName: "Test_doTar", format: "tar", help: false, verbose: true, fileName: zip1, args: []string{zip1}}
	doZip(zip1, tmpFile)
	closeFile(t, tmpFile)

	tmpReader := openZip(t, tmpName)
	testReader := openZip(t, zip1)
	for i, thisTestFile := range testReader.File {
		t.Logf("File %d: %s", i, thisTestFile.Name)
		thisTmpFile := tmpReader.File[i]
		t.Logf("tmp header %+v", thisTmpFile)
		modTime := thisTmpFile.Modified
		if modTime.UnixMicro() != 0 {
			t.Fatal("Timestamp is not zero")
		}
		if thisTestFile.UncompressedSize64 != thisTmpFile.UncompressedSize64 {
			t.Fatal("Lengths are not equal.")
		}
		thisTestReader := openFileInZip(t, thisTestFile)
		thisTmpReader := openFileInZip(t, thisTmpFile)
		compareContent(thisTestFile.UncompressedSize64, thisTestFile.Name, thisTestReader.Read, thisTmpReader.Read)
	}
}

func openFileInZip(t *testing.T, thisTestFile *zip.File) io.ReadCloser {
	reader, err := thisTestFile.Open()
	if err != nil {
		t.Fatalf("Error (%s) opening %s in %s", err, thisTestFile.Name, zip1)
	}
	return reader
}

func closeFile(t *testing.T, tmpFile *os.File) {
	err := tmpFile.Close()
	if err != nil {
		t.Logf("Error closing outout: %s", err)
	}
}

func openZip(t *testing.T, tmpName string) *zip.ReadCloser {
	reader, err := zip.OpenReader(tmpName)
	if err != nil {
		t.Fatalf("Failed to open %s", zip1)
	}
	return reader
}

func Test_initFlags(t *testing.T) {
	programName := "Test_initFlags"
	conf, err := initFlags(programName, []string{})
	if err != nil {
		t.Errorf("For empty command line initFlags returned error: %s", err)
	} else {
		if conf.help || conf.verbose || conf.programName != programName || conf.format != "" || len(conf.args) != 0 {
			t.Errorf("Empty command line parsed to unexpected configuration value %+v", conf)
		}
	}
	conf2, err2 := initFlags(programName, []string{"-v", "-help", "-format", "tar", tar1})
	if err2 != nil {
		t.Errorf("For full command line initFlags returned error: %s", err2)
	} else {
		if !conf2.help || !conf2.verbose || conf2.programName != programName || conf2.format != "tar" || len(conf2.args) != 1 || conf2.args[0] != tar1 {
			t.Errorf("Full command line parsed to unexpected configuration value %+v", conf)
		}
	}
	_, err3 := initFlags(programName, []string{"-bogus"})
	if err3 == nil {
		t.Error("Expected error by got none.")
	}
}

func Test_systemTar(t *testing.T) {
	fileReader1, _ := openTestTarFile(t, tar1)
	hasher1 := sha512.New()
	conf = &config{programName: "Test_doTar", format: "tar", help: false, verbose: true, fileName: fileReader1.Name(), args: []string{fileReader1.Name()}}
	doTar(fileReader1, hasher1)
	hash1 := hasher1.Sum(nil)

	fileReader2, _ := openTestTarFile(t, tar2)
	hasher2 := sha512.New()
	conf = &config{programName: "Test_doTar", format: "tar", help: false, verbose: true, fileName: fileReader2.Name(), args: []string{fileReader2.Name()}}
	doTar(fileReader2, hasher2)
	hash2 := hasher2.Sum(nil)

	if bytes.Compare(hash1, hash2) != 0 {
		t.Errorf("Hashes are unequal.")
	}
}

func Test_systemZip(t *testing.T) {
	hasher1 := sha512.New()
	conf = &config{programName: "Test_doTar", format: "tar", help: false, verbose: true, fileName: zip1, args: []string{zip1}}
	doZip(zip1, hasher1)
	hash1 := hasher1.Sum(nil)

	hasher2 := sha512.New()
	conf = &config{programName: "Test_doTar", format: "tar", help: false, verbose: true, fileName: zip2, args: []string{zip2}}
	doZip(zip2, hasher2)
	hash2 := hasher2.Sum(nil)

	if bytes.Compare(hash1, hash2) != 0 {
		t.Errorf("Hashes are unequal.")
	}
}

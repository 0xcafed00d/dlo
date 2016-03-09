package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/simulatedsimian/assert"
)

func TestFileIndex(t *testing.T) {
	verify := assert.Make(t)

	fi := MakeFileIndex("data")

	verify(fi.makeFileName(0)).Equal("data/00000/000")
	verify(fi.makeFileName(999)).Equal("data/00000/999")
	verify(fi.makeFileName(1000)).Equal("data/00001/000")
	verify(fi.makeFileName(12345678)).Equal("data/12345/678")
}

func exist(fname string) bool {
	_, err := os.Stat(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func TestStore(t *testing.T) {

	verify := assert.Make(t)

	path := filepath.Join("/tmp", "testdata", strconv.FormatInt(time.Now().UnixNano(), 10))

	fi := MakeFileIndex(path)

	err := fi.MakeDummyFiles(3999)
	verify(err).IsNil()

	fi.RefeshFileCount()
	verify(fi.fileCount).Equal(uint64(3999))

	verify(exist(filepath.Join(path, "00003", "999"))).Equal(false)
	index := fi.ReserveFileIndex()
	verify(fi.StoreFile(index, "test file")).NoError()
	verify(exist(filepath.Join(path, "00003", "999"))).Equal(true)

	verify(exist(filepath.Join(path, "00004", "000"))).Equal(false)
	index = fi.ReserveFileIndex()
	verify(fi.StoreFile(index, "test file")).NoError()
	verify(exist(filepath.Join(path, "00004", "000"))).Equal(true)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileIndex struct {
	sync.Mutex
	dataRoot  string
	fileCount uint64
}

func MakeFileIndex(dataRoot string) *FileIndex {
	fi := FileIndex{}

	fi.dataRoot = dataRoot
	fi.fileCount = 0

	return &fi
}

func (fi *FileIndex) RefeshFileCount() error {
	return nil
}

func (fi *FileIndex) MakeDummyFiles(count int) error {

	return nil
}

func (fi *FileIndex) makeFileName(index uint64) string {
	return filepath.Join(fi.dataRoot,
		fmt.Sprintf("%05v", index/1000),
		fmt.Sprintf("%03v", index%1000))
}

func (fi *FileIndex) StoreFile(index uint64, text string) error {
	f, err := os.Create(fi.makeFileName)
	if err != nil {

	}

	return nil
}

func (fi *FileIndex) ReserveFileIndex() uint64 {
	fi.Mutex.Lock()
	defer fi.Unlock()

	fi.fileCount++
	return fi.fileCount
}

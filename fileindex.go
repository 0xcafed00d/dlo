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

func (fi *FileIndex) RefeshFileCount() {
	fi.Mutex.Lock()
	defer fi.Unlock()

	fi.fileCount = 0

	folders, err := filepath.Glob(filepath.Join(fi.dataRoot, "*"))
	if err != nil {
		return
	}

	if len(folders) == 0 {
		return
	}
	files, err := filepath.Glob(filepath.Join(folders[len(folders)-1], "*"))

	fi.fileCount = uint64(((len(folders) - 1) * 1000) + len(files))
}

func (fi *FileIndex) MakeDummyFiles(count int) error {
	for n := 0; n < count; n++ {
		index := fi.ReserveFileIndex()
		err := fi.StoreFile(index, fmt.Sprintf("file index %v", index))
		if err != nil {
			return err
		}
	}
	return nil
}

func (fi *FileIndex) makeFileName(index uint64) string {
	return filepath.Join(fi.dataRoot,
		fmt.Sprintf("%05v", index/1000),
		fmt.Sprintf("%03v", index%1000))
}

func (fi *FileIndex) StoreFile(index uint64, text string) error {

	// first file in a set? create folder.
	if index%1000 == 0 {
		err := os.MkdirAll(filepath.Join(fi.dataRoot, fmt.Sprintf("%05v", index/1000)), os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(fi.makeFileName(index))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprint(f, text)
	if err != nil {
		return err
	}
	return nil
}

func (fi *FileIndex) ReserveFileIndex() uint64 {
	fi.Mutex.Lock()
	defer fi.Unlock()

	n := fi.fileCount
	fi.fileCount++
	return n
}

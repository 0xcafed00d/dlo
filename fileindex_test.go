package main

import (
	"testing"

	"github.com/simulatedsimian/assert"
)

func TestFileIndex(t *testing.T) {
	fi := MakeFileIndex("data")

	assert.Equal(t, fi.makeFileName(0), "data/00000/000")
	assert.Equal(t, fi.makeFileName(999), "data/00000/999")
	assert.Equal(t, fi.makeFileName(1000), "data/00001/000")
	assert.Equal(t, fi.makeFileName(12345678), "data/12345/678")
}

package freewritechecker

import (
	"io/ioutil"
	"testing"
)

func TestFreeWriteChecker(t *testing.T) {
	file, err := ioutil.ReadFile("test.txt")
	if err != nil {
		t.Error(err)
	}

	got, err := WordCount(file)
	want := 54
	if got != want {
		t.Errorf("want %d; got %d", want, got)
	}
}

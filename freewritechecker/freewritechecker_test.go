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

	t.Run("Returns Correct Count", func(t *testing.T) {
		got := WordCount(file)
		want := 48
		if got != want {
			t.Errorf("want %d; got %d", want, got)
		}
	})
}
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

	t.Run("Returns Correct Word Count", func(t *testing.T) {
		got := WordCount(file)
		want := 48
		if got != want {
			t.Errorf("want %d; got %d", want, got)
		}
	})

	t.Run("Returns string of text with notes removed", func(t *testing.T) {
		notesRemoved, err := ioutil.ReadFile("notesRemoved.txt")
		if err != nil {
			t.Error(err)
		}
		got := NoteRemover(file)
		want := string(notesRemoved)
		if got != want {
			t.Errorf("\nwant:\n%v\ngot:\n%v", want, got)
		}
	})

}

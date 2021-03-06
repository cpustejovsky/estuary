package freewritechecker_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/cpustejovsky/estuary/freewritechecker"
)

func TestFreeWriteChecker(t *testing.T) {
	file, err := ioutil.ReadFile("test.txt")
	if err != nil {
		t.Error(err)
	}

	notesRemoved, err := ioutil.ReadFile("notesRemoved.txt")
	if err != nil {
		t.Error(err)
	}
	text := string(file)

	t.Run("Returns Correct Word Count", func(t *testing.T) {
		got := freewritechecker.WordCount(text)
		want := 48
		if got != want {
			t.Errorf("want %d; got %d", want, got)
		}
	})

	t.Run("Returns string of text with notes removed", func(t *testing.T) {
		got, _ := freewritechecker.NoteChecker(text)
		want := string(notesRemoved)
		if got != want {
			t.Errorf("\nwant:\n%v\ngot:\n%v", want, got)
		}
	})

	t.Run("Returns an array of notes", func(t *testing.T) {
		_, got := freewritechecker.NoteChecker(text)
		want := []string{"this is a note", "this is a note"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\nwant:\n%v\ngot:\n%v", want, got)
		}
	})

}

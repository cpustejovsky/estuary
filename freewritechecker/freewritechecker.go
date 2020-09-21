package freewritechecker

import (
	"regexp"
)

func WordCount(file []byte) int {
	text := string(file)
	//TODO: determine if error checking is needed
	total := regexp.MustCompile(`\w+`)
	matches := total.FindAllStringIndex(text, -1)

	notes := regexp.MustCompile(`[#][n]|[#][N]`)
	noteMatches := notes.FindAllStringIndex(text, -1)

	wordCount := len(matches) - len(noteMatches)
	return wordCount
}

func NoteRemover(file []byte) string {
	return string(file)
}

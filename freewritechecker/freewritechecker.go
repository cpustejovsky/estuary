package freewritechecker

import (
	"regexp"
)

func WordCount(file []byte) (int, error) {
	text := string(file)

	total := regexp.MustCompile(`\w+`)
	matches := total.FindAllStringIndex(text, -1)

	notes := regexp.MustCompile(`[#][n]|[#][N]`)
	noteMatches := notes.FindAllStringIndex(text, -1)

	wordCount := len(matches) - len(noteMatches)
	return wordCount, nil
}

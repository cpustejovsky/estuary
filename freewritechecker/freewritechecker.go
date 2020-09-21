package freewritechecker

import (
	"regexp"
)

func WordCount(file []byte) int {
	text := string(file)
	//TODO: determine if error checking is needed
	total := regexp.MustCompile(`\w+`)
	matches := total.FindAllStringIndex(text, -1)

	noteSymbol := regexp.MustCompile(`[#][n]|[#][N]`)
	symbolMatches := noteSymbol.FindAllStringIndex(text, -1)

	wordCount := len(matches) - len(symbolMatches)
	return wordCount
}

func NoteRemover(file []byte) string {
	text := string(file)
	notes := regexp.MustCompile(`\s[#][n].*|\s[#][N].*`)
	newtext := notes.ReplaceAllString(text, "")
	return newtext
}

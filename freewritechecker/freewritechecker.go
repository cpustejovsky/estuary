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

func NoteChecker(file []byte) []string {
	text := string(file)
	re := regexp.MustCompile(`[#][n].*|[#][N].*`)
	symb := regexp.MustCompile(`[#][n].|[#][N].`)
	t := re.FindAllString(text, -1)
	var notes []string
	for _, s := range t {
		modified := symb.ReplaceAllString(s, "")
		notes = append(notes, modified)
	}
	return notes
}

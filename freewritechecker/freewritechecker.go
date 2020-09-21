package freewritechecker

import (
	"regexp"
)

func WordCount(file []byte) int {
	text := string(file)
	//TODO: determine if error checking is needed
	t := regexp.MustCompile(`\w+`)
	total := t.FindAllStringIndex(text, -1)

	n := regexp.MustCompile(`[#][n]|[#][N]`)
	symbols := n.FindAllStringIndex(text, -1)

	return len(total) - len(symbols)
}

func NoteRemover(file []byte) string {
	text := string(file)
	re := regexp.MustCompile(`\s[#][n].*|\s[#][N].*`)
	return re.ReplaceAllString(text, "")
}

func NoteChecker(file []byte) []string {
	text := string(file)
	var notes []string

	re := regexp.MustCompile(`[#][n].*|[#][N].*`)
	t := re.FindAllString(text, -1)

	symb := regexp.MustCompile(`[#][n].|[#][N].`)
	for _, s := range t {
		modified := symb.ReplaceAllString(s, "")
		notes = append(notes, modified)
	}

	return notes
}

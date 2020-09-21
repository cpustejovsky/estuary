package freewritechecker

import (
	"regexp"
)

func WordCount(file []byte) int {
	text := string(file)
	//TODO: determine if error checking is needed
	t := regexp.MustCompile(`\w+`)
	total := t.FindAllStringIndex(text, -1)

	n := regexp.MustCompile(`[#][nN]`)
	symbols := n.FindAllStringIndex(text, -1)

	return len(total) - len(symbols)
}

func NoteRemover(file []byte) string {
	text := string(file)
	re := regexp.MustCompile(`\s[#][nN].*`)
	return re.ReplaceAllString(text, "")
}

func NoteChecker(file []byte) []string {
	text := string(file)
	var notes []string

	re := regexp.MustCompile(`[#][nN].*`)
	t := re.FindAllString(text, -1)

	symb := regexp.MustCompile(`[#][nN].`)
	for _, s := range t {
		modified := symb.ReplaceAllString(s, "")
		notes = append(notes, modified)
	}

	return notes
}

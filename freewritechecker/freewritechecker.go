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

func NoteChecker(file []byte) (string, []string) {
	t := string(file)
	var notes []string
	re := regexp.MustCompile(`\s[#][nN].*`)
	text := re.ReplaceAllString(t, "")
	n := re.FindAllString(t, -1)

	symb := regexp.MustCompile(`\s[#][nN].`)
	for _, s := range n {
		modified := symb.ReplaceAllString(s, "")
		notes = append(notes, modified)
	}

	return text, notes
}

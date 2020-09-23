package freewritechecker

import (
	"regexp"
)

func WordCount(text string) int {
	//TODO: determine if error checking is needed
	t := regexp.MustCompile(`\w+`)
	total := t.FindAllStringIndex(text, -1)

	n := regexp.MustCompile(`[#][nN]`)
	symbols := n.FindAllStringIndex(text, -1)

	return len(total) - len(symbols)
}

func NoteChecker(text string) (string, []string) {
	var notes []string
	re := regexp.MustCompile(`\s[#][nN].*`)
	modified := re.ReplaceAllString(text, "")
	n := re.FindAllString(text, -1)

	symb := regexp.MustCompile(`\s[#][nN].`)
	for _, s := range n {
		m := symb.ReplaceAllString(s, "")
		notes = append(notes, m)
	}

	return modified, notes
}

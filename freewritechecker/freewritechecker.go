package freewritechecker

import (
	"regexp"
)

func WordCount(file []byte) (int, error) {
	re := regexp.MustCompile(`\s`)
	text := string(file)
	split := re.Split(text, -1)
	return len(split), nil
}

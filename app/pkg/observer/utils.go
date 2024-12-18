package observer

import "strings"

// CheckForIgnore checks for blacklist words and ignore tag
func CheckForIgnore(word string, tagValue string) bool {
	if tagValue == TagIgnoreVal {
		return true
	}
	for _, w := range blackListWords {
		if strings.Contains(word, w) {
			return true
		}
	}
	return false
}

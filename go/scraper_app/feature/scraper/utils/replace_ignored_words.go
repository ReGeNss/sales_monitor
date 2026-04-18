package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ReplaceIgnoredWords(title string, wordsToIgnore []string) string {
	title = strings.TrimSpace(title)
    if len(wordsToIgnore) == 0 {
        return title
    }

    pattern := fmt.Sprintf("(?i)(%s)", strings.Join(wordsToIgnore, "|")) 
    re := regexp.MustCompile(pattern)
    
    result := re.ReplaceAllString(title, "")
    
    return result
}
package paramformatter

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TitleFormatter struct{}

const alwaysLowerWord = "a an and in of on the to"

func (tf *TitleFormatter) IsSingleWord(param string) bool {
	single_word_re := regexp.MustCompile(`^[[:alpha:]]+$`)

	return single_word_re.MatchString(param)
}

func (tf *TitleFormatter) IsMultiWord(param string) bool {
	multi_word_pattern := `^(\w+\_)+\w+$`
	multi_word_re := regexp.MustCompile(multi_word_pattern)

	return multi_word_re.MatchString(param)
}

func (tf *TitleFormatter) Run(param string, firstWordOnly bool) string {
	c := cases.Title(language.Und)
	firstWordTitle := c.String(strings.ToLower(param))

	if firstWordOnly {
		return firstWordTitle
	}

	// Return title case phrase
	words := strings.Split(param, "_")
	for i, w := range words {
		lc := strings.ToLower(w)
		// But keep some words always lower case except when they are the first word in a title
		if w != words[0] && tf.isAlwaysLowerWord(lc) {
			words[i] = lc
		} else {
			words[i] = c.String(strings.ToLower(w))
		}
	}
	fmt.Println(words)
	return strings.Join(words, "_")
}

func (tf *TitleFormatter) isAlwaysLowerWord(word string) bool {
	return strings.Contains(alwaysLowerWord, " "+word+" ")
}

func NewTitleFormatter() *TitleFormatter {
	return &TitleFormatter{}
}

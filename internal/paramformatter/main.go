package paramformatter

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TitleFormatter struct{}

func (pf *TitleFormatter) IsSingleWord(param string) bool {
	single_word_re := regexp.MustCompile(`^[[:alpha:]]+$`)

	return single_word_re.MatchString(param)
}

func (pf *TitleFormatter) IsMultiWord(param string) bool {
	multi_word_pattern := `^(\w+\_)+\w+$`
	multi_word_re := regexp.MustCompile(multi_word_pattern)

	return multi_word_re.MatchString(param)
}

func (pf *TitleFormatter) Run(param string, firstWordOnly bool) string {
	c := cases.Title(language.Und)
	firstWordTitle := c.String(strings.ToLower(param))

	if firstWordOnly {
		return firstWordTitle
	}

	words := strings.Split(param, "_")
	for i, w := range words {
		words[i] = c.String(strings.ToLower(w))
	}

	return strings.Join(words, "_")
}

func NewTitleFormatter() *TitleFormatter {
	return &TitleFormatter{}
}

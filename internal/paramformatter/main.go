package paramformatter

import (
	"regexp"
	"strings"
)

type TitleFormatter struct{}

func (pf *TitleFormatter) Run(param string) string {
	pattern := `^(\w+\_)+\w+$`

	re := regexp.MustCompile(pattern)

	if re.MatchString(param) {
		return pf.titleize(param)
	} else {
		return param
	}

}

func (pf *TitleFormatter) titleize(param string) string {
	chars := strings.Split(param, "_")
	for i, c := range chars {
		chars[i] = strings.Title(strings.ToLower(c))
	}

	return strings.Join(chars, "_")
}

func NewTitleFormatter() *TitleFormatter {
	return &TitleFormatter{}
}

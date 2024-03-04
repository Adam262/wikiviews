package paramvalidator

import (
	"fmt"
	"regexp"
	"strings"
)

const forbiddenChars = "#<>[]{}|"

type TitleValidator struct{}

func (tv *TitleValidator) Run(param string) (isValid bool, err error) {
	// Check for an empty param
	if len(param) == 0 {
		err = fmt.Errorf("error: article param is invalid: param cannot be empty")
		return
	}

	// Check for a lower-case first character
	alphaRe := regexp.MustCompile(`^[[:alpha:]]$`)
	firstChar := string(param[0])

	if alphaRe.MatchString(firstChar) && strings.ToLower(firstChar) == firstChar {
		err = fmt.Errorf("error: article param %s is invalid: param must not begin with a lower case character", param)
		return
	}

	// Check for forbidden characters
	forbiddenRe := regexp.MustCompile(`.*(\#|\<|\>|\[|\]|\{|\}|\|)+.*`)
	if forbiddenRe.MatchString(param) {
		err = fmt.Errorf("error: article param %s is invalid: param must not contain chars %s", param, forbiddenChars)
		return
	}

	return true, nil
}

func NewTitleValidator() *TitleValidator {
	return &TitleValidator{}
}

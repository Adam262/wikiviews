package paramvalidator

import (
	"fmt"
	"regexp"
	"strings"
)

const forbiddenChars = "#<>[]{}|"

type TitleValidator struct{}

func (pf *TitleValidator) Run(param string) (isValid bool, err error) {
	// Check for an empty param
	if len(param) == 0 {
		err = fmt.Errorf("error: article param %s is invalid: param cannot be empty", param)
		return false, err
	}

	// Check for a lower-case first character
	alphaRe := regexp.MustCompile(`^[[:alpha:]]$`)
	firstChar := string(param[0])

	if alphaRe.MatchString(firstChar) && strings.ToLower(firstChar) == firstChar {
		err = fmt.Errorf("error: article param %s is invalid: param must not begin with a lower case character", param)
		return false, err
	}

	// Check for forbidden characters
	forbiddenRe := regexp.MustCompile(`.*(\#|\<|\>|\[|\]|\{|\}|\|)+.*`)
	if forbiddenRe.MatchString(param) {
		err = fmt.Errorf("error: article param %s is invalid: param must not contain chars %s", param, forbiddenChars)
		return false, err
	}

	return true, nil
}

func NewTitleValidator() *TitleValidator {
	return &TitleValidator{}
}

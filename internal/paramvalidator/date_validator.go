package paramvalidator

import (
	"fmt"
	"regexp"
)

type DateValidator struct{}

const yearMonth = `^[12]\d{3}(0[1-9]|1[0-2])$`

func (dv *DateValidator) Run(date string) (isValid bool, err error) {
	if len(date) == 0 {
		err = fmt.Errorf("error: date param is invalid: param cannot be empty. Please enter in form YYYYMM")
		return
	}

	re := regexp.MustCompile(yearMonth)
	if !re.MatchString(date) {
		err = fmt.Errorf("error: date param is invalid: please enter a valid year and month in form YYYYMM")
		return
	}

	return true, nil
}

func NewDateValidator() *DateValidator {
	return &DateValidator{}
}

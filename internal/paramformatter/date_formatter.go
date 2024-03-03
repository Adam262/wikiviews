package paramformatter

import (
	"fmt"
	"regexp"
	"strconv"
)

type (
	DateFormatter struct{}
)

var monthLengths = map[string]string{
	"01":   "31",
	"02":   "28",
	"02LY": "29",
	"03":   "31",
	"04":   "30",
	"05":   "31",
	"06":   "30",
	"07":   "31",
	"08":   "31",
	"09":   "30",
	"10":   "31",
	"11":   "30",
	"12":   "31",
}

func (df *DateFormatter) Run(date string) (start, end string, err error) {
	year := date[:4]
	month := date[4:]

	// Check for edge case of February in leap year
	if month == "02" {
		isLeapYear, lyErr := df.isLeapYear(year)
		if lyErr != nil {
			err = fmt.Errorf("error: %s", lyErr)
			return
		}

		if isLeapYear {
			month = "02LY"
		}
	}

	start = fmt.Sprintf("%s%s", date, "01")
	end = fmt.Sprintf("%s%s", date, monthLengths[month])

	return
}

func (df *DateFormatter) isLeapYear(year string) (isLeapYear bool, err error) {
	yearRe := regexp.MustCompile(`^\d{4}$`)
	if !yearRe.MatchString(year) {
		err = fmt.Errorf("error: year param is invalid: must be in form YYYY")
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		err = fmt.Errorf("error converting year %s to int", year)
		return
	}

	return yearInt%4 == 0, nil
}

func NewDateFormatter() *DateFormatter {
	return &DateFormatter{}
}

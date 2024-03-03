package paramvalidator

import (
	"testing"
)

func TestDateValidator_Run(t *testing.T) {
	validator := NewDateValidator()

	testCases := []struct {
		param   string
		isValid bool
	}{
		{"202401", true},
		{"202402", true},
		{"302402", false},
		{"002402", false},
		{"202410", true},
		{"202411", true},
		{"202412", true},
		{"202413", false},
		{"", false},
		{"20240", false},
		{"2024", false},
		{"202", false},
		{"20", false},
		{"2", false},
	}

	for _, tc := range testCases {
		isValid, _ := validator.Run(tc.param)

		if isValid != tc.isValid {
			t.Errorf("TestDateValidator.Run(%q) returns isValid = %t; Expected %t", tc.param, isValid, tc.isValid)
		}
	}
}

package paramvalidator

import (
	"testing"
)

func TestTitleValiadator_Run(t *testing.T) {
	validator := NewTitleValidator()

	testCases := []struct {
		param   string
		isValid bool
	}{
		{"Abc", true},
		{"A#bc", false},
		{"A<bc", false},
		{"A>bc", false},
		{"A[bc", false},
		{"A]bc", false},
		{"A{bc", false},
		{"A}bc", false},
		{"A|bc", false},
		{"A#", false},
		{"#a", false},
		{"123", true},
		{"", false},
		{"abc", false},
	}

	for _, tc := range testCases {
		isValid, _ := validator.Run(tc.param)

		if isValid != tc.isValid {
			t.Errorf("TestTitleValidator.Run(%q) returns isValid = %t; Expected %t", tc.param, isValid, tc.isValid)
		}
	}
}

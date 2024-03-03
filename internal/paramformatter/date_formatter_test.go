package paramformatter

import "testing"

func TestDateFormatter_Run(t *testing.T) {
	formatter := NewDateFormatter()

	testCases := []struct {
		date          string
		expectedStart string
		expectedEnd   string
	}{
		{"202301", "20230101", "20230131"},
		{"202401", "20240101", "20240131"},
		{"202302", "20230201", "20230228"},
		{"202402", "20240201", "20240229"},
	}

	for _, tc := range testCases {
		actualStart, actualEnd, _ := formatter.Run(tc.date)

		if actualStart != tc.expectedStart {
			t.Errorf("TestDateFormatter.Run(%q) returns start date %q; Expected %q", tc.date, actualStart, tc.expectedStart)
		}

		if actualEnd != tc.expectedEnd {
			t.Errorf("TestDateFormatter.Run(%q) returns end date %q; Expected %q", tc.date, actualEnd, tc.expectedEnd)
		}
	}
}

func TestDateFormatter_isLeapYear(t *testing.T) {
	formatter := NewDateFormatter()

	testCases := []struct {
		year           string
		expectedResult bool
	}{
		{"2020", true},
		{"2021", false},
		{"2022", false},
		{"2023", false},
		{"2024", true},
	}

	for _, tc := range testCases {
		actualResult, _ := formatter.isLeapYear(tc.year)

		if actualResult != tc.expectedResult {
			t.Errorf("TestFormatter.IsSingleWord(%q) returns %t; Expected %t", tc.year, actualResult, tc.expectedResult)
		}
	}
}

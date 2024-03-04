package paramformatter

import "testing"

func TestTitleFormatter_Run(t *testing.T) {
	formatter := NewTitleFormatter()

	testCases := []struct {
		param         string
		firstWordOnly bool
		newParam      string
	}{
		{"michael_phelps", true, "Michael_phelps"},
		{"michael_phelps", false, "Michael_Phelps"},
		{"!!!", true, "!!!"},
		{"!!!", false, "!!!"},
		{"MICHAEL_PHELPS", true, "Michael_phelps"},
		{"MICHAEL_PHELPS", false, "Michael_Phelps"},
		{"mIcHaeL_pHeLPS", true, "Michael_phelps"},
		{"mIcHaeL_pHeLPS", false, "Michael_Phelps"},
		{"michael_Phelps", true, "Michael_phelps"},
		{"michael_Phelps", false, "Michael_Phelps"},
		{"MichaelPhelps", true, "Michaelphelps"},
		{"MichaelPhelps", false, "Michaelphelps"},
		{"man_page", true, "Man_page"},
		{"man_page", false, "Man_Page"},
		{"Man_page", true, "Man_page"},
		{"Man_page", false, "Man_Page"},
		{"man_page", true, "Man_page"},
		{"man_page", false, "Man_Page"},
		{"MAN_PAGE", true, "Man_page"},
		{"MAN_PAGE", false, "Man_Page"},
		{"Orca", true, "Orca"},
		{"orca", false, "Orca"},
		{"call_of_the_wild", false, "Call_of_the_Wild"},
		{"call_of_the_wild", true, "Call_of_the_wild"},
		{"Call_Of_the_wild", false, "Call_of_the_Wild"},
		{"on_golden_pond", false, "On_Golden_Pond"},
		{"On_golden_pond", true, "On_golden_pond"},
	}

	for _, tc := range testCases {
		actualNewParam := formatter.Run(tc.param, tc.firstWordOnly)

		if actualNewParam != tc.newParam {
			t.Errorf("TestTitleFormatter.Run(%q) returns new param %q; Expected %q", tc.param, actualNewParam, tc.newParam)
		}
	}
}

func TestTitleFormatter_IsSingleWord(t *testing.T) {
	formatter := NewTitleFormatter()

	testCases := []struct {
		param          string
		expectedResult bool
	}{
		{"michael_phelps", false},
		{"Michael_Phelps", false},
		{"Orca", true},
		{"orca", true},
		{"", false},
		{"!!!", false},
	}

	for _, tc := range testCases {
		actualResult := formatter.IsSingleWord(tc.param)

		if actualResult != tc.expectedResult {
			t.Errorf("TestFormatter.IsSingleWord(%q) returns %t; Expected %t", tc.param, actualResult, tc.expectedResult)
		}
	}
}

func TestTitleFormatter_IsMultiWord(t *testing.T) {
	formatter := NewTitleFormatter()

	testCases := []struct {
		param          string
		expectedResult bool
	}{
		{"michael_phelps", true},
		{"Michael_Phelps", true},
		{"New_York_city", true},
		{"Orca", false},
		{"orca", false},
		{"", false},
		{"!!!", false},
	}

	for _, tc := range testCases {
		actualResult := formatter.IsMultiWord(tc.param)

		if actualResult != tc.expectedResult {
			t.Errorf("TestFormatter.IsMultiWord(%q) returns %t; Expected %t", tc.param, actualResult, tc.expectedResult)
		}
	}
}

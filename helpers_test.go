package saga

import (
	"reflect"
	"testing"
)

func compareTable(t *testing.T, actual, expected *Table) bool {
	// validate tables
	if err := actual.validate(); err != nil {
		t.Errorf("Actual table is invalid.\n")
		t.Errorf("    %v\n", err)
		return false
	}

	if err := expected.validate(); err != nil {
		t.Errorf("Expected table is invalid.\n")
		t.Errorf("    %v\n", err)
		return false
	}

	// should have same number of headers
	{
		if len(actual.headers) != len(expected.headers) {
			t.Errorf("Number of headers do not match.\n")
			t.Errorf("    actual:   %v\n", actual.headers)
			t.Errorf("    expected: %v\n", expected.headers)
			return false
		}

		// find headers that are in actual but not in expected
		actualExtraHeaders := []string{}
		for colName := range actual.headers {
			if _, ok := expected.headers[colName]; !ok {
				actualExtraHeaders = append(actualExtraHeaders, colName)
			}
		}

		// find headers that are in expected but not in actual
		expectedExtraHeaders := []string{}
		for colName := range expected.headers {
			if _, ok := actual.headers[colName]; !ok {
				expectedExtraHeaders = append(expectedExtraHeaders, colName)
			}
		}

		if len(actualExtraHeaders) != 0 || len(expectedExtraHeaders) != 0 {
			t.Errorf("Found extra header columns.\n")
			t.Errorf("    actual:   %v\n", actualExtraHeaders)
			t.Errorf("    expected: %v\n", expectedExtraHeaders)
			return false
		}
	}

	// should have same number of columns
	if len(actual.cols) != len(expected.cols) {
		t.Errorf("Number of columns do not match.\n")
		t.Errorf("    actual:   %v\n", len(actual.cols))
		t.Errorf("    expected: %v\n", len(expected.cols))
		return false
	}

	// at this point we know the Table's are valid and passed basic sanity
	// checks. Start comparing values
	{
		mismatchedColumns := []string{}
		for colName := range actual.headers {
			actualIdx := actual.headers[colName]
			actualValues := actual.cols[actualIdx]

			expectedIdx := expected.headers[colName]
			expectedValues := expected.cols[expectedIdx]

			if !reflect.DeepEqual(actualValues, expectedValues) {
				mismatchedColumns = append(mismatchedColumns, colName)
			}
		}

		if len(mismatchedColumns) > 0 {
			t.Errorf("Some column values mismatched.\n")
			t.Errorf("    columns that mismatched: %v\n", mismatchedColumns)
			return false
		}
	}

	return true
}

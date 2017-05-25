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

func TestNew(t *testing.T) {
	cases := map[string]struct {
		headers   []string
		rowValues [][]interface{}
		expected  *Table
	}{
		"headers_only": {
			headers:   []string{"id", "name", "active"},
			rowValues: nil,
			expected: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{},
					[]interface{}{},
					[]interface{}{},
				},
			},
		},

		"headers_and_values": {
			headers: []string{"id", "name", "active"},
			rowValues: [][]interface{}{
				{0, "Name 0", true},
				{1, "Name 1", false},
			},
			expected: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
					[]interface{}{true, false},
				},
			},
		},

		"no_headers": {
			headers:   nil,
			rowValues: nil,
			expected: &Table{
				headers: map[string]int{},
				cols:    [][]interface{}{},
			},
		},

		"no_headers_with_values": {
			headers: nil,
			rowValues: [][]interface{}{
				{0, "Name 0", true},
				{1, "Name 1", false},
			},
			expected: &Table{
				headers: map[string]int{},
				cols:    [][]interface{}{},
			},
		},

		"headers_with_fewer_number_of_cols": {
			headers: []string{"id", "name"},
			rowValues: [][]interface{}{
				{0, "Name 0", true},
				{1, "Name 1", false},
			},
			expected: &Table{
				headers: map[string]int{"id": 0, "name": 1},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := New(tc.headers, tc.rowValues...)
			compareTable(t, actual, tc.expected)
		})
	}
}

func TestInsertRows(t *testing.T) {
	cases := map[string]struct {
		headers   []string
		rowValues [][]interface{}
		initial   *Table
		expected  *Table
	}{
		"more_rows": {
			initial: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
					[]interface{}{true, false},
				},
			},
			headers: []string{"id", "name", "active"},
			rowValues: [][]interface{}{
				{2, "Name 2", false},
				{3, "Name 3", true},
			},
			expected: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1, 2, 3},
					[]interface{}{"Name 0", "Name 1", "Name 2", "Name 3"},
					[]interface{}{true, false, false, true},
				},
			},
		},

		"extra_headers": {
			initial: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
					[]interface{}{true, false},
				},
			},
			headers: []string{"id", "name", "active", "ieee754"},
			rowValues: [][]interface{}{
				{2, "Name 2", false, 1.23},
				{3, "Name 3", true, 2.34},
			},
			expected: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1, 2, 3},
					[]interface{}{"Name 0", "Name 1", "Name 2", "Name 3"},
					[]interface{}{true, false, false, true},
				},
			},
		},

		// TODO - initial is empty
		// TODO - subset of headers
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := tc.initial.InsertRows(tc.headers, tc.rowValues...)
			compareTable(t, actual, tc.expected)
		})
	}
}

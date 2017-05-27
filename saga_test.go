package saga

import (
	"testing"
)

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
				{2, "Name 2", false, 1.23, nil},
				{3, "Name 3", true, 2.34, nil},
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

		"extra_column_values": {
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

		"all_headers_missing_from_table": {
			initial: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
					[]interface{}{true, false},
				},
			},
			headers: []string{"colA", "colB", "colC"},
			rowValues: [][]interface{}{
				{2, "Name 2", false},
				{3, "Name 3", true},
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

		"subset_of_headers": {
			initial: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
					[]interface{}{true, false},
				},
			},
			headers: []string{"id"},
			rowValues: [][]interface{}{
				{2},
				{3},
			},
			expected: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1, 2, 3},
					[]interface{}{"Name 0", "Name 1", nil, nil},
					[]interface{}{true, false, nil, nil},
				},
			},
		},

		"nil_rowValues": {
			initial: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
					[]interface{}{true, false},
				},
			},
			headers:   []string{"id"},
			rowValues: nil,
			expected: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{0, 1},
					[]interface{}{"Name 0", "Name 1"},
					[]interface{}{true, false},
				},
			},
		},

		"empty_initial": {
			initial: &Table{
				headers: map[string]int{"id": 0, "name": 1, "active": 2},
				cols: [][]interface{}{
					[]interface{}{},
					[]interface{}{},
					[]interface{}{},
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
					[]interface{}{2, 3},
					[]interface{}{"Name 2", "Name 3"},
					[]interface{}{false, true},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := tc.initial.InsertRows(tc.headers, tc.rowValues...)
			compareTable(t, actual, tc.expected)
		})
	}
}

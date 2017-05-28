package saga

import (
	"log"
)

type Table struct {
	headers map[string]int
	cols    [][]interface{}

	// TODO - Document invariants for Table that must always hold
	// See Table.validate
}

// New creates a new table by using values from rowValues to fill the Table.
// Upto a maxiumum of N columnNames can be provided, where N is the length of a
// row. If fewer columnNames are provided only the first C column values from
// rowValues are used.
func New(headers []string, rowValues ...[]interface{}) *Table {
	// NOTE: assume that all rows in rowValues have length >= len(headers)

	t := &Table{
		headers: make(map[string]int, len(headers)),
		cols:    make([][]interface{}, len(headers)),
	}

	// pre-allocate cols
	for i := range t.cols {
		t.cols[i] = make([]interface{}, len(rowValues))
	}

	// set headers
	for i, colName := range headers {
		t.headers[colName] = i
	}

	// set column values; pivot rowValues of dimensions MxN to NxM
	for rowIdx, vals := range rowValues {
		for colIdx, val := range vals[:len(headers)] {
			t.cols[colIdx][rowIdx] = val
		}
	}

	return t
}

// InsertRows adds rowValues into Table. Values for any headers that does not
// exist in Table are skipped; headers which are not provided will have their
// values set to nil. Likewise, extra values at the end of the slice are also
// ignored
func (t *Table) InsertRows(headers []string, rowValues ...[]interface{}) *Table {
	if len(rowValues) == 0 {
		return t
	}

	// ensure that at least one element in headers exists in t.headers
	found := false
	for _, h := range headers {
		if _, ok := t.headers[h]; ok {
			found = true
			break
		}
	}
	if !found {
		return t
	}

	offset := len(t.cols[0])

	// pre-allocate memory
	for i := range t.cols {
		t.cols[i] = append(t.cols[i], make([]interface{}, len(rowValues))...)
	}

	for rowIdx, vals := range rowValues {
		for i, header := range headers {
			if colIdx, ok := t.headers[header]; ok {
				t.cols[colIdx][offset+rowIdx] = vals[i]
			}
		}
	}

	return t
}

// AddColumn adds a new column with an optional defaultValue (nil if not
// provided); error if column already exists. defaultValue can be a plain value
// or gnerator `func() interface{}`
func (t *Table) AddColumn(colName string, defaultValue ...interface{}) *Table {
	// check if colName already exists
	if _, ok := t.headers[colName]; ok {
		log.Panicf("Column %v already exists in table", colName)
		return t
	}

	// optional defaultValue can be a generator function or plain value; normalize
	// everything to a gnerator function to make things simpler when setting values
	defaultFn := func() interface{} { return nil }
	if len(defaultValue) > 0 {
		switch val := defaultValue[0].(type) {
		case nil:
			// "empty" out defaultVaule; it will be checked later when we are
			// adding values
			defaultValue = defaultValue[0:0]
		case func() interface{}:
			defaultFn = val
		default:
			defaultFn = func() interface{} { return val }
		}
	}

	// update headers with new colName
	headersLen := len(t.headers)
	t.headers[colName] = headersLen

	// add len(cols[0]) number of values
	values := []interface{}{}
	if len(t.cols) > 0 {
		values = make([]interface{}, len(t.cols[0]))

		// values will all be nil by default, so only populate if the caller
		// provided a default value. This should avoid some unnecessary function calls
		if len(defaultValue) > 0 {
			for i := range values {
				values[i] = defaultFn()
			}
		}
	}

	t.cols = append(t.cols, values)

	return t
}

// UpdateColumn updates all values of colName to value. If the column does not
// exist it is added
func (t *Table) UpdateColumn(colName string, value interface{}) *Table {
	// idx indicates which column values to update
	idx, ok := t.headers[colName]

	// create a new
	if !ok {
		idx = len(t.headers)
		t.headers[colName] = idx

		var values []interface{}
		if len(t.cols) > 0 {
			values = make([]interface{}, len(t.cols[0]))
		}

		t.cols = append(t.cols, values)
	}

	defaultFn := func() interface{} { return value }
	switch val := value.(type) {
	case func() interface{}:
		defaultFn = val
	}

	for j := range t.cols[idx] {
		t.cols[idx][j] = defaultFn()
	}

	return t
}

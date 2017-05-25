package saga

import "log"

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

func NewFromMap(data ...map[string]interface{}) *Table {
	// TODO - We expecte data to be valid, is this okay?

	// sanity check, always return a valid Table
	if len(data) == 0 {
		return &Table{
			headers: map[string]int{},
		}
	}

	// colOrder is used to maintain order in which Table.cols will added;
	// index is table.cols offset, and value is the lookup key from data
	colOrder := make([]string, 0, len(data[0]))
	headers := make(map[string]int, len(data[0]))

	for colName := range data[0] {
		colOrder = append(colOrder, colName)

		// Since len(colOrder) increases by one every loop, we can use that
		// as the offset marker for colName instead of maintaining a separate
		// counter
		headers[colName] = len(colOrder) - 1
	}

	table := &Table{
		headers: headers,
		cols:    make([][]interface{}, len(data[0])),
	}

	// pre-allocate cols
	for i := range table.cols {
		table.cols[i] = make([]interface{}, len(data))
	}

	// add data to table
	for rowIdx, row := range data {
		for colIdx, colName := range colOrder {
			val := row[colName]
			table.cols[colIdx][rowIdx] = val
		}
	}

	return table
}

func (t *Table) InsertRows(headers []string, rowValues ...[]interface{}) *Table {
	// DONE - What if we get extra headers? - Ignore them
	// TODO - What if we get a subset of t.headers (t.cols is jagged now) - Set values to nil
	// TODO - Pre-allocate?

	for _, vals := range rowValues {
		for i, header := range headers {
			if colIdx, ok := t.headers[header]; ok {
				t.cols[colIdx] = append(t.cols[colIdx], vals[i])
			}
		}
	}

	return t
}

func (t *Table) InsertMap(row map[string]interface{}) *Table {
	// TODO - does not add new columns from row; should we be a little smarter
	// about this?

	if len(t.cols) != len(row) {
		log.Panicf("Column size does not match: %v != %v\n", len(t.cols), len(row))
		return t
	}

	for colName, i := range t.headers {
		val := row[colName]
		t.cols[i] = append(t.cols[i], val)
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

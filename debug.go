package saga

import "fmt"

func (t *Table) validate() error {
	if t == nil {
		return nil
	}

	// we should have the same number of headers as columns
	if len(t.headers) != len(t.cols) {
		return fmt.Errorf("number of headers does not match number of columns (%v != %v)",
			len(t.headers), len(t.cols))
	}

	// no index should be repeated in t.headers
	{
		idxMap := map[int][]string{}
		for colName, idx := range t.headers {
			colNames := idxMap[idx]
			idxMap[idx] = append(colNames, colName)
		}

		for idx, colNames := range idxMap {
			if len(colNames) == 1 {
				delete(idxMap, idx)
			}
		}

		if len(idxMap) > 0 {
			return fmt.Errorf("found duplicate indexes for some columns. %v", idxMap)
		}
	}

	// check that all cols have the same number of values
	if len(t.cols) > 0 {
		totalCols := len(t.cols[0])

		colsMismatch := map[string]int{}
		for colName, idx := range t.headers {
			if numCols := len(t.cols[idx]); numCols != totalCols {
				colsMismatch[colName] = numCols
			}
		}

		if len(colsMismatch) > 0 {
			return fmt.Errorf("some columns have unequal number of values (expected %v values). %v",
				totalCols, colsMismatch)
		}
	}

	return nil
}

func (t *Table) print() {
	if err := t.validate(); err != nil {
		panic(err)
	}

	fmt.Println("headers: ", t.headers)

	for name, idx := range t.headers {
		fmt.Printf("  %v (%v): ", name, len(t.cols[idx]))
		for _, val := range t.cols[idx] {
			fmt.Printf("  %v, ", val)
		}
		fmt.Println()
	}
	fmt.Println()
}

package godrudge

// returns a subslice of the slice arr where start is the index to start on
// removing every other index and returning the resulting slice
func sliceEveryOther[T any](arr []T, start int) []T {
	newIndexCount := 0
	for count := start; count < len(arr); count += 2 {
		arr[newIndexCount] = arr[count]
		newIndexCount++
	}
	arr = arr[:newIndexCount]
	return arr
}

// finds the maximum length out of all inner slices
func determineMaximumColumnSize[T any](columns [][]T) int {
	maxRows := 0
	for _, col := range columns {
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}
	return maxRows
}

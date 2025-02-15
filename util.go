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

package godrudge

import (
	"reflect"
	"testing"
)

func TestSliceEveryOther(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		start    int
		expected []int
	}{
		{"start at 0", []int{1, 2, 3, 4, 5, 6}, 0, []int{1, 3, 5}},
		{"start at 1", []int{1, 2, 3, 4, 5, 6}, 1, []int{2, 4, 6}},
		{"empty slice", []int{}, 0, []int{}},
		{"single element", []int{1}, 0, []int{1}},
		{"start out of bounds", []int{1, 2, 3}, 5, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sliceEveryOther(tt.arr, tt.start)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDetermineMaximumColumnSize(t *testing.T) {
	tests := []struct {
		name     string
		columns  [][]int
		expected int
	}{
		{"multiple columns", [][]int{{1, 2}, {3, 4, 5}, {6}}, 3},
		{"single column", [][]int{{1, 2, 3}}, 3},
		{"empty columns", [][]int{}, 0},
		{"columns with empty slices", [][]int{{}, {}, {}}, 0},
		{"mixed empty and non-empty", [][]int{{1, 2}, {}, {3, 4, 5}}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineMaximumColumnSize(tt.columns)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

package main

import "testing"

func TestIntSliceToString(t *testing.T) {
	tests := []struct {
		inputSlice []int
		expected   string
	}{
		{[]int{1, 2, 3}, "1 2 3"},
		{[]int{1}, "1"},
		{[]int{}, ""},
	}

	for i, test := range tests {
		res := IntSliceToString(test.inputSlice)

		if res != test.expected {
			t.Errorf("test %v, expected %v, got %v", i, test.expected, res)
		}
	}
}

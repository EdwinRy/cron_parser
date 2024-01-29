package main

import (
	"fmt"
	"reflect"
	"testing"
)


func TestCronTaskCompile(t *testing.T) {
	tests := []struct {
		inputCronStr string
		expectedTask CronTask
		expectedError error
	}{
		{ // Check basic use
			"1 1 1 1 1 /usr/bin/find",
			CronTask{
				Minutes:     []int{1},
				Hours:       []int{1},
				DaysOfMonth: []int{1},
				Months:      []int{1},
				DaysOfWeek:  []int{1},
				Command:     "/usr/bin/find",
			},
			nil,
		},
		{ // Check complicated string
			"*/15 0 1,15 * 1-5 /usr/bin/find",
			CronTask{
				Minutes:     []int{0, 15, 30, 45},
				Hours:       []int{0},
				DaysOfMonth: []int{1, 15},
				Months:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				DaysOfWeek:  []int{1, 2, 3, 4, 5},
				Command:     "/usr/bin/find",
			},
			nil,
		},
		{ // Check utf-8
			"1 1 1 1 1 /usr/bin/💀",
			CronTask{
				Minutes:     []int{1},
				Hours:       []int{1},
				DaysOfMonth: []int{1},
				Months:      []int{1},
				DaysOfWeek:  []int{1},
				Command:     "/usr/bin/💀",
			},
			nil,
		},
		{ // Empty string
			"",
			CronTask{},
			fmt.Errorf("failed to tokenize your cron string: didn't find any valid characters in the cron string"),
		},
		{ // No command
			"1 1 1 1 1",
			CronTask{},
			fmt.Errorf("could not parse cron task: couldn't parse time field 5: expected a space after time expression - got EOF() instead after parsing a complete time expression \"1\" for this field"),
		},
		{ // No time fields
			"test",
			CronTask{},
			fmt.Errorf("failed to tokenize your cron string: invalid Tokens found in the cron string: [t e s t], need 5 time space-separated time fields followed by a command"),
		},
		{ // starting with space
			" 1 1 1 1 1",
			CronTask{},
			fmt.Errorf("could not parse cron task: couldn't parse time field 1: couldn't parse time expression"),
		},
		{ // single time field
			"*",
			CronTask{},
			fmt.Errorf("could not parse cron task: couldn't parse time field 1: expected a space after time expression - got EOF() instead after parsing a complete time expression \"*\" for this field"),
		},
		{ // single time field with a space
			"* ",
			CronTask{},
			fmt.Errorf("could not parse cron task: couldn't parse time field 2: couldn't parse time expression"),
		},
		{ // invalid steps value
			"1/10 1 1 1 1 test",
			CronTask{},
			fmt.Errorf("could not parse cron task: couldn't parse time field 1: expected a space after time expression - got Slash(/) instead after parsing a complete time expression \"1\" for this field, it's possible you have provided an invalid value (Number(10)) for the step number or the value range (Number(1))?"),
		},
		{ // invalid range for a field
			"1 40-50 1 1 1 test",
			CronTask{},
			fmt.Errorf("failed to extract valid cron task from syntax: time range needs to be between 0 and 23, got 40 and 50"),
		},
		{ // invalid steps for a field
			"1 40-50/5 1 1 1 test",
			CronTask{},
			fmt.Errorf("failed to extract valid cron task from syntax: steps time range needs to be between 0 and 23, got 40 and 50"),
		},
	}

	for i, test := range tests {
		task, err := CronTaskCompile(test.inputCronStr)

		// Should return the correct error value
		if test.expectedError != nil && err == nil {
			t.Errorf("test %v, expected error \"%v\", got nil", i, test.expectedError)
		} else if test.expectedError == nil && err != nil {
			t.Errorf("test %v, expected no error, got \"%v\"", i, err)
		} else if test.expectedError != nil && err != nil && err.Error() != test.expectedError.Error() {
			t.Errorf("test %v, expected error \"%v\", got \"%v\"", i, test.expectedError, err)
		}
		if err != nil {
			continue
		}

		if !reflect.DeepEqual(*task, test.expectedTask) {
			t.Errorf("test %v, expected %v, got %v", i, test.expectedTask, *task)
		}
	}
}


package main

import (
	"reflect"
	"testing"
)

func TestGetCronTask(t *testing.T) {

	tests := []struct {
		inputNode     AstNode
		expectedTask  CronTask
		expectedError error
	}{
		{
			AstNode{
				NodeType: AstNodeTask,
				Value:    "",
				Children: []AstNode{
					{
						NodeType: AstNodeField,
						Children: []AstNode{
							{
								NodeType: AstTimeExpr,
								Children: []AstNode{
									{
										NodeType: AstTimeRange,
										Children: []AstNode{
											{
												NodeType: AstTimeVal,
												Value:    "1",
											},
											{
												NodeType: AstTimeVal,
												Value:    "10",
											},
										},
									},
								},
							},
						},
					},
					{
						NodeType: AstNodeField,
						Children: []AstNode{
							{
								NodeType: AstTimeExpr,
								Children: []AstNode{
									{
										NodeType: AstTimeSteps,
										Children: []AstNode{
											{
												NodeType: AstTimeRange,
												Children: []AstNode{
													{
														NodeType: AstTimeVal,
														Value:    "1",
													},
													{
														NodeType: AstTimeVal,
														Value:    "10",
													},
												},
											},
											{
												NodeType: AstTimeVal,
												Value:    "2",
											},
										},
									},
								},
							},
						},
					},
					{
						NodeType: AstNodeField,
						Children: []AstNode{
							{
								NodeType: AstTimeExpr,
								Children: []AstNode{
									{
										NodeType: AstTimeSteps,
										Children: []AstNode{
											{
												NodeType: AstTimeRange,
												Children: []AstNode{
													{
														NodeType: AstTimeVal,
														Value:    "1",
													},
													{
														NodeType: AstTimeVal,
														Value:    "10",
													},
												},
											},
											{
												NodeType: AstTimeVal,
												Value:    "2",
											},
										},
									},
									{
										NodeType: AstTimeVal,
										Value:    "4",
									},
								},
							},
						},
					},
					{
						NodeType: AstNodeField,
						Value:    "",
						Children: []AstNode{
							{
								NodeType: AstTimeExpr,
								Value:    "*",
								Children: []AstNode{
									{
										NodeType: AstAsterisk,
										Value:    "*",
									},
								},
							},
						},
					},
					{
						NodeType: AstNodeField,
						Value:    "",
						Children: []AstNode{
							{
								NodeType: AstTimeExpr,
								Value:    "1",
								Children: []AstNode{
									{
										NodeType: AstTimeVal,
										Value:    "1",
									},
								},
							},
						},
					},
					{
						NodeType: AstNodeCommand,
						Value:    "/usr/bin/find",
					},
				},
			},
			CronTask{
				Minutes:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				Hours:       []int{1, 3, 5, 7, 9},
				DaysOfMonth: []int{1, 3, 4, 5, 7, 9},
				Months:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				DaysOfWeek:  []int{1},
				Command:     "/usr/bin/find",
			},
			nil,
		},
	}

	for i, test := range tests {

		res, err := GetCronTask(&test.inputNode)

		if err != test.expectedError {
			t.Errorf("test %v, expected error %v, got %v", i, test.expectedError, err)
		}

		if !reflect.DeepEqual(test.expectedTask.Minutes, res.Minutes) {
			t.Errorf("test %v, expected minutes %v, got %v", i, test.expectedTask.Minutes, res.Minutes)
		}
		if !reflect.DeepEqual(test.expectedTask.Hours, res.Hours) {
			t.Errorf("test %v, expected hours %v, got %v", i, test.expectedTask.Hours, res.Hours)
		}
		if !reflect.DeepEqual(test.expectedTask.DaysOfMonth, res.DaysOfMonth) {
			t.Errorf("test %v, expected days of month %v, got %v", i, test.expectedTask.DaysOfMonth, res.DaysOfMonth)
		}
		if !reflect.DeepEqual(test.expectedTask.Months, res.Months) {
			t.Errorf("test %v, expected months %v, got %v", i, test.expectedTask.Months, res.Months)
		}
		if !reflect.DeepEqual(test.expectedTask.DaysOfWeek, res.DaysOfWeek) {
			t.Errorf("test %v, expected days of week %v, got %v", i, test.expectedTask.DaysOfWeek, res.DaysOfWeek)
		}
	}
}

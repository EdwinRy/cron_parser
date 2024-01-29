package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type CronTask struct {
	Minutes     []int
	Hours       []int
	DaysOfMonth []int
	Months      []int
	DaysOfWeek  []int
	Command     string
}

func (t CronTask) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-14v %v\n", "minute", IntSliceToString(t.Minutes)))
	sb.WriteString(fmt.Sprintf("%-14s %v\n", "hour", IntSliceToString(t.Hours)))
	sb.WriteString(fmt.Sprintf("%-14s %v\n", "day of month", IntSliceToString(t.DaysOfMonth)))
	sb.WriteString(fmt.Sprintf("%-14s %v\n", "month", IntSliceToString(t.Months)))
	sb.WriteString(fmt.Sprintf("%-14s %v\n", "day of week", IntSliceToString(t.DaysOfWeek)))
	sb.WriteString(fmt.Sprintf("%-14s %v\n", "command", t.Command))
	return sb.String()
}

func getTimeVal(fieldValues map[int]struct{}, timeValue int, minVal int, maxVal int) (success bool) {
	if timeValue < minVal || timeValue > maxVal {
		return false
	}
	fieldValues[timeValue] = struct{}{}
	return true
}

func listNodeTimeRange(node AstNode) (start int, end int, err error) {
	if len(node.Children) != 2 {
		err = fmt.Errorf("invalid time range format: %v", node)
		return
	}
	if node.Children[0].NodeType != AstTimeVal || node.Children[1].NodeType != AstTimeVal {
		err = fmt.Errorf("time range needs to consist of 2 integers, got %v and %v", node.Children[0].NodeType, node.Children[1].NodeType)
		return
	}

	start, err = strconv.Atoi(node.Children[0].Value)
	if err != nil {
		return
	}
	end, err = strconv.Atoi(node.Children[1].Value)
	if err != nil {
		return
	}
	if start > end {
		err = fmt.Errorf("time range needs to start from a lower to a higher value, got %v", node.Value)
		return
	}
	return
}

func getTimeRange(fieldValues map[int]struct{}, start int, end int, steps int) {
	for i := start; i <= end; i += steps {
		fieldValues[i] = struct{}{}
	}
}

func getExpressionPart(ast AstNode, fieldValues map[int]struct{}, minVal int, maxVal int) error {
	switch ast.NodeType {
	case AstAsterisk:
		getTimeRange(fieldValues, minVal, maxVal, 1)

	case AstTimeVal:
		timeValue, err := strconv.Atoi(ast.Value)
		if err != nil {
			return fmt.Errorf("steps value needs to be a valid number, got %v", ast.Children[1].Value)
		}
		if timeValue < minVal || timeValue > maxVal {
			return fmt.Errorf("time value needs to be between %v and %v, got %v", minVal, maxVal, timeValue)
		}
		getTimeVal(fieldValues, timeValue, minVal, maxVal)

	case AstTimeRange:
		start, end, err := listNodeTimeRange(ast)
		if err != nil {
			return err
		}
		if start > maxVal || end > maxVal {
			return fmt.Errorf("time range needs to be between %v and %v, got %v and %v", minVal, maxVal, start, end)
		}
		getTimeRange(fieldValues, max(minVal, start), min(maxVal, end), 1)

	case AstTimeSteps:
		if len(ast.Children) != 2 {
			return fmt.Errorf("invalid time steps format: %v", fieldValues)
		}
		if ast.Children[1].NodeType != AstTimeVal {
			return fmt.Errorf("steps value needs to be a valid number, got %v", ast.Children[1].Value)
		}

		steps, err := strconv.Atoi(ast.Children[1].Value)
		if err != nil {
			return fmt.Errorf("steps value needs to be a valid number, got %v", ast.Children[1].Value)
		}

		switch ast.Children[0].NodeType {
		// If it's steps for an asterisk (e.g. */5)
		case AstAsterisk:
			getTimeRange(fieldValues, minVal, maxVal, steps)
		// If it's steps for a range (e.g. 1-10/5)
		case AstTimeRange:
			start, end, err := listNodeTimeRange(ast.Children[0])
			if err != nil {
				return err
			}
			if start > maxVal || end > maxVal {
				return fmt.Errorf("steps time range needs to be between %v and %v, got %v and %v", minVal, maxVal, start, end)
			}
			getTimeRange(fieldValues, max(minVal, start), min(maxVal, end), steps)
		default:
			return fmt.Errorf("invalid time steps format: %v", fieldValues)
		}
	}

	return nil
}

func getCronTimeField(ast AstNode, minVal int, maxVal int) ([]int, error) {
	// Expect a time field to consist of a single expression
	if len(ast.Children) != 1 {
		return nil, fmt.Errorf("invalid time expression format")
	}
	if ast.Children[0].NodeType != AstTimeExpr {
		return nil, fmt.Errorf("invalid time expression format")
	}

	expression := ast.Children[0]
	fieldRange := maxVal - minVal + 1
	fieldValues := make(map[int]struct{}, 0)

	// get values for each part of the expression
	for _, expr := range expression.Children {
		err := getExpressionPart(expr, fieldValues, minVal, maxVal)
		if err != nil {
			return nil, err
		}
		// We're using every possible value for this field, we can return
		if len(fieldValues) == fieldRange {
			break
		}
	}

	// convert map keys to a slice
	values := make([]int, len(fieldValues))
	i := 0
	for val := range fieldValues {
		values[i] = val
		i++
	}
	sort.Ints(values)
	return values, nil
}

func GetCronTask(ast *AstNode) (*CronTask, error) {
	// Expect 6 children: 5 time fields and 1 command
	if len(ast.Children) != 6 {
		return nil, fmt.Errorf("invalid cron format, expected 5 time fields and a command")
	}
	for i := 0; i < 5; i++ {
		if ast.Children[i].NodeType != AstNodeField {
			return nil, fmt.Errorf("expected 5 time fields followed by a command")
		}
	}
	if ast.Children[5].NodeType != AstNodeCommand {
		return nil, fmt.Errorf("expected after 5 time fields")
	}

	task := CronTask{}

	// Get trigger times from time fields
	var err error
	task.Minutes, err = getCronTimeField(ast.Children[0], 0, 59)
	if err != nil {
		return nil, err
	}
	task.Hours, err = getCronTimeField(ast.Children[1], 0, 23)
	if err != nil {
		return nil, err
	}
	task.DaysOfMonth, err = getCronTimeField(ast.Children[2], 1, 31)
	if err != nil {
		return nil, err
	}
	task.Months, err = getCronTimeField(ast.Children[3], 1, 12)
	if err != nil {
		return nil, err
	}
	task.DaysOfWeek, err = getCronTimeField(ast.Children[4], 0, 6)
	if err != nil {
		return nil, err
	}

	task.Command = ast.Children[5].Value
	return &task, nil
}

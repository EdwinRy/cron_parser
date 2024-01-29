package main

import (
	"encoding/json"
	"fmt"
)

type AstNodeType byte

const (
	AstNil AstNodeType = iota
	AstNodeTask
	AstNodeField
	AstNodeCommand
	AstTimePart
	AstTimeSteps
	AstTimeRange
	AstTimeVal
	AstAsterisk
	AstTimeExpr
)

func (t AstNodeType) String() string {
	switch t {
	case AstNil:
		return "Nil"
	case AstNodeTask:
		return "Task"
	case AstNodeField:
		return "Field"
	case AstNodeCommand:
		return "Command"
	case AstTimePart:
		return "TimePart"
	case AstTimeSteps:
		return "TimeSteps"
	case AstTimeRange:
		return "TimeRange"
	case AstTimeVal:
		return "TimeVal"
	case AstAsterisk:
		return "Asterisk"
	case AstTimeExpr:
		return "TimeExpr"
	default:
		return "Unknown"
	}
}

func (t AstNodeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

type AstNode struct {
	NodeType AstNodeType
	Value    string
	Children []AstNode
}

func (n AstNode) String() string {
	return fmt.Sprintf("{\"%v|%v\":%v},", n.NodeType, n.Value, n.Children)
}

func lookahead(tokens []Token, tokensPtr int, expected []TokenType) bool {
	if len(tokens) < tokensPtr+len(expected) {
		return false
	}
	for i, tokType := range expected {
		if tokens[tokensPtr+i].tokType != tokType {
			return false
		}
	}
	return true
}

func parseTimeVal(tokens []Token, tokensPtr int) (node AstNode, newTokenPtr int, success bool) {

	switch tokens[tokensPtr].tokType {

	case TokenAsterisk:
		node = AstNode{AstAsterisk, "*", []AstNode{}}
		newTokenPtr = tokensPtr + 1
		success = true
		return

	case TokenNumber:
		node = AstNode{AstTimeVal, string(tokens[tokensPtr].value), []AstNode{}}
		newTokenPtr = tokensPtr + 1
		success = true
		return

	default:
		newTokenPtr = tokensPtr
		success = false
		return
	}
}

func parseTimeRange(tokens []Token, tokensPtr int) (node AstNode, newTokenPtr int, success bool) {
	tkptr := tokensPtr

	// look ahead to see if the range production matches
	if !lookahead(tokens, tkptr, []TokenType{TokenNumber, TokenDash, TokenNumber}) {
		success = false
		newTokenPtr = tokensPtr
		return
	}

	// Parse the start value
	rangeFrom := AstNode{AstTimeVal, string(tokens[tkptr].value), []AstNode{}}
	tkptr++

	// Skip the dash
	tkptr++

	// Parse the end value
	rangeTo := AstNode{AstTimeVal, string(tokens[tkptr].value), []AstNode{}}
	tkptr++

	rangeValue := fmt.Sprintf("%v-%v", rangeFrom.Value, rangeTo.Value)
	timeRange := AstNode{AstTimeRange, rangeValue, []AstNode{rangeFrom, rangeTo}}
	return timeRange, tkptr, true
}

func parseTimeSteps(tokens []Token, tokensPtr int) (node AstNode, newTokenPtr int, success bool) {
	tkptr := tokensPtr

	// Parse either an Asterisk or a number range
	var stepsRange AstNode
	if tokens[tkptr].tokType == TokenAsterisk {
		stepsRange = AstNode{AstAsterisk, "*", []AstNode{}}
		tkptr++
	} else {
		var gotRange bool
		stepsRange, tkptr, gotRange = parseTimeRange(tokens, tkptr)
		if !gotRange {
			newTokenPtr = tokensPtr
			success = false
			return
		}
	}

	// Expect the slash
	if tokens[tkptr].tokType != TokenSlash {
		newTokenPtr = tokensPtr
		success = false
		return
	}
	tkptr++

	// Parse the step value
	if tokens[tkptr].tokType != TokenNumber {
		newTokenPtr = tokensPtr
		success = false
		return
	}
	stepVal := AstNode{AstTimeVal, string(tokens[tkptr].value), []AstNode{}}
	tkptr++

	// Return the time steps node
	stepsValue := fmt.Sprintf("%v/%v", stepsRange.Value, stepVal.Value)
	timeSteps := AstNode{AstTimeSteps, stepsValue, []AstNode{stepsRange, stepVal}}
	return timeSteps, tkptr, true
}

func parseTimePart(tokens []Token, tokensPtr int) (node AstNode, newTokenPtr int, success bool) {
	tkptr := tokensPtr

	var partNode AstNode
	var parseSuccess bool

	// Try parsing a field with steps first as it's the longest match
	partNode, tkptr, parseSuccess = parseTimeSteps(tokens, tkptr)
	if parseSuccess {
		return partNode, tkptr, true
	}

	// Try parsing a time range next
	partNode, tkptr, parseSuccess = parseTimeRange(tokens, tkptr)
	if parseSuccess {
		return partNode, tkptr, true
	}

	// Try parsing a single time value
	partNode, tkptr, parseSuccess = parseTimeVal(tokens, tkptr)
	if parseSuccess {
		return partNode, tkptr, true
	}

	// If we got here, we couldn't parse a time part
	newTokenPtr = tokensPtr
	success = false
	return
}

func parseTimeExpr(tokens []Token, tokensPtr int) (node AstNode, newTokenPtr int, success bool) {
	tkptr := tokensPtr

	timeExpr := AstNode{AstTimeExpr, "", []AstNode{}}

	// keep parsing time parts until we don't see a comma at the end
	for {
		var timePart AstNode
		var gotPart bool
		timePart, tkptr, gotPart = parseTimePart(tokens, tkptr)
		if !gotPart {
			break
		}
		timeExpr.Children = append(timeExpr.Children, timePart)
		timeExpr.Value += timePart.Value

		// stop parsing if we don't have a comma
		if tokens[tkptr].tokType != TokenComma {
			break
		} else {
			timeExpr.Value += ","
		}
		tkptr++
	}

	// If we didn't parse any time parts, expression is not valid
	if len(timeExpr.Children) == 0 {
		newTokenPtr = tokensPtr
		success = false
		return
	}

	return timeExpr, tkptr, true
}

func parseTimeField(tokens []Token, tokensPtr int) (node AstNode, newTokenPtr int, err error) {
	tkptr := tokensPtr

	// Parse the time expression
	var timeExpr AstNode
	var gotExpr bool
	timeExpr, tkptr, gotExpr = parseTimeExpr(tokens, tkptr)
	if !gotExpr {
		newTokenPtr = tokensPtr
		err = fmt.Errorf("couldn't parse time expression")
		return
	}

	// Expect a space
	if tokens[tkptr].tokType != TokenSpace {
		newTokenPtr = tokensPtr

		err = fmt.Errorf(
			"expected a space after time expression - got %v instead after "+
				"parsing a complete time expression \"%v\" for this field",
			tokens[tkptr].String(), timeExpr.Value)

		// provide additional error context for step syntax
		if tokens[tkptr].tokType == TokenSlash && tkptr+1 < len(tokens) && tkptr > 0 {
			err = fmt.Errorf(
				"%v, it's possible you have provided an invalid value (%v) for the step number or "+
					"the value range (%v)?",
				err, tokens[tkptr+1].String(), tokens[tkptr-1].String())
		}
		return
	}
	tkptr++

	return AstNode{AstNodeField, timeExpr.Value, []AstNode{timeExpr}}, tkptr, nil
}

func parseTask(tokens []Token, tokensPtr int) (node AstNode, newTokenPtr int, err error) {
	tkptr := tokensPtr
	task := AstNode{AstNodeTask, "", []AstNode{}}

	// Parse 5 time fields
	for i := 0; i < 5; i++ {
		var field AstNode
		field, tkptr, err = parseTimeField(tokens, tkptr)
		if err != nil {
			newTokenPtr = tokensPtr
			err = fmt.Errorf("couldn't parse time field %v: %v", i+1, err)
			return
		}
		task.Children = append(task.Children, field)
	}

	// Parse the command
	if tokens[tkptr].tokType != TokenCommand {
		newTokenPtr = tokensPtr
		err = fmt.Errorf("expected 5 space-separated time fields followed by a command")
		return
	}
	command := AstNode{AstNodeCommand, string(tokens[tkptr].value), []AstNode{}}
	tkptr++
	task.Children = append(task.Children, command)

	return task, tkptr, nil
}

func Parse(tokens []Token) (*AstNode, error) {
	root, tokenPtr, err := parseTask(tokens, 0)
	if err != nil {
		return nil, err
	}
	if tokenPtr+1 != len(tokens) {
		return nil, fmt.Errorf("incorrect format: expected 5 space-separated time fields followed by a command")
	}
	return &root, nil
}

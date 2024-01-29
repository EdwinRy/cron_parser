package main

import (
	"reflect"
	"testing"
)

func TestLookahead(t *testing.T) {
	tests := []struct {
		inputTokens   []Token
		inputPtr      int
		inputExpected []TokenType
		successMatch  bool
	}{
		{ // matching single token
			[]Token{{TokenAsterisk, []rune("*")}},
			0,
			[]TokenType{TokenAsterisk},
			true,
		},
		{ // matching multiple tokens
			[]Token{{TokenAsterisk, []rune("*")}, {TokenNumber, []rune("432")}},
			0,
			[]TokenType{TokenAsterisk, TokenNumber},
			true,
		},
		{ // non-matching multiple tokens
			[]Token{{TokenAsterisk, []rune("*")}, {TokenNumber, []rune("432")}},
			0,
			[]TokenType{TokenAsterisk, TokenAsterisk},
			false,
		},
	}

	for i, test := range tests {
		res := lookahead(test.inputTokens, test.inputPtr, test.inputExpected)

		if res != test.successMatch {
			t.Errorf("test %v, expected %v, got %v", i, test.successMatch, res)
		}
	}
}

func TestParseTimeVal(t *testing.T) {
	tests := []struct {
		inputTokens []Token
		inputPtr    int

		expectedNode    AstNode
		expectedPtr     int
		expectedSuccess bool
	}{
		{ // Check for asterisk
			[]Token{{TokenAsterisk, []rune("*")}},
			0,
			AstNode{AstAsterisk, "*", []AstNode{}},
			1,
			true,
		},
		{ // Check for numbers
			[]Token{{TokenNumber, []rune("3")}},
			0,
			AstNode{AstTimeVal, "3", []AstNode{}},
			1,
			true,
		},
		{ // Reject invalid range values
			[]Token{{TokenComma, []rune(",")}},
			0,
			AstNode{AstNodeCommand, ",", []AstNode{}},
			0,
			false,
		},
	}

	for i, test := range tests {
		node, newPtr, success := parseTimeVal(test.inputTokens, test.inputPtr)

		// Should return the correct success value
		if success != test.expectedSuccess {
			t.Errorf("test %v, expected success to be %v, got %v", i, test.expectedSuccess, success)
		}

		if newPtr != test.expectedPtr {
			t.Errorf("test %v, expected pointer %v, got %v", i, test.expectedPtr, newPtr)
		}

		// No need to test the value
		if success == false {
			continue
		}

		if !reflect.DeepEqual(node, test.expectedNode) {
			t.Errorf("test %v, expected nodes to be %v, got %v", i, test.expectedNode, node)
		}
	}
}

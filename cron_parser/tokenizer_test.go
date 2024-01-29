package main

import (
	"fmt"
	"testing"
)

func TestTokenizeNumber(t *testing.T) {

	tests := []struct {
		input           []rune
		expectedValue   string
		expectedEnd     int
		expectedSuccess bool
	}{
		{make([]rune, 0), "", 0, false},
		{[]rune("1"), "1", 1, true},
		{[]rune("1234"), "1234", 4, true},
		{[]rune("+-="), "", 0, false},
	}

	for i, test := range tests {
		cronStrRunes := []rune(test.input)
		token, end, success := TokenizeNumber(&cronStrRunes, 0)

		// Should return the correct success value
		if success != test.expectedSuccess {
			t.Errorf("test %v, expected success to be %v, got %v", i, test.expectedSuccess, success)
		}

		// No need to test the value
		if success == false {
			continue
		}

		// Should return a TokenNumber
		if token.tokType != TokenNumber {
			t.Errorf("test %v, expected TokenNumber, got %v", i, token.tokType)
		}

		// Should tokenize the value correctly
		if string(token.value) != test.expectedValue {
			t.Errorf("test %v, expected %v, got %v", i, test.expectedValue, string(token.value))
		}

		// Expect to parse the whole number available
		if end != test.expectedEnd {
			t.Errorf("test %v, expected end to be %v, got %v", i, test.expectedEnd, end)
		}
	}
}

func TestTokenize(t *testing.T) {

	var tests = []struct {
		input          string
		expectedTokens []Token
		expectedError  error
	}{
		{"", []Token{{TokenEOF, []rune("")}}, fmt.Errorf("didn't find any valid characters in the cron string")},
		{"1", []Token{{TokenNumber, []rune("1")}, {TokenEOF, []rune("")}}, nil},
		{"123", []Token{{TokenNumber, []rune("123")}, {TokenEOF, []rune("")}}, nil},
		{"*", []Token{{TokenAsterisk, []rune("*")}, {TokenEOF, []rune("")}}, nil},
		{"*,-/54", []Token{
			{TokenAsterisk, []rune("*")},
			{TokenComma, []rune(",")},
			{TokenDash, []rune("-")},
			{TokenSlash, []rune("/")},
			{TokenNumber, []rune("54")},
			{TokenEOF, []rune("")},
		}, nil},
		{"&", []Token{{TokenEOF, []rune("")}}, fmt.Errorf("invalid Tokens found in the cron string: [&], need 5 time space-separated time fields followed by a command")},
	}

	for i, test := range tests {
		tokens, err := Tokenize(test.input)

		// Should return the correct error value
		if test.expectedError != nil && err == nil {
			t.Errorf("test %v, expected error %v, got nil", i, test.expectedError)
		} else if test.expectedError == nil && err != nil {
			t.Errorf("test %v, expected no error, got %v", i, err)
		} else if test.expectedError != nil && err != nil && err.Error() != test.expectedError.Error() {
			t.Errorf("test %v, expected error %v, got %v", i, test.expectedError, err)
		}
		if err != nil {
			continue
		}

		// Should return the correct tokens
		for j, token := range tokens {
			if token.tokType != test.expectedTokens[j].tokType {
				t.Errorf("test %v, expected type %v, got %v",
					i, test.expectedTokens[j].tokType, token.tokType)
			}

			if string(token.value) != string(test.expectedTokens[j].value) {
				t.Errorf("test %v, expected value %v, got %v",
					i, string(test.expectedTokens[j].value), string(token.value))
			}
		}
	}
}

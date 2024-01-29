package main

import (
	"fmt"
	"unicode"
)

type TokenType byte

const (
	TokenAsterisk TokenType = iota
	TokenComma
	TokenDash
	TokenSlash

	TokenNumber
	TokenSpace
	TokenCommand

	TokenEOF
)

func (t TokenType) String() string {
	switch t {
	case TokenAsterisk:
		return "Asterisk"
	case TokenComma:
		return "Comma"
	case TokenDash:
		return "Dash"
	case TokenSlash:
		return "Slash"
	case TokenNumber:
		return "Number"
	case TokenSpace:
		return "Space"
	case TokenCommand:
		return "Command"
	case TokenEOF:
		return "EOF"
	default:
		return "Unknown"
	}
}

type Token struct {
	tokType TokenType
	value   []rune
}

func (t Token) String() string {
	return fmt.Sprintf("%v(%v)", t.tokType, string(t.value))
}

func TokenizeNumber(cronStrRunes *[]rune, start int) (tk Token, end int, success bool) {

	end = start

	// Check for out of bounds
	if len(*cronStrRunes) <= start {
		success = false
		return
	}

	// Keep fetching all digits until we hit a non-digit
	for ; end < len(*cronStrRunes); end++ {
		if !unicode.IsDigit((*cronStrRunes)[end]) {
			break
		}
	}

	// Check if we fetched any digits
	if end == start {
		success = false
		return
	}

	tk = Token{TokenNumber, (*cronStrRunes)[start:end]}
	success = true
	return
}

func Tokenize(cronStr string) ([]Token, error) {

	tokens := make([]Token, 0)
	invalidTokens := make([]string, 0)
	space_count := 0

	invalid := false

	runes := []rune(cronStr)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		// If we've seen 5 spaces, we're done and the rest is the command
		if space_count >= 5 && i < len(runes) {
			tokens = append(tokens, Token{TokenCommand, runes[i:]})
			break
		}

		switch {
		case char == '*':
			tokens = append(tokens, Token{TokenAsterisk, runes[i : i+1]})
		case char == ',':
			tokens = append(tokens, Token{TokenComma, runes[i : i+1]})
		case char == '-':
			tokens = append(tokens, Token{TokenDash, runes[i : i+1]})
		case char == '/':
			tokens = append(tokens, Token{TokenSlash, runes[i : i+1]})
		case char == ' ':
			tokens = append(tokens, Token{TokenSpace, runes[i : i+1]})
			space_count++
		case unicode.IsDigit(char):
			numToken, end, success := TokenizeNumber(&runes, i)
			if success {
				tokens = append(tokens, numToken)
			}
			i = end - 1
		default:
			invalidTokens = append(invalidTokens, string(char))
			invalid = true
		}
	}

	if invalid {
		return nil, fmt.Errorf(
			"invalid Tokens found in the cron string: %v, need 5 time space-separated time fields followed by a command",
			invalidTokens)
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("didn't find any valid characters in the cron string")
	}

	tokens = append(tokens, Token{TokenEOF, []rune("")})

	return tokens, nil
}

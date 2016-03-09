package parse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type parseState uint

const (
	None parseState = iota
	String
	Hyphen
)

func ParseInputString(input string) []string {
	var words []string

	for len(input) > 0 {
		nextUtf8Token, i := getNextUtf8Token(input)
		if nextUtf8Token != "" {
			words = append(words, nextUtf8Token)
		}
		input = input[i:]
	}
	return words
}

func getNextUtf8Token(input string) (string, int) {
	// Future improvements: Parse languages that use non-ASCII characters.
	var state parseState = None
	token := ""
	i := 0

loop:
	for w := 0; i < len(input); i += w {
		runeValue, width := utf8.DecodeRuneInString(input[i:])
		switch {
		// Append alphanumeric characters.
		case (runeValue >= 'a' && runeValue <= 'z') || (runeValue >= 'A' && runeValue <= 'Z'):
			token += string(runeValue)
			state = String
		case runeValue >= '0' && runeValue <= '9':
			token += string(runeValue)
			state = String
		// Punctuation
		case unicode.IsPunct(runeValue):
			if runeValue == '-' {
				// If a single hyphen is encountered mid-word, treat it as part of the
				// current word. Otherwise treat it as punctuation.
				if state == Hyphen {
					break loop
				}
				token += string(runeValue)
				state = Hyphen
			} else if runeValue == '\'' && state == String {
				// Allow apostrophes inside of words.
				token += string(runeValue)
			} else if state == String {
				// Break on punctuation only if parsing mid-word.
				break loop
			}
		// Break on whitespace only if parsing mid-word.
		case unicode.IsSpace(runeValue):
			if state == String {
				break loop
			}
		// Append any other characters.
		default:
			token += string(runeValue)
		}
		w = width
	}
	// Trim any whitespace and hyphens from the output.
	token = strings.Trim(token, "-")
	token = strings.Trim(token, " ")
	return token, i
}

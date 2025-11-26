package commands

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Type represents a tokenizing type
type Type int

// Token represents an argument arbitrary within the source
type Token struct {
	Type
	Literal any
}

const (
	STRING Type = iota
	NUMBER
	BOOL
	SPACE
	BLANK

	// ANY supports all the token types
	ANY
)

// Tokenize implements a base tokenization process
func Tokenize(text string, tokens []*Token) ([]*Token, error) {
	index := strings.Split(text, "")
	for pos := 0; pos < len(index); pos++ {
		switch rune(index[pos][0]) {

		case ' ':
			if pos == 0 {
				tokens = append(tokens, &Token{Literal: "", Type: BLANK})
			}
			
			tokens = append(tokens, &Token{
				Literal: index[pos],
				Type: SPACE,
			})

			if pos + 1 >= len(index) {
				tokens = append(tokens, &Token{Literal: "", Type: BLANK})
				return tokens, nil
			}

		case '"', '\'':
			tokens = append(tokens, &Token{
				Type: STRING,
				Literal: index[pos],
			})

			for _, char := range index[pos + 1:] {
				tokens[len(tokens) - 1].Literal = fmt.Sprintf("%s" + char, tokens[len(tokens) - 1].Literal)
				if char[0] == index[pos][0] {
					break
				}
			}

			pos += len(fmt.Sprint(tokens[len(tokens) - 1].Literal)) - 1
			if len(index) <= pos {
				break
			}

			tokens[len(tokens) - 1].Literal = strings.ReplaceAll(fmt.Sprint(tokens[len(tokens) - 1].Literal), index[pos], "")

		default:
			if unicode.IsDigit(rune(index[pos][0])) {
				tokens = append(tokens, &Token{
					Type: NUMBER,
					Literal: index[pos],
				})
	
				for _, char := range index[pos + 1:] {
					if !unicode.IsDigit(rune(char[0])) && rune(char[0]) != '.' {
						break
					}

					tokens[len(tokens) - 1].Literal = fmt.Sprintf("%s" + char, tokens[len(tokens) - 1].Literal)
				}
	
				pos += len(fmt.Sprint(tokens[len(tokens) - 1].Literal)) - 1
				if strings.Contains(fmt.Sprint(tokens[len(tokens) - 1].Literal), ".") {
					tokens[len(tokens) - 1].Type = STRING
					continue
				}

				digit, err := strconv.Atoi(fmt.Sprint(tokens[len(tokens) - 1].Literal))
				if err != nil {
					return nil, err
				}

				tokens[len(tokens) - 1].Literal = digit
			} else if unicode.IsLetter(rune(index[pos][0])) {
				tokens = append(tokens, &Token{
					Type: STRING,
					Literal: index[pos],
				})
	
				for _, char := range index[pos + 1:] {
					if unicode.IsSpace(rune(char[0])) || rune(char[0]) == '"' {
						break
					}

					tokens[len(tokens) - 1].Literal = fmt.Sprintf("%s" + char, tokens[len(tokens) - 1].Literal)
				}

				pos += len(fmt.Sprint(tokens[len(tokens) - 1].Literal)) - 1
				switch fmt.Sprint(tokens[len(tokens) - 1].Literal) {
				case "true", "false", "on", "off":
					operator, err := strconv.ParseBool(fmt.Sprint(tokens[len(tokens) - 1].Literal))
					if err != nil {
						return nil, err
					}

					tokens[len(tokens) - 1].Type = BOOL
					tokens[len(tokens) - 1].Literal = operator
				}
			}
		}
	}

	return tokens, nil
}

// ToString will return the string representation of the token
func (t *Token) ToString() string {
	return fmt.Sprint(t.Literal)
}

// Join will return the string representation of the tokens with the sep
func ToString(tokens []*Token) []string {
	dest := make([]string, 0)
	for _, token := range tokens {
		dest = append(dest, token.ToString())
	}

	return dest
}
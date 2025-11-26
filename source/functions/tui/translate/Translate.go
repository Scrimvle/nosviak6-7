package translate

import (
	"Nosviak4/source/swash"
	"fmt"
	"strings"
)

type Tag struct {
	Line           int
	Name           string
	Attr           []*swash.Token
	Content        string
}

type Translator struct {
	Tags      []Tag
	tokenizer *swash.Tokenizer
}

// NewTranslator creates a brand new translator instance
func NewTranslator(source string) *Translator {
	return &Translator{
		Tags:      make([]Tag, 0),
		tokenizer: swash.NewTokenizer(source, false),
	}
}

// Parse will initialize the new translator instance
func (translator *Translator) Parse() error {
	var tail *Tag = nil

	for {
		token, ok := translator.tokenizer.PeekNext()
		if token == nil || !ok {
			break
		}

		switch token.TokenType {

		case swash.LESSTHAN:
			nextInline, ok := translator.tokenizer.PeekNext()
			if !ok || nextInline == nil {
				return fmt.Errorf("EOF")
			}

			switch nextInline.TokenType {

			case swash.DIVIDE:
				assign, ok := translator.tokenizer.PeekNext()
				if !ok || assign == nil || assign.TokenType != swash.INDENT {
					return fmt.Errorf("EOF")
				}

				closure, ok := translator.tokenizer.PeekNext()
				if !ok || closure == nil || closure.TokenType != swash.GREATERTHAN {
					return fmt.Errorf("EOF")
				}

				/* rebuilds the tail */
				translator.Tags = append(translator.Tags, *tail)
				tail = &Tag{
					Attr: make([]*swash.Token, 0),
				}

			case swash.INDENT:
				tail = &Tag{
					Attr: make([]*swash.Token, 0),
				}

				for {
					data, ok := translator.tokenizer.PeekNext()
					if !ok || data == nil {
						break
					}

					if data.TokenType == swash.GREATERTHAN {
						break
					}

					tail.Attr = append(tail.Attr, data)
				}

				tail.Line = token.TokenLine + 1
				tail.Name = nextInline.TokenLiteral
				tail.Content = strings.Split(translator.tokenizer.Line(), "</")[0]
			}
		}
	}

	return nil
}

// Process will process all the fields attached onto the current tag
func (tag *Tag) Process() map[string]any {
	mimic := &swash.Tokenizer{
		TokenizerStream:      tag.Attr,
	}

	kvs := make(map[string]any)

	for {
		index, ok := mimic.PeekNext()
		if !ok || index == nil || index.TokenType != swash.INDENT {
			break
		}

		equalizer, ok := mimic.PeekNext()
		if !ok || equalizer == nil || equalizer.TokenType != swash.EQUAL {
			break
		}

		value, ok := mimic.PeekNext()
		if !ok || value == nil {
			break
		}

		kvs[index.TokenLiteral] = value.TokenType.Go(value.TokenLiteral)
	}

	return kvs
}
package translate_test

import (
	"Nosviak4/source/functions/tui/translate"
	"fmt"
	"strings"
	"testing"
)

func TestTranslate(t *testing.T) {
	translate := translate.NewTranslator("<text>┌─────────────────────────────────────────────────────────────────────────────┐</text>")
	content, err := translate.Analyze()
	if err != nil {
		panic(err)
	}

	fmt.Println(strings.Join(content, "\r\n"))
}

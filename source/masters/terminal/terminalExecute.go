package terminal

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions/iplookup"
	"Nosviak4/source/swash"
	"Nosviak4/source/swash/evaluator"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sync"
)

var (
	// FakeSlaves controls our amount of fake connections established
	FakeSlaves int = 0
	Mutex      sync.Mutex
)

// ExecuteString will execute the swash evaluator through the terminal frontend
func (t *Terminal) ExecuteString(vars map[string]any, content string) error {
	tokenizer := swash.NewTokenizer(content, true).Strip()
	if err := tokenizer.Parse(); err != nil {
		return err
	}

	vars["clear"] = t.ClearString
	vars["exec"] = func(cmd string) {
		go func() {
			t.Signal.Queue <- append([]byte(cmd), 13)
		}()
	}

	return ExecuteStringToWriter(t, vars, tokenizer)
}

// ExecuteBranding will directly interact with the VIEWS map
func (t *Terminal) ExecuteBranding(vars map[string]any, content ...string) error {
	branding, ok := source.OPTIONS.Config.Renders[filepath.Join(content...)]
	if !ok || branding == nil {
		return fmt.Errorf("branding file not found (%s)", filepath.Join(content...))
	}

	if err := t.ExecuteString(vars, string(branding)); err != nil {
		return fmt.Errorf("error occurred while executing (%s)[%v]", filepath.Join(content...), err)
	}

	return nil
}

// ExecuteBrandingToString will execute the branding file into a string directly
func (t *Terminal) ExecuteBrandingToString(vars map[string]any, content ...string) (string, error) {
	branding, ok := source.OPTIONS.Config.Renders[filepath.Join(content...)]
	if !ok || branding == nil {
		return "", nil
	}

	tokenizer := swash.NewTokenizer(string(branding), true).Strip()
	if err := tokenizer.Parse(); err != nil {
		return "", err
	}

	vars["clear"] = t.ClearString
	vars["exec"] = func(cmd string) {
		t.Signal.Queue <- append([]byte(cmd), 13)
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	if err := ExecuteStringToWriter(buffer, vars, tokenizer); err != nil {
		return "", fmt.Errorf("error occurred while executing (%s)[%v]", filepath.Join(content...), err)
	}

	return buffer.String(), nil
}

// ExecuteStringToString will execute the text into a string directly
func (t *Terminal) ExecuteStringToString(vars map[string]any, content string) (string, error) {
	tokenizer := swash.NewTokenizer(content, true).Strip()
	if err := tokenizer.Parse(); err != nil {
		return "", err
	}

	vars["clear"] = t.ClearString
	vars["exec"] = func(cmd string) {
		t.Signal.Queue <- append([]byte(cmd), 13)
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	if err := ExecuteStringToWriter(buffer, vars, tokenizer); err != nil {
		return "", fmt.Errorf("error occurred while executing (%s)[%v]", content, err)
	}

	return buffer.String(), nil
}

// ExecuteBrandingToString will execute the branding to the string output
func ExecuteBrandingToString(vars map[string]any, content ...string) (string, error) {
	branding, ok := source.OPTIONS.Config.Renders[filepath.Join(content...)]
	if !ok || branding == nil {
		return "", fmt.Errorf("branding file not found (%s)", filepath.Join(content...))
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	tokenizer := swash.NewTokenizer(string(branding), true).Strip()
	if err := tokenizer.Parse(); err != nil {
		return "", err
	}

	if err := ExecuteStringToWriter(buf, vars, tokenizer); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ExecuteStringToWriter will append into the settings basis
func ExecuteStringToWriter(wr io.Writer, vars map[string]any, tokenizer *swash.Tokenizer) error {
	eval := evaluator.NewEvaluator(tokenizer, wr)
	eval.Memory.AllocateMap(map[string]any{
		"version": source.VERSION,
		"cnc": func() string {
			return source.OPTIONS.String("app_name")
		},

		"iplookup": func(ip string) *iplookup.Internet {
			internet, err := iplookup.Lookup(ip)
			if err != nil {
				return &iplookup.Internet{}
			}

			return internet
		},

		"uptime":      database.DB.Connected.Unix(),
		"fake_slaves": FakeSlaves,
	})

	if err := eval.Memory.AllocateMap(vars); err != nil {
		return err
	}

	return eval.Execute()
}

func (t *Terminal) ExecuteStringToWriter(vars map[string]any, tokenizer *swash.Tokenizer) error {
	vars["clear"] = t.ClearString
	return ExecuteStringToWriter(t, vars, tokenizer)
}

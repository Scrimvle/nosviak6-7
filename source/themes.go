package source

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	// Themes is a list of all the themes
	Themes map[string]*Theme = make(map[string]*Theme)

	// ErrInvalidTheme is returned when a theme is invalid
	ErrInvalidTheme error = errors.New("invalid theme")
)

// themes represents a dialog for the theme
type Theme struct {
	Name             string
	Hidden           bool     `toml:"hidden"`
	Description      string   `toml:"description"`
	Branding         string   `toml:"branding"`
	DisabledCommands []string `toml:"disabledCommands"`
	Glamour     struct {
		Enabled  bool    `toml:"enabled"`
		Colours  [][]int `toml:"colours"`
		TabMenus bool    `toml:"tab_menus"`
	} `toml:"glamour"`

	/* CustomCommands validates all the customCommands which are loaded on a per theme basis, also includes attacks */
	CustomCommands []any
}

// GetTheme returns the theme which is indexed via the Themes
func GetTheme(theme string) (*Theme, error) {
	index, ok := Themes[theme]
	if !ok {
		return Themes[OPTIONS.String("default_theme")], nil
	}

	return index, nil
}

// configureThemes configures the themes handling
func configureThemes(bytes []byte, core map[string]any) (map[string]any, error) {
	destRead := make(map[string]any)
	if err := toml.Unmarshal(bytes, &destRead); err != nil {
		return nil, err
	}

	/* options represents any configs inside the file */
	options := make(map[string]any)

	/* themes represents any theme configs inside the file */
	Themes = make(map[string]*Theme)

	for key, val := range destRead {
		context, ok := val.(map[string]interface{})
		if !ok || context == nil {
			options[key] = val
			continue
		}

		wr, err := json.Marshal(context)
		if err != nil {
			return nil, err
		}

		dest := new(Theme)
		if err := json.Unmarshal(wr, &dest); err != nil {
			return nil, err
		}

		dest.Name = key
		Themes[key] = dest
		if err := Themes[key].handleThemesCustomCommands(); err != nil {
			return nil, err
		}
	}

	return options, nil
}

// GetCustomCommand will attempt to execute a custom command based on it's header information
func (t *Theme) GetCustomCommand(alias string) any {
	for _, command := range t.CustomCommands {
		switch index := command.(type) {

		case *Text:
			if !strings.EqualFold(alias, index.Name) {
				continue
			}

			return index
		}
	}

	return nil
}

package source

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source/swash"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Text struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Permissions string `toml:"permissions"`
	Themes      string `toml:"themes"`
	Tokenizer   *swash.Tokenizer
}

type Bin struct {
	Name        []string `toml:"name"`
	Description string   `toml:"description"`
	Permissions []string `toml:"permissions"`
	Runtime     string   `toml:"runtime"`
	Env         []string `toml:"vars"`
}

// handleThemesCustomCommands will load all the different types of commands based on the scope.
func (t *Theme) handleThemesCustomCommands() error {
	t.CustomCommands = make([]any, 0)
	return filepath.Walk(filepath.Join(ASSETS, COMMANDS), func(path string, info fs.FileInfo, err error) error {
		index := path[len(filepath.Join(ASSETS, COMMANDS)):]
		if err != nil || strings.Count(index, string(filepath.Separator)) < 2 {
			return err
		}

		switch strings.Split(path[len(filepath.Join(ASSETS, COMMANDS)):][1:], string(filepath.Separator))[0] {

		case "configs":
			return nil

		// text command
		case "text":
			content, err := os.ReadFile(path)
			if err != nil || len(strings.Split(string(content), "============================ START ============================")) <= 1 {
				return err
			}

			var metaData = new(Text)
			if err := toml.Unmarshal([]byte(strings.Split(string(content), "============================ START ============================")[0]), &metaData); err != nil || len(metaData.Themes) > 0 && metaData.Themes != t.Name {
				if err == nil {
					return nil
				}

				LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "non-fatal error occurred with %s: %v", path, err)
				return err
			}

			metaData.Tokenizer = swash.NewTokenizer(strings.Join(strings.Split(strings.Split(string(content), "============================ START ============================")[1], "\r\n")[1:], "\r\n"), true).Strip()
			if err := metaData.Tokenizer.Parse(); err != nil {
				LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "non-fatal error occurred with %s: %v", path, err)
				return err
			}

			t.CustomCommands = append(t.CustomCommands, metaData)
			return nil

		// bin command
		case "bin":
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var config *Bin = new(Bin)
			if err := toml.Unmarshal(content, &config); err != nil {
				return nil
			}

			t.CustomCommands = append(t.CustomCommands, config)
			return nil

		default:
			return fmt.Errorf("unknown path provided inside a dynamic dir: %s", path)
		}
	})
}

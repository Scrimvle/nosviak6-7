package commands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"errors"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

// ErrCommandNotFound is returned when the command is not found
var ErrCommandNotFound = errors.New("command not found")

// CommandsLog is our log of all commands executed
var CommandsLog = source.LOGGER.NewFileLogger(filepath.Join(source.ASSETS, "logs", "commands.log"), int64(source.OPTIONS.Ints("branding", "recycle_log")))

// Command represents the config for each individual command
type Command struct {
	parent      *Command
	Aliases     []string
	Description string
	Permissions []string
	Subcommands []*Command
	Args        []*Arg
	CommandFunc
	Callback
	Keypress
}

// CommandFunc is called upon whenever the command is executed
type CommandFunc func(*ArgContext, *sessions.Session) error

type Callback func(*ArgContext, *sessions.Session, int) []string

type Keypress func([]byte, *terminal.Read, *sessions.Session) ([]byte, bool)

// NewCommand will register a new command into the subcommand for the root
func (root *Command) NewCommand(cmd *Command) *Command {
	if cmd.Subcommands == nil {
		cmd.Subcommands = make([]*Command, 0)
	}

	middlewareBreak := cmd.CommandFunc

	// middleware for logging on command execution
	cmd.CommandFunc = func(ac *ArgContext, s *sessions.Session) error {
		CommandsLog.WriteLog("\"%s\" executed by \"%s\"", strings.Join(ac.Text, " "), s.User.Username)
		source.LOGGER.AggregateTerminal().WriteLog(gologr.DEBUG, "[COMMAND-EXEC] \"%s\" executed by \"%s\"", strings.Join(ac.Text, " "), s.User.Username)
		return middlewareBreak(ac, s)
	}

	cmd.parent = root
	root.Subcommands = append(root.Subcommands, cmd)
	return cmd
}

// indexSubcommands indexes and categorizes subcommands
func (root *Command) indexSubcommands(alias string) *Command {
	if root == nil { return nil }
	for _, cmd := range root.Subcommands {
		if !slices.Contains(cmd.Aliases, strings.Split(alias, "=")[0]) {
			continue
		}

		return cmd
	}

	return nil
}

// IndexCommands indexes and categorizes commands & subcommands
func (root *Command) IndexCommands(args []string, s *sessions.Session) *Command {
	if len(args) == 0 {
		if root == nil || len(s.Theme.CustomCommands) <= 0 {
			return ROOT
		}

		command := &Command{
			Aliases: root.Aliases,
			Permissions: root.Permissions,
			Description: root.Description,
			Subcommands: root.Subcommands,
		}

		for _, custom := range s.Theme.CustomCommands {
			switch typeCommand := custom.(type) {

			case *source.Text:
				command.Subcommands = append(command.Subcommands, &Command{
					Aliases: strings.Split(typeCommand.Name, ","),
					Permissions: strings.Split(typeCommand.Permissions, ","),
					Description: typeCommand.Description,
				})

				// whenever the permissions are nil
				if len(strings.Split(typeCommand.Permissions, ",")) == 1 && len(strings.Split(typeCommand.Permissions, ",")[0]) == 0 {
					command.Subcommands[len(command.Subcommands) - 1].Permissions = make([]string, 0)
				}
			}
		}


		return command
	}
	
	for {
		index := root.indexSubcommands(args[0])
		if index == nil && index != ROOT {
			return root
		}

		args = args[1:]
		if len(args) < 1 && index != ROOT {
			return index
		}

		root = index
	}
}

// Parents will find all the parents of a command
func (root *Command) Parents() []*Command {
	appendCommands := make([]*Command, 0)

	for {
		if root == nil || root.parent == nil {
			return appendCommands
		}

		appendCommands = append(appendCommands, root.parent)
		root = appendCommands[len(appendCommands) - 1]
	}
}

// Keypress will indirectly guide the handle direction of each command to a subcontext
func (node *Command) CallbackKeypress(session *sessions.Session) func([]byte, byte, *terminal.Read) ([]byte, bool) {
	return func(content []byte, key byte, state *terminal.Read) ([]byte, bool) {
		command := node.IndexCommands(strings.Split(string(content), " "), session)
		if command == nil || command.Keypress == nil {
			return nil, false
		}

		return command.Keypress(content, state, session)
	}
}

package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/sessions"
	"strings"

	"golang.org/x/exp/slices"
)

// Execute will execute a command through the given session
func (node *Command) Execute(session *sessions.Session, args []string) error {
	if len(args) == 0 || len(args[0]) == 0 {
		return nil
	}

	cmd := node.IndexCommands(args, session)
	if cmd == nil || strings.ToLower(args[0]) != "credits" && slices.Contains(session.Theme.DisabledCommands, strings.ToLower(strings.Join(args[:len(cmd.Parents())], "."))) {
		return CommandNotFound(session, args, cmd)
	}

	parents := cmd.Parents()
	for _, parent := range append(parents, cmd) {
		if functions.CanAccessThemPermissions(session.User, parent.Permissions...) {
			continue
		}

		return session.ExecuteBranding(make(map[string]any), "command_access_denied.tfx")
	}

	ctx, err := cmd.ParseArgs(args[len(parents):], args, session, true)
	if err != nil {
		return session.ExecuteBranding(map[string]any{"command": cmd.Aliases[0], "reason": err.Error()}, "invalid_argument.tfx")
	}

	if len(parents)-1 >= 1 {
		ctx.Header = strings.Join(strings.Split(args[len(parents)-1], "=")[1:], "=")
	}

	if cmd.CommandFunc == nil {
		return nil
	}

	err = cmd.CommandFunc(ctx, session)
	if err == nil || err == ErrCommandNotFound {
		if err == ErrCommandNotFound {
			return CommandNotFound(session, args, cmd)
		}

		return nil
	}

	if (len(parents) - 1) <= 0 {
		return session.ExecuteBranding(map[string]any{"command": cmd.Aliases[0], "err": err.Error()}, "command_error.tfx")
	} else {
		return session.ExecuteBranding(map[string]any{"command": cmd.Aliases[0], "err": err.Error()}, "subcommand_error.tfx")
	}
}

// CommandNotFound will directly and imprecisely check if there are any commands under this alias
func CommandNotFound(session *sessions.Session, args []string, cmd *Command) error {
	index := session.Theme.GetCustomCommand(args[0])
	if index == nil {
		return session.ExecuteBranding(map[string]any{"command": cmd.Aliases[0]}, "command_not_found.tfx")
	}

	switch concurrent := index.(type) {

	default:
		return session.ExecuteBranding(map[string]any{"command": cmd.Aliases[0]}, "command_not_found.tfx")

	case *source.Text:
		return session.Terminal.ExecuteStringToWriter(session.AppendDefaultSession(map[string]any{"clear": session.Terminal.ClearString}), concurrent.Tokenizer)
	}
}

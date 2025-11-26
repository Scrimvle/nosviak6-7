package commands

import (
	"Nosviak4/modules/gotable2"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"strings"
)

var Commands = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"commands", "cmds", "help"},
	Permissions: make([]string, 0),
	Description: "list all commands",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		tablet := gotable2.NewGoTable(&gotable2.Style{BorderValues: 1})
		tablet.Head(&gotable2.Row{
			Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "commands", "aliases.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "commands", "description.tfx"),
					Align: gotable2.AlignCenter,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "commands", "roles.tfx"),
					Align: gotable2.AlignCenter,
				},
			},
		})

		for _, command := range commands.ROOT.IndexCommands(commands.ToString(context.Args[0].Values), session).Subcommands {
			if !functions.CanAccessThemPermissions(session.User, command.Permissions...) {
				continue
			}

			// enable "fill" to support that if the array is empty, everyone can access it
			roles, err := interactions.PopulateStringWithRoles(session.Terminal, command.Permissions...)
			if err != nil {
				return err
			}

			tablet.Append(&gotable2.Row{Columns: []*gotable2.Column{
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"aliases": strings.Join(command.Aliases, ",")}, "commands", "commands", "value_aliases.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"description": command.Description}, "commands", "commands", "value_description.tfx"),
					Align: gotable2.AlignLeft,
				},
				{
					Text:  session.ExecuteBrandingToStringNoErr(map[string]any{"roles": roles}, "commands", "commands", "value_roles.tfx"),
					Align: gotable2.AlignLeft,
				},
			}})
		}

		return session.Table(tablet, context.Command.Aliases[0])
	},

	Args: []*commands.Arg{{
		Name:        "args",
		Type:        commands.STRING,
		OpenEnded:   true,
		Description: "filter commands",
		NotProvided: func(s *sessions.Session, _ []string) (string, error) {
			return "", nil
		},

		Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
			args := make([]string, 1)
			if len(ac.Args[0].Values) > 0 {
				args = append(make([]string, 0), commands.ToString(ac.Args[0].Values)...)
			}

			cmd := commands.ROOT.IndexCommands(args, s)
			if cmd == nil {
				return make([]string, 0)
			}

			ctx, _ := cmd.ParseArgs(args[len(cmd.Parents()):], args, s, false)
			if ctx == nil {
				return make([]string, 0)
			}

			return cmd.IndexPrefixReturnsBundle(args[len(args)-1], ctx, s, len(ctx.Tokens))
		},
	}},
})

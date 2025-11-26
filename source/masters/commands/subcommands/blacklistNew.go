package subcommands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/attacks"
	cmd "Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/commands/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"strconv"
	"time"
)

// BlacklistNew will create a brand new blacklisted target which expires after x amount of time
var BlacklistNew = commands.Blacklists.NewCommand(&cmd.Command{
	Aliases: []string{"new"},
	Description: "add to the blacklist",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(ac *cmd.ArgContext, s *sessions.Session) error {
		target := attacks.NewTarget(ac.Args[0].Values[0].ToString(), &source.Method{IPAllowed: true, URLAllowed: true})
		if !target.Validate() {
			return s.ExecuteBranding(map[string]any{"target": ac.Args[0].Values[0].ToString()}, "commands", "blacklist", "invalid_target.tfx")
		}

		// conv is a unit presented in days
		conv, err := strconv.Atoi(ac.Args[1].Values[0].ToString())
		if err != nil || conv == 0 {
			return s.ExecuteBranding(map[string]any{"target": ac.Args[0].Values[0].ToString()}, "commands", "blacklist", "invalid_int.tfx")
		}

		// can't have a expiry of blacklist which is larger than your plan
		if s.User.Expiry < int64(conv * 86400) {
			return s.ExecuteBranding(map[string]any{"target": ac.Args[0].Values[0].ToString()}, "commands", "blacklist", "expires_after_yourself.tfx")
		}

		// Performs the query to insert it into the database
		if err := database.DB.NewBlacklist(&database.Blacklist{User: s.User.ID, Target: ac.Args[0].Values[0].ToString(), Created: time.Now().Unix(), Expires: int64(conv * 86400)}); err != nil {
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "blacklist", "error_occurred.tfx")
		}

		return s.ExecuteBranding(map[string]any{"target": ac.Args[0].Values[0].ToString()}, "commands", "blacklist", "target_blacklisted.tfx")
	},

	Args: []*cmd.Arg{
		{
			Name: "target",
			Type: cmd.STRING,
			OpenEnded: false,
			Callback: func(ac *cmd.ArgContext, s *sessions.Session, i int) []string {
				targetBuf := make([]string, 0)
				targets, err := database.DB.GetUserAttacks(s.User.Username)
				if err != nil {
					return targetBuf
				}

				for _, attack := range targets {
					targetBuf = append(targetBuf, attack.Target)
				}

				return targetBuf
			},
		},

		{
			Name: "duration",
			Type: cmd.NUMBER,
			OpenEnded: false,
		},
	},
})
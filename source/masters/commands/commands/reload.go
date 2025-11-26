package commands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
)

var Reload = commands.ROOT.NewCommand(&commands.Command{
	Aliases: []string{"reload"},
	Description: "reloads the assets",
	Permissions: []string{interactions.ADMIN},
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		if err := source.OpenOptions(); err != nil {
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "reload", "reload_failed.tfx")
		}

		if err := sessions.PushConcurrentChangesAcrossSessions(s); err != nil {
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "reload", "reload_failed.tfx")
		}

		prompt, err := s.ExecuteBrandingToString(make(map[string]any), "prompt.tfx")
		if err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error occurred when %s tried to reload: %v", s.User.Username, err)
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "reload", "reload_failed.tfx")
		}

		s.Reader.Prompt = prompt
		return s.ExecuteBranding(make(map[string]any), "commands", "reload", "reload_success.tfx")
	},
})
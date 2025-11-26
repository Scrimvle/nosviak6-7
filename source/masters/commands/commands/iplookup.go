package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/functions/iplookup"
	"Nosviak4/source/masters/attacks"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var IpLookup = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"iplookup"},
	Description: "retrieves IP address information",
	Permissions: make([]string, 0),
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		if len(ac.Args) != 1 {
			return s.ExecuteBranding(make(map[string]any), "commands", "iplookup", "invalid_syntax.tfx")
		}

		target := ac.Args[0].Values[0].ToString()
		if !attacks.NewTarget(target, &source.Method{IPAllowed: true, URLAllowed: true}).Validate() {
			return s.ExecuteBranding(map[string]any{"ip": target}, "commands", "iplookup", "invalid_ip.tfx")
		}

		ip, err := iplookup.Lookup(target)
		if err != nil {
			return nil
		}

		return s.ExecuteBranding(map[string]any{"result": ip}, "commands", "iplookup", "result.tfx")
	},

	Args: []*commands.Arg{{
		Type:        commands.STRING,
		Name:        "ip",
		Description: "target for iplookup",
	}},
})

package subcommands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal/interactions"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

var PlansTokensCreateBatch = PlansTokens.NewCommand(&commands.Command{
	Aliases: []string{"createbatch"},
	Description: "create x amount of tokens",
	Permissions: []string{interactions.ADMIN, interactions.MOD},
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		plan, ok := source.Presets[ac.Args[1].Values[0].ToString()]
		if !ok || plan == nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "plans", "invalid_plan.tfx")
		}

		/* error converting to integer, caught a command handler error  */
		duration, err := strconv.Atoi(ac.Args[2].Values[0].ToString())
		if err != nil {
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "plans", "error_occurred.tfx")
		}

		creating, err := strconv.Atoi(ac.Args[0].Values[0].ToString())
		if err != nil {
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "plans", "error_occurred.tfx")
		}

		tokens := make([]string, 0)

		// loops through how many tokens are being created
		for i := 0; i < creating; i++ {
			octets := make([]string, 4)
			for index := range octets {
				octets[index] = hex.EncodeToString(*database.NewSalt(3))
			}
	
			token := &database.Token{
				Token: strings.Join(octets, "-"),
				Plan: ac.Args[1].Values[0].ToString(),
				Created: time.Now().Unix(),
				Expiry: int64(duration * 86400),
			}

			if err := database.DB.NewToken(token, s.User); err != nil {
				return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "plans", "error_occurred.tfx")
			}

			tokens = append(tokens, token.Token)
		}

		return s.Page(tokens)
	},

	Args: []*commands.Arg{
		{
			Name: "amount",
			Type: commands.NUMBER,
			Description: "amount of tokens to be created",
		},
		{
			Name: "plan",
			Type: commands.STRING,
			Description: "the plan to be given to the token",
			Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
				return maps.Keys(source.Presets)
			},

			NotProvided: func(s *sessions.Session, args []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "plan.tfx"), s.Reader.Context)
				content, err := read.ReadLine()
				if err != nil {
					return "", err
				}
	
				if len(content) == 0 {
					return "", fmt.Errorf("not allowed")
				}
	
				return string(content), nil
			},
		},
		{
			Name: "duration",
			Type: commands.NUMBER,
			Description: "the duration of the tokens lifetime",
			Callback: func(ac *commands.ArgContext, s *sessions.Session, i int) []string {
				return maps.Keys(source.Presets)
			},

			NotProvided: func(s *sessions.Session, args []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "plans", "tokens_length.tfx"), s.Reader.Context)
				content, err := read.ReadLine()
				if err != nil {
					return "", err
				}
	
				if len(content) == 0 {
					return "", fmt.Errorf("not allowed")
				}
	
				return string(content), nil
			},
		},
	},
})
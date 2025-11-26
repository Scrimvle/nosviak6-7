package commands

import (
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Connect allow users to connect there Nosviak4 account to a telegram account
var Connect = commands.ROOT.NewCommand(&commands.Command{
	Aliases: []string{"connect"},
	Description: "link your account with telegram",
	Permissions: make([]string, 0),
	CommandFunc: func(ac *commands.ArgContext, s *sessions.Session) error {
		commands.Mutex.Lock()
		defer commands.Mutex.Unlock()

		tracker, ok := commands.Trackers[ac.Args[0].Values[0].Literal.(int)]
		if !ok || tracker == nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "connect", "invalid_code.tfx")
		}

		message, err := s.ExecuteBrandingToString(make(map[string]any), "telegram", "linked.tfx")
		if err != nil {
			return err
		}
		
		s.User.Telegram = int(tracker.User)
		if err := database.DB.EditUser(s.User, s.User, commands.Conn.SendWebhook); err != nil {
			return s.ExecuteBranding(map[string]any{"err": err.Error()}, "commands", "connect", "error_occurred.tfx")
		}

		/* tries to alert the original owner of the code */
		if _, err := tracker.Bot.SendMessage(tracker.ID, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"}); err != nil {
			return s.ExecuteBranding(make(map[string]any), "commands", "connect", "error_linking.tfx")
		}

		return s.ExecuteBranding(make(map[string]any), "commands", "connect", "linked.tfx")
	},

	Args: []*commands.Arg{
		{
			Name: "code",
			Type: commands.NUMBER,
			Description: "code provided on the telegram bot",
			NotProvided: func(s *sessions.Session, args []string) (string, error) {
				read := s.Terminal.NewReadWithContext(s.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "connect", "code.tfx"), s.Reader.Context)
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
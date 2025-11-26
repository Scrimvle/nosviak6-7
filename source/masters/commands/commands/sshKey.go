package commands

import (
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
)

var SshKey = commands.ROOT.NewCommand(&commands.Command{
	Aliases:     []string{"sshkey"},
	Permissions: make([]string, 0),
	Description: "register for ssh key authentication",
	CommandFunc: func(context *commands.ArgContext, session *sessions.Session) error {
		if err := session.ExecuteBranding(make(map[string]any), "commands", "sshkey", "banner.tfx"); err != nil {
			return err
		}

		key, err := session.Terminal.NewReadWithContext(session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "sshkey", "prompt.tfx"), session.Reader.Context).ReadLine()
		if err != nil {
			return err
		}
		

		session.User.SSHKey = key
		if err := database.DB.EditUser(session.User, session.User, commands.Conn.SendWebhook); err != nil {
			return session.ExecuteBranding(make(map[string]any), "commands", "sshkey", "error.tfx")
		}

		return session.ExecuteBranding(make(map[string]any), "commands", "sshkey", "success.tfx")
	},
})

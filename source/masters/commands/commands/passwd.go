package commands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"bytes"
)

var Passwd = commands.ROOT.NewCommand(&commands.Command{
	Aliases: []string{"passwd"},
	Description: "change your password",
	Permissions: make([]string, 0),
	CommandFunc: func(ac *commands.ArgContext, session *sessions.Session) error {
		password, err := session.Terminal.NewReadWithContext(session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "passwd", "password.tfx"), session.Reader.Context).Mask([]byte(source.OPTIONS.String("password_mask"))).ChangeMaxLen(source.OPTIONS.Ints("maximum_password_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_password_length")).ReadLine()
		if err != nil {
			return err
		}

		confirmPassword, err := session.Terminal.NewReadWithContext(session.ExecuteBrandingToStringNoErr(make(map[string]any), "commands", "passwd", "confirm_password.tfx"), session.Reader.Context).Mask([]byte(source.OPTIONS.String("password_mask"))).ChangeMaxLen(source.OPTIONS.Ints("maximum_password_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_password_length")).ReadLine()
		if err != nil {
			return err
		}

		if !bytes.Equal(password, confirmPassword) {
			return session.ExecuteBranding(make(map[string]any), "commands", "passwd", "dont_equal.tfx")
		}

		session.User.Password = database.NewHash(confirmPassword, &session.User.Key)
		if err := database.DB.EditUser(session.User, session.User, commands.Conn.SendWebhook); err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Error occurred while changing password %s: %v", session.User.Username, err)
			return session.ExecuteBranding(make(map[string]any), "commands", "passwd", "error_occurred.tfx")
		}

		return session.ExecuteBranding(make(map[string]any), "commands", "passwd", "password_changed.tfx")
	},
})
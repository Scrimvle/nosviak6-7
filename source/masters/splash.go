package masters

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/sessions"
	"context"
	"math/rand"
	"strconv"
	"time"
)

// splashMiddleware happens inbetween the login and prompt being invoked, this prints home screens etc
func splashMiddleware(session *sessions.Session) error {

	/* handles the kick system */
	kicks, err := database.DB.GetOngoingKicks(session.User)
	if err != nil || kicks == nil || len(kicks) >= 1 {
		issuer, err := database.DB.GetUserWithID(kicks[len(kicks) - 1].Issuer)
		if err != nil {
			return err
		}

		if err := session.ExecuteBranding(map[string]any{"issuer": issuer.User(), "reason": kicks[len(kicks) - 1].Reason}, "session_user_kicked.tfx"); err != nil {
			return err
		}

		time.Sleep(4 * time.Second)
		return session.Terminal.Channel.Close()
	}
	
	/* handles the warns system */
	warns, err := database.DB.GetOngoingWarnings(session.User)
	if err != nil || warns == nil || source.OPTIONS.Ints("default_user", "warns_before_kick") < len(warns) && len(warns) > 0 {
		issuer, err := database.DB.GetUserWithID(warns[len(warns) - 1].Issuer)
		if err != nil {
			return err
		}

		if err := session.ExecuteBranding(map[string]any{"issuer": issuer.User(), "reason": warns[len(warns) - 1].Reason}, "session_warned_threshold.tfx"); err != nil {
			return err
		}

		time.Sleep(4 * time.Second)
		return session.Terminal.Channel.Close()
	}

	/* whenever the user is banned */
	if session.User.Banned {
		if err := session.ExecuteBranding(make(map[string]any), "login", "user_banned.tfx"); err != nil {
			return err
		}

		time.Sleep(4 * time.Second)
		return session.Terminal.Channel.Close()
	}

	/* Whenever the user is expired. */
	if session.User.Created + session.User.Expiry <= time.Now().Unix() {
		if err := session.ExecuteBranding(make(map[string]any), "login", "user_expired.tfx"); err != nil {
			return err
		}

		time.Sleep(4 * time.Second)
		return session.Terminal.Channel.Close()
	}

	if len(sessions.IndexSessions(session.User.Username)) - 1 > session.User.Sessions {
		if err := session.ExecuteBranding(make(map[string]any), "login", "max_sessions.tfx"); err != nil {
			return err
		}

		time.Sleep(4 * time.Second)
		return session.Terminal.Channel.Close()
	}

	// Forces them forwards for the captcha prompt
	if session.Terminal.Config.Captcha.Enabled && !functions.CanAccessThemPermissions(session.User, session.Terminal.Config.Captcha.IgnorePerms...) {
		x := rand.Intn(session.Terminal.Config.Captcha.Max - session.Terminal.Config.Captcha.Min) + session.Terminal.Config.Captcha.Min
		y := rand.Intn(session.Terminal.Config.Captcha.Max - session.Terminal.Config.Captcha.Min) + session.Terminal.Config.Captcha.Min
		if err := session.ExecuteBranding(map[string]any{"number1": x, "number2": y}, "login", "captcha", "banner.tfx"); err != nil {
			return err
		}

		session.Reader = session.Terminal.NewRead(session.ExecuteBrandingToStringNoErr(make(map[string]any), "login", "captcha", "prompt.tfx"))
		session.Reader.Context, session.Cancel = context.WithCancel(context.Background())
		text, err := session.Reader.ReadLine()
		if err != nil {
			return err
		}

		answer, err := strconv.Atoi(string(text))
		if err != nil || x + y != answer {
			return session.ExecuteBranding(make(map[string]any), "login", "captcha", "incorrect.tfx")
		}
	}

	/* newUser forces the user to change there password */
	if session.User.NewUser {
		err := NewUser(session)
		if err != nil {
			return err
		}
	}

	return prompt(session)
}
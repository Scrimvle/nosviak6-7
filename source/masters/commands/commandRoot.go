package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/attacks"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"Nosviak4/source/masters/terminal/interactions"
	"Nosviak4/source/swash"
	"Nosviak4/source/swash/packages"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// ROOT stores all the registered commands
var ROOT *Command = &Command{
	Aliases:     []string{"<root_node>"},
	Description: "root node for the command handler",
	Subcommands: make([]*Command, 0),

	// CommandFunc is called upon whenever they press enter
	CommandFunc: func(ac *ArgContext, session *sessions.Session) error {
		if !strings.HasPrefix(ac.Text[0], source.MethodConfig.Attacks.AttackPrefix) || len(ac.Text[0]) <= len(source.MethodConfig.Attacks.AttackPrefix) {
			return executeCustomCommands(ac, session)
		}

		method, ok := source.Methods[ac.Text[0][len(source.MethodConfig.Attacks.AttackPrefix):]]
		if !ok || method == nil {
			return session.ExecuteBranding(make(map[string]any), "attack_invalid_method.tfx")
		}

		if !functions.CanAccessThemPermissions(session.User, method.Access...) && !database.DB.IsSuperuser(session.User) {
			return session.ExecuteBranding(make(map[string]any), "attack_access_denied.tfx")
		}

		/* handles the parsing and validating of the duration param */
		duration, err := strconv.Atoi(ac.Args[2].Values[0].ToString())
		if err != nil && !source.OPTIONS.Bool("attacks", "port_then_duration") {
			return session.ExecuteBranding(make(map[string]any), "attack_bad_duration.tfx")
		}

		/* handles the parsing of the port param */
		port, err := strconv.Atoi(ac.Args[3].Values[0].ToString())
		if err != nil && !source.OPTIONS.Bool("attacks", "port_then_duration") {
			return session.ExecuteBranding(make(map[string]any), "attack_bad_port.tfx")
		}

		/* handles if the pos changes */
		if source.OPTIONS.Bool("attacks", "port_then_duration") {
			duration, err = strconv.Atoi(ac.Args[3].Values[0].ToString())
			if err != nil {
				return session.ExecuteBranding(make(map[string]any), "attack_bad_duration.tfx")
			}

			port, err = strconv.Atoi(ac.Args[2].Values[0].ToString())
			if err != nil {
				return session.ExecuteBranding(make(map[string]any), "attack_bad_port.tfx")
			}
		}

		// changed: allowed for 0 port
		if port < 0 || port > 65536 {
			return session.ExecuteBranding(make(map[string]any), "attack_bad_port.tfx")
		}

		// checks the input length to ensure we don't get any unhandled error appearing
		if len(ac.Args) <= 1 || len(ac.Args[1].Values) == 0 {
			return session.ExecuteBranding(make(map[string]any), "attack_bad_target.tfx")
		}

		suggestedMethod, suggestedMethodName, err := attacks.NewTarget(ac.Args[1].Values[0].ToString(), method).Suggest()
		if err != nil {
			return session.ExecuteBranding(make(map[string]any), "attack_bad_target.tfx")
		}

		// checks if we should display the suggestion forced message
		if suggestion, ok := source.Suggestions[suggestedMethodName]; ok && suggestedMethod != nil && suggestion.Forced {
			return session.ExecuteBranding(map[string]any{"suggested_method": suggestedMethodName, "target": ac.Args[1].Values[0].ToString()}, "attack_suggestion_enforced.tfx")
		}

		kvs, err := functions.HandleKeyValues(ac.Text[1:], method)
		if err != nil {
			return session.ExecuteBranding(map[string]any{"err": err.Error()}, "attack_key_value_void.tfx")
		}

		// NewAttack receives information about the duration and port which have been verified.
		err = attacks.NewAttack(session, method, ac.Text[0][len(source.MethodConfig.Attacks.AttackPrefix):], ac.Args[1].Values[0].ToString(), duration, port, kvs, Conn.SendWebhook)
		switch err {

		case attacks.ErrUnauthorizedSender:
			return session.ExecuteBranding(make(map[string]any), "attack_access_denied.tfx")

		case attacks.ErrAttacksDisabled:
			return session.ExecuteBranding(make(map[string]any), "attack_disabled.tfx")

		case attacks.ErrInvalidTarget:
			return session.ExecuteBranding(make(map[string]any), "attack_bad_target.tfx")

		case attacks.ErrMaxtimeGreaterThanEnhanced:
			return session.ExecuteBranding(map[string]any{"maxtime": method.Options.EnhancedMaxtime}, "attack_duration_greater_method.tfx")

		case attacks.ErrMaxtimeGreaterThanOverride:
			return session.ExecuteBranding(map[string]any{"maxtime": method.Options.MaxtimeOverride}, "attack_duration_greater_override.tfx")

		case attacks.ErrMaxtimeGreaterThanSelf:
			return session.ExecuteBranding(make(map[string]any), "attack_duration_greater_self.tfx")

		case attacks.ErrMinimumMaxTime:
			return session.ExecuteBranding(make(map[string]any), "attack_minimum_duration.tfx")

		case attacks.ErrOngoingMethodCap:
			return session.ExecuteBranding(make(map[string]any), "attack_no_method_slots.tfx")

		case attacks.ErrConcurrentLimitReached:
			running, err := database.DB.GetOngoingUser(session.User)
			if err != nil {
				return err
			}

			var nextToFinish *database.Attack = running[0]

			for _, attack := range running[1:] {
				if attack.Created+attack.Duration <= nextToFinish.Created+nextToFinish.Duration {
					continue
				}

				nextToFinish = attack
			}

			return session.ExecuteBranding(map[string]any{"wait": packages.Until(nextToFinish.Created + nextToFinish.Duration)}, "attack_concurrent_limit_reached.tfx")

		case attacks.ErrGroupConcurrentLimitReached:
			running, err := database.DB.GetOngoingGroup(method.Options.MethodGroup)
			if err != nil {
				return err
			}

			var nextToFinish *database.Attack = running[0]

			for _, attack := range running[1:] {
				if attack.Created+attack.Duration <= nextToFinish.Created+nextToFinish.Duration {
					continue
				}

				nextToFinish = attack
			}

			return session.ExecuteBranding(map[string]any{"wait": packages.Until(nextToFinish.Created + nextToFinish.Duration)}, "attack_group_slots_limit_reached.tfx")

		// user cooldown
		case attacks.ErrCooldownPeriodEngaged:
			running, err := database.DB.GetOngoingUser(session.User)
			if err != nil {
				return err
			}

			var nextToFinish *database.Attack = running[0]

			for _, attack := range running[1:] {
				if attack.Created+attack.Duration <= nextToFinish.Created+nextToFinish.Duration {
					continue
				}

				nextToFinish = attack
			}

			return session.ExecuteBranding(map[string]any{"wait": packages.Until(nextToFinish.Created + session.User.Cooldown)}, "attack_wait_cooldown.tfx")

		case attacks.ErrUserBlacklisted:
			target, err := database.DB.GetBlacklistedTarget(ac.Args[1].Values[0].ToString())
			if err != nil || target == nil {
				return err
			}

			return session.ExecuteBranding(map[string]any{"target": target.Target, "expires": int(target.Created + target.Expires)}, "attack_user_blacklisted_target.tfx")

		case attacks.ErrSystemBlacklisted:
			return session.ExecuteBranding(map[string]any{"target": ac.Args[1].Values[0].ToString()}, "attack_blacklisted_target.tfx")

		case attacks.ErrGlobalConns:
			return session.ExecuteBranding(make(map[string]any), "attack_no_global_slots.tfx")

		case attacks.ErrGlobalMaxtime:
			return session.ExecuteBranding(make(map[string]any), "attack_duration_greater_global.tfx")

		case attacks.ErrMaxAttacksToday:
			return session.ExecuteBranding(make(map[string]any), "attack_max_daily_attacks.tfx")

		// global cooldown
		case attacks.ErrGlobalCooldown:
			running, err := database.DB.GetOngoing()
			if err != nil {
				return err
			}

			var nextToFinish *database.Attack = running[0]
			for _, attack := range running[0:] {
				if attack.Created+attack.Duration <= nextToFinish.Created+nextToFinish.Duration {
					continue
				}

				nextToFinish = attack
			}

			return session.ExecuteBranding(map[string]any{"wait": packages.Until(nextToFinish.Created + session.User.Cooldown)}, "attack_global_cooldown_enabled.tfx")

		case attacks.ErrMethodDisabled:
			return session.ExecuteBranding(make(map[string]any), "attack_method_disabled.tfx")

		case attacks.ErrPowersaving:
			running, err := database.DB.GetOngoingTarget(ac.Args[1].Values[0].ToString())
			if err != nil || len(running) == 0 {
				return session.ExecuteBranding(make(map[string]any), "attack_failed.tfx")
			}

			var nextToFinish *database.Attack = running[0]
			for _, attack := range running[1:] {
				if attack.Created+attack.Duration <= nextToFinish.Created+nextToFinish.Duration {
					continue
				}

				nextToFinish = attack
			}

			return session.ExecuteBranding(map[string]any{"method": running[0].Method, "wait": packages.Until(nextToFinish.Created + nextToFinish.Duration)}, "attack_powersaving_enabled.tfx")

		// Group cooldown
		case attacks.ErrGroupCooldown:
			running, err := database.DB.GetOngoingGroup(method.Options.MethodGroup)
			if err != nil {
				return err
			}

			var nextToFinish *database.Attack = running[0]
			for _, attack := range running[0:] {
				if attack.Created+attack.Duration <= nextToFinish.Created+nextToFinish.Duration {
					continue
				}

				nextToFinish = attack
			}

			return session.ExecuteBranding(map[string]any{"wait": packages.Until(nextToFinish.Created + session.User.Cooldown)}, "attack_group_cooldown.tfx")

		// Attack launched towards the target successfully
		case nil:
			props := make(map[string]any)
			props["target"] = ac.Args[1].Values[0].ToString()
			props["port"] = port
			props["duration"] = duration
			props["method"] = ac.Text[0][len(source.MethodConfig.Attacks.AttackPrefix):]
			props["created"] = ac.init.Unix()

			// kv("len")
			props["kv"] = func(key string) string {
				return fmt.Sprint(kvs[key])
			}

			if len(method.Options.AttackSend) >= 1 {
				return session.ExecuteBranding(make(map[string]any), method.Options.AttackSend)
			}

			return session.ExecuteBranding(props, "attack_sent.tfx")

		default:
			if slices.Contains(session.User.Roles, interactions.ADMIN) && database.DB.IsSuperuser(session.User) {
				urls := make([]string, 0)
				for _, url := range method.URLs {
					tokenizer := swash.NewTokenizer(string(url), true).Strip()
					if err := tokenizer.Parse(); err != nil {
						return session.ExecuteBranding(map[string]any{"err": err.Error()}, "attack_unknown_error.tfx")
					}

					params := map[string]any{
						"username": session.User.Username,
						"method":   ac.Text[0][len(source.MethodConfig.Attacks.AttackPrefix):],
						"target":   ac.Args[1].Values[0].ToString(),
						"time":     strconv.Itoa(duration),
						"port":     strconv.Itoa(port),
						"kv": func(key string) string {
							value, ok := kvs[key]
							if !ok || value == "<nil>" {
								return ""
							}

							return fmt.Sprint(value)
						},
					}

					url := bytes.NewBuffer(make([]byte, 0))
					if err := terminal.ExecuteStringToWriter(url, params, tokenizer); err != nil {
						return session.ExecuteBranding(map[string]any{"err": err.Error()}, "attack_unknown_error.tfx")
					}

					urls = append(urls, url.String())
				}

				// builds the error payload
				errBuf := bytes.NewBuffer(make([]byte, 0))
				fmt.Fprintf(errBuf, "[ATTACK FAILED TO SEND DUE TO REMOTE ERROR]\r\n")
				fmt.Fprintf(errBuf, "[Attack] Method: \"%s\" Port: %d Duration: %d Target: \"%s\"\r\n", ac.Text[0][len(source.MethodConfig.Attacks.AttackPrefix):], port, duration, ac.Args[1].Values[0].ToString())
				fmt.Fprintf(errBuf, "[API URL] %v\r\n", strings.Join(urls, ","))
				fmt.Fprintf(errBuf, "[ERROR] %v\r\n\r\n", err)
				fmt.Fprintf(errBuf, "What now? The error was not an issue with Nosviak. Either you have made a mistake in your api.json file, or your API provider has made a mistake. DO NOT CONTACT FB OR HARRY; please reach out to our support team if you have any queries.")

				return session.ExecuteBranding(map[string]any{"err": errBuf.String()}, "attack_failed_superuser.tfx")
			}

			return session.ExecuteBranding(make(map[string]any), "attack_failed.tfx")
		}
	},

	// Keypress implements concurrent live alerts on the attacks
	Keypress: func(content []byte, reader *terminal.Read, session *sessions.Session) ([]byte, bool) {
		args := strings.Split(string(content), " ")
		if len(args) == 0 || !strings.HasPrefix(args[0], source.MethodConfig.Attacks.AttackPrefix) || !session.Terminal.XTerm {
			return nil, false
		}

		// checks if the method exists or not, this will render alerts etc.
		method, ok := source.Methods[args[0][len(source.MethodConfig.Attacks.AttackPrefix):]]
		if !ok || method == nil {
			reader.PostAlert(&terminal.Alert{
				AlertCode:    terminal.MESSAGE,
				AlertMessage: session.ExecuteBrandingToStringNoErr(make(map[string]any), "alerts", "attack_method_not_found.tfx"),
			})

			/*
				continues to loop reading in trying to see if the user corrected themselves, if they don't the alert
				persists but if the remove everything the alert disappears.
			*/

			for {
				// reads in from the session using a signal worker.
				buf, err := session.Terminal.Signal.ReadWithContext(session.Reader.Context)
				if err != nil {
					return nil, true
				}

				// handles through the buf so we can access the content again.
				ok, err := reader.Buf(buf, true)
				if err != nil || ok {
					reader.DisgardAlert()
					if _, err := session.Terminal.Write([]byte("\r\n")); err != nil {
						return nil, true
					}

					return reader.Content(), true
				}

				// whenever the prefix is removed.
				args = strings.Split(string(reader.Content()), " ")
				if len(args) == 0 || !strings.HasPrefix(args[0], source.MethodConfig.Attacks.AttackPrefix) {
					reader.DisgardAlert()
					return nil, false
				}

				// if method found, we disconnect from the for loop and return everything
				method, ok = source.Methods[args[0][len(source.MethodConfig.Attacks.AttackPrefix):]]
				if ok && method != nil {
					reader.DisgardAlert()
					break
				}

				// keep looping as they haven't corrected themselves
			}

		}

		// missing arguments to continue verifying
		if len(args) <= 1 {
			return nil, false
		}

		// tries to verify the target directly
		target := attacks.NewTarget(args[1], method)
		if !target.Validate() {
			reader.PostAlert(&terminal.Alert{
				AlertCode:    terminal.MESSAGE,
				AlertMessage: session.ExecuteBrandingToStringNoErr(make(map[string]any), "alerts", "attack_invalid_target.tfx"),
			})

			/*
				continues to loop reading in trying to see if the user corrected themselves, if they don't the alert
				persists but if the remove everything the alert disappears.
			*/

			for {
				// reads in from the session using a signal worker.
				buf, err := session.Terminal.Signal.ReadWithContext(session.Reader.Context)
				if err != nil {
					return nil, true
				}

				// handles through the buf so we can access the content again.
				ok, err := reader.Buf(buf, true)
				if err != nil || ok {
					reader.DisgardAlert()
					if _, err := session.Terminal.Write([]byte("\r\n")); err != nil {
						return nil, true
					}

					return reader.Content(), true
				}

				args = strings.Split(string(reader.Content()), " ")
				if len(args) <= 1 || len(args[1]) == 0 {
					reader.DisgardAlert()
					return nil, false
				}

				// checks if the target is not valid
				if attacks.NewTarget(args[1], method).Validate() {
					reader.DisgardAlert()
					break
				}
			}
		}

		// continues to check if the target is now valid or not.
		if target.Blacklisted() && !functions.CanAccessThemPermissions(session.User, source.OPTIONS.Strings("bypass_blacklist")...) {
			reader.PostAlert(&terminal.Alert{
				AlertCode:    terminal.MESSAGE,
				AlertMessage: session.ExecuteBrandingToStringNoErr(make(map[string]any), "alerts", "attack_target_blacklisted.tfx"),
			})

			/*
				continues to loop reading in trying to see if the user corrected themselves, if they don't the alert
				persists but if the remove everything the alert disappears.
			*/

			for {
				// reads in from the session using a signal worker.
				buf, err := session.Terminal.Signal.ReadWithContext(session.Reader.Context)
				if err != nil {
					return nil, true
				}

				// handles through the buf so we can access the content again.
				ok, err := reader.Buf(buf, true)
				if err != nil || ok {
					reader.DisgardAlert()
					if _, err := session.Terminal.Write([]byte("\r\n")); err != nil {
						return nil, true
					}

					return reader.Content(), true
				}

				args = strings.Split(string(reader.Content()), " ")
				if len(args) <= 1 || len(args[1]) == 0 {
					reader.DisgardAlert()
					return nil, false
				}

				// target is no longer blacklisted.
				if !attacks.NewTarget(args[1], method).Blacklisted() {
					reader.DisgardAlert()
					return nil, false
				}
			}
		}

		// implements method suggestions
		methodSuggested, methodName, err := target.Suggest()
		if err != nil || method == nil || methodSuggested == nil || methodSuggested == method {
			return nil, false
		}

		reader.PostAlert(&terminal.Alert{
			AlertCode:    terminal.MESSAGE,
			AlertMessage: session.ExecuteBrandingToStringNoErr(map[string]any{"method": args[0][len(source.MethodConfig.Attacks.AttackPrefix):], "method_recommended": methodName, "target": args[1]}, "alerts", "attack_suggestion.tfx"),
		})

		for {
			// reads in from the session using a signal worker.
			buf, err := session.Terminal.Signal.ReadWithContext(session.Reader.Context)
			if err != nil {
				return nil, true
			}

			// handles through the buf so we can access the content again.
			ok, err := reader.Buf(buf, true)
			if err != nil || ok {
				reader.DisgardAlert()
				if _, err := session.Terminal.Write([]byte("\r\n")); err != nil {
					return nil, true
				}

				return reader.Content(), true
			}

			args = strings.Split(string(reader.Content()), " ")
			if len(args) <= 1 || len(args[1]) == 0 {
				reader.DisgardAlert()
				return nil, false
			}

			method, methodName, err := target.Suggest()
			if err != nil || method == nil || methodSuggested == nil || len(methodName) == 0 {
				return nil, reader.DisgardAlert() != nil
			}
		}
	},
}

// executeCustomCommands will execute any custom command types
func executeCustomCommands(ac *ArgContext, s *sessions.Session) error {
	for _, command := range s.Theme.CustomCommands {
		switch index := command.(type) {

		case *source.Text:
			if !strings.EqualFold(ac.Text[0], index.Name) {
				continue
			}

			if !functions.CanAccessThemPermissions(s.User, index.Permissions) {
				return s.ExecuteBranding(make(map[string]any), "command_access_denied.tfx")
			}

			return s.Terminal.ExecuteStringToWriter(s.AppendDefaultSession(make(map[string]any)), index.Tokenizer)

		case *source.Bin:
			if !slices.Contains(index.Name, strings.ToLower(ac.Text[0])) {
				continue
			}

			if !functions.CanAccessThemPermissions(s.User, index.Permissions...) {
				return s.ExecuteBranding(make(map[string]any), "command_access_denied.tfx")
			}

			return executeBinCommand(index, s)
		}
	}

	return ErrCommandNotFound
}

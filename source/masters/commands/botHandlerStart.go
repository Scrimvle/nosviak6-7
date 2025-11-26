package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/terminal"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// StartCommandBot is executed when the poll results with `/start`
func StartCommandBot(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := database.DB.GetUserTelegram(int(ctx.EffectiveUser.Id))
	if err != nil || user == nil {
		return StartCommandLink(b, ctx)
	}

	message, err := terminal.ExecuteBrandingToString(map[string]any{"user": user.User()}, source.ASSETS, source.BRANDING, "telegram", "start.tfx")
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, message, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	return err
}

// StartCommandLink is when the `/start` command is executed without being linked
func StartCommandLink(b *gotgbot.Bot, ctx *ext.Context) error {
	executedMessage, err := terminal.ExecuteBrandingToString(make(map[string]any), source.ASSETS, source.BRANDING, "telegram", "not_linked.tfx")
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, executedMessage, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	return err
}
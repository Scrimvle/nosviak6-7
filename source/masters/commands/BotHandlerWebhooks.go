package commands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/functions"
	"Nosviak4/source/masters/terminal"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"golang.org/x/exp/slices"
)

// SendWebhook will send the webhook to whatever endpoints required
func (b *Bot) SendWebhook(event string, vars map[string]any) {
	go func() {
		users, err := database.DB.GetUsers()
		if err != nil {
			return
		}

		roles, triggers := source.OPTIONS.Strings("web", "telegram", "logs", "log_roles"), source.OPTIONS.Strings("web", "telegram", "logs", "triggers")
		if len(triggers) == 0 || !slices.Contains(triggers, event) {
			return
		}

		content, err := terminal.ExecuteBrandingToString(vars, source.ASSETS, source.BRANDING, "telegram", "webhooks", fmt.Sprintf("%s.tfx", event))
		if err != nil {
			source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "[WEBHOOK] issue occurred with `%s` event error: %v", err.Error())
			return
		}

		for _, user := range users {
			if !functions.CanAccessThemPermissions(user, roles...) || user.Telegram == 0 {
				continue
			}

			// attempts to send
			b.conn.SendMessage(int64(user.Telegram), content, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
		}
	}()
}
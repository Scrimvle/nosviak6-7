package commands

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/source"
	"path/filepath"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// Bot is a direct implementation of the bot interface
type Bot struct {
	conn *gotgbot.Bot
	init time.Time
}

var (
	// conn is the concurrent connection to the bot
	Conn *Bot = new(Bot)
	mutex sync.Mutex
)

// OpenBotConn spawns a telegram session
func OpenBotConn() error {
	if !source.OPTIONS.Bool("web", "telegram", "enabled") {
		return nil
	}

	logger := source.LOGGER.NewFileLogger(filepath.Join(source.ASSETS, "logs", "telegram.log"), int64(source.OPTIONS.Ints("branding", "recycle_log")))
	if logger.Err != nil {
		return logger.Err
	}

	defer logger.Close()
	t := logger.WithTerminal()

	/* launches the new bot session with the package */
	session, err := gotgbot.NewBot(source.OPTIONS.String("web", "telegram", "token"), &gotgbot.BotOpts{DefaultRequestOpts: &gotgbot.RequestOpts{Timeout: gotgbot.DefaultTimeout, APIURL: gotgbot.DefaultAPIURL}})
	if err != nil {
		return err
	}

	t.WriteLog(gologr.ALERT, "- Telegram bot is now available: @%s", session.Username)

	mutex.Lock()
	Conn = &Bot{
		conn: session,
		init: time.Now(),
	}

	mutex.Unlock()
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			MaxRoutines: ext.DefaultMaxRoutines,

			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				t.WriteLog(gologr.ERROR, "Error occurred with the telegram API client: %v", err)
				return ext.DispatcherActionNoop
			},
		}),
	})

	// short list of all the registered commands
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", StartCommandBot))
	updater.Dispatcher.AddHandler(handlers.NewCommand("launch", LaunchCommandBot))
	updater.Dispatcher.AddHandler(handlers.NewCommand("connect", ConnectCommandBot))

	return updater.StartPolling(session, &ext.PollingOpts{DropPendingUpdates: true, GetUpdatesOpts: gotgbot.GetUpdatesOpts{}})
}
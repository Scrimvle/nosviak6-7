package commands

import (
	"Nosviak4/source"
	"Nosviak4/source/database"
	"Nosviak4/source/masters/terminal"
	"crypto/rand"
	"encoding/binary"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Tracker is when someone enters the screen, there can only be one Tracker per ID
type Tracker struct {
	Bot  *gotgbot.Bot
	ID   int64
	User int64
	Code int
}

var (
	Trackers map[int]*Tracker = make(map[int]*Tracker)
	Mutex    sync.Mutex
)

// ConnectCommandBot is when the user executes the connect command
func ConnectCommandBot(b *gotgbot.Bot, ctx *ext.Context) error {
	user, err := database.DB.GetUserTelegram(int(ctx.EffectiveUser.Id))
	if err == nil && user != nil {
		return StartCommandBot(b, ctx)
	}

	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return ConnectCommandBot(b, ctx)
	}

	tracker := &Tracker{
		Bot:  b,
		ID:   ctx.EffectiveSender.ChatId,
		User: ctx.EffectiveUser.Id,
		Code: int(binary.BigEndian.Uint16(buf)),
	}

	Mutex.Lock()
	defer Mutex.Unlock()
	Trackers[tracker.Code] = tracker
	
	// Expires the code after 10 mins
	time.AfterFunc(10 * time.Minute, func() {
		if _, ok := Trackers[tracker.Code]; !ok {
			return
		}

		Mutex.Lock()
		defer Mutex.Unlock()
		delete(Trackers, tracker.Code)

		/* tells the chat when the code expires */
		b.SendMessage(ctx.EffectiveSender.ChatId, "The code has expired.", &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	})

	executedMessage, err := terminal.ExecuteBrandingToString(map[string]any{"code": tracker.Code}, source.ASSETS, source.BRANDING, "telegram", "link_pending.tfx")
	if err != nil {
		return err
	}

	_, err = ctx.EffectiveMessage.Reply(b, executedMessage, &gotgbot.SendMessageOpts{ParseMode: "markdown"})
	return err
}
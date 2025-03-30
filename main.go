package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

func main() {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

// start introduces the bot.
func start(b *gotgbot.Bot, ctx *ext.Context) error {
 file := gotgbot.InputFileByID("config.StartImageID")

	b.SendPhoto(ctx.EffectiveChat.Id, file, nil)
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s.\nI am a sample bot", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

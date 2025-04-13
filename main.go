package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/database"
	"github.com/Vivekkumar-IN/EditguardianBot/filters"
	"github.com/Vivekkumar-IN/EditguardianBot/modules"
)

func main() {
	b, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	defer database.Disconnect()

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			msg := fmt.Sprintf(
				"❗ <b>Error occurred</b>\n\n<b>Error:</b> <code>%vscode>\n<b>From User:</b> %vsn<b>Chat:</b> %v\n<b>Data:</b> <pre>%v</pre>",
				err,
				ctx.EffectiveUser.FirstName,
				ctx.EffectiveChat.Title,
			)

			b.SendMessage(config.LoggerId, msg, &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			})

			return ext.DispatcherActionContinueGroups
		},
		MaxRoutines: 500,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	// Handlers

	for _, h := range modules.Handlers {
		dispatcher.AddHandler(h)
	}
	dispatcher.AddHandlerToGroup(handlers.NewMessage(
		filters.Invert(filters.ChatAdmins(b)),
		modules.DeleteEditedMessage,
	).SetAllowEdited(true), -1)

	dispatcher.AddHandler(handlers.NewMessage(filters.Invert(filters.ChatAdmins(b)), modules.DeleteLongMessage))

	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Command(b, "ping"), modules.PingHandle), 1)

	// Allowed updates
	allowedUpdates := []string{
		"message",
		"my_chat_member",
		"chat_member",
		"edited_message",
		"callback_query",
	}

	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout:        9,
			AllowedUpdates: allowedUpdates,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("Failed to start polling: " + err.Error())
	}

	log.Printf("%s has been started...\n", b.User.Username)
	_, _ = b.SendMessage(
		config.LoggerId,
		fmt.Sprintf("%s has started\n", b.User.Username),
		nil,
	)

	updater.Idle()
}

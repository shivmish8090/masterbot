package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"

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
		Error: func(_ *gotgbot.Bot, _ *ext.Context, err error) ext.DispatcherAction {
			log.Println("An error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: 500,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// Handlers

	for _, h := range modules.CommandHandlers {
		dispatcher.AddHandler(h)
	}
	for _, h := range modules.MessageHandlers {
		dispatcher.AddHandler(h)
	}
	for _, h := range modules.CallbackQueryHandlers {
		dispatcher.AddHandler(h)
	}

	dispatcher.AddHandler(handlers.NewMyChatMember(
		func(u *gotgbot.ChatMemberUpdated) bool {
			wasMember, isMember := ExtractJoinLeftStatusChange(u)
			return !wasMember && isMember
		},
		AddedToGroups,
	))

	evalHandler := handlers.NewMessage(
		filters.AndFilter(filters.Owner, filters.Command(b, "eval")),
		EvalHandler,
	).SetAllowEdited(true)

	dispatcher.AddHandler(evalHandler)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(
		filters.Invert(filters.ChatAdmins(b)),
		deleteEditedMessage,
	).SetAllowEdited(true), -1)

	dispatcher.AddHandler(handlers.NewMessage(
		filters.And(filters.LongMessage, filters.Invert(filters.ChatAdmins(b))),
		deleteLongMessage,
	))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("help_callback"), helpCB))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("help_echo"), echoCB))

	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Equal("close"), closeCallback))

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

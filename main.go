package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/filters"
)

func main() {
	b, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: 500,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	filters.Init(b)

	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewMyChatMember(
		func(u *gotgbot.ChatMemberUpdated) bool {
			wasMember, isMember := ExtractJoinLeftStatusChange(u)
			return !wasMember && isMember
		},
		AddedToGroups,
	))

	evalHandler := handlers.NewMessage(
		filters.AndFilter(filters.Owner, filters.Command("eval")),
		EvalHandler,
	).SetAllowEdited(true)

	dispatcher.AddHandler(evalHandler)
	dispatcher.AddHandler(handlers.NewCommand("echo", EcoHandler))
	dispatcher.AddHandlerToGroup(handlers.NewMessage(
		filters.Invert(filters.ChatAdmins),
		deleteEditedMessage,
	).SetAllowEdited(true), -1)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(
		filters.LongMessage,
		deleteLongMessage,
	), -1)

	allowedUpdates := []string{
		"message",
		"my_chat_member",
		"chat_member",
		"edited_message",
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
		panic("failed to start polling: " + err.Error())
	}

	log.Printf("%s has been started...\n", b.User.Username)
	b.SendMessage(
		config.LoggerId,
		fmt.Sprintf("%s has started\n", b.User.Username),
		nil,
	)
	updater.Idle()
}

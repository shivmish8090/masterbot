package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func init() {
	Register(handlers.NewCommand("reload", ReloadHandler))
}

func ReloadHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	ChatId := ctx.EffectiveChat.Id

	msg.Delete(b, nil)
	b.SendMessage(ChatId, "Refreshing Cache of chat admins...", nil)
	return nil
}

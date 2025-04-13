package modules

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/config/helpers"
)

func init() {
	Register(handlers.NewCommand("reload", ReloadHandler))
}

func ReloadHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chatId := ctx.EffectiveChat.Id

	msg.Delete(b, nil)

	x, err := b.SendMessage(chatId, "Refreshing Cache of chat admins...", nil)
	if err != nil {
		return err
	}

	chatMembers, e := b.GetChatAdministrators(chatId, nil)
	if e != nil {
		x.EditText(b, fmt.Sprintf("⚠️ Cache refresh failed — %v", e), nil)
		return e
	}

	var admins []int64
	for _, m := range chatMembers {
		status := m.GetStatus()
		if status == "administrator" || status == "creator" {
			admins = append(admins, m.GetUser().Id)
		}
	}

	config.Cache.Store(fmt.Sprintf("admins:%d", chatId), admins)

	var text string
	if helpers.Contains(admins, ctx.EffectiveUser.Id) {
		text = "✅ Successfully refreshed the cache of chat admins!"
	} else {
		text = "⚠️ Tried refreshing the admin cache... but You're not an admin!"
	}

	_, _, err := x.EditText(b, text, nil)

	return orCont(err)
}

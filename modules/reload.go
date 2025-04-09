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
	x, err := b.SendMessage(ChatId, "Refreshing Cache of chat admins...", nil)
if err != nil {

return err
}

                        chatmembers, e := b.GetChatAdministrators(ChatId, nil)
                        if e != nil {


x.EditTexg(b, fmt.Sprintf("Oops! Cache refresh failed — %v", e), nil)
                                return e
                        }

                        var admins []int64
                        for _, m := range chatmembers {
                                status := m.GetStatus()
                                if status == "administrator" || status == "creator" {
                                        admins = append(admins, m.GetUser().Id)
                                }
                        }

                        config.Cache.Store(cacheKey, admins)

var text string
                        if !config.Contains(admins, ctx.EffectiveUser.Id) {

text = "✅ Successfully refreshed the cache of chat admins!"

} else {

text = "⚠️ Tried refreshing the admin cache... but oops! You're not an admin."
}

	x.EditTexg(b, text , nil)
}

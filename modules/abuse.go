package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	goaway "github.com/TwiN/go-away"
)

func init() {
	Register(handlers.NewMessage(nil, DeleteAbuseHandler))
}

func DeleteAbuseHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	if message == nil || message.Text == "" {
		return nil
	}

	if goaway.IsProfane(message.Text) {
		_, err := message.Delete(bot, nil)
		if err != nil {
			return nil
		}

		censored := goaway.Censor(message.Text)
		warning := "⚠️ <b>Watch your language!</b>\nYour message was removed:\n\n`" + censored + "`"

		ctx.EffectiveChat.SendMessage(bot, warning, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	}

	return ext.ContinueGroups
}

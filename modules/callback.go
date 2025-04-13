package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

func init() {
	Register(handlers.NewCallback(callbackquery.Equal("close"), close))
}

func close(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Message.Delete(b, nil)
	return orCont(err)
}

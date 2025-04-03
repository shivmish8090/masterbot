package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func helpCB(b *gotgbot.Bot, ctx *ext.Context) error {
	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Back",
					CallbackData: "start_callback",
				},
				{
					Text:         "Close",
					CallbackData: "close",
				},
			},
		},
	}
	ctx.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{Caption: "Soon...", ReplyMarkup: keyboard})
	return nil
}

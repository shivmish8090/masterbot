package main

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func helpCB(b *gotgbot.Bot, ctx *ext.Context) error {
	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Back", CallbackData: "start_callback"},
				{Text: "Close", CallbackData: "close"},
			},
		},
	}

	_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
		Caption:     "Soon...",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Println("Failed to edit caption:", err)
	}
	return nil
}

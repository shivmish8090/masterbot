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
	helpText := `<b>/echo text</b> - Saves messages over 800 chars to Telegraph, deletes the original, and replies if used in response.`

_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
        Caption:     helpText,
        ReplyMarkup: keyboard,
        ParseMode:   "HTML",
})

	if err != nil {
		log.Println("Failed to edit caption:", err)
	}
	return nil
}

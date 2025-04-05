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
				{Text: "📝 Echo", CallbackData: "help_echo"},
				{Text: "✍️ EditMode", CallbackData: "help_editmode"},
			},
			{
				{Text: "⬅️ Back", CallbackData: "start_callback"},
				{Text: "❌ Close", CallbackData: "close"},
			},
		},
	}

	helpText := `📚 <b>Bot Command Help</b>

Here you'll find details for all available plugins and features.

👇 Tap the buttons below to view help for each module:`

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
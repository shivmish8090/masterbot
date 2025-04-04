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
{Text: "Echo", CallbackData: "help_echo"},
{Text: "EditMode", CallbackData: "help_editmode"}

},
			{
				{Text: "Back", CallbackData: "start_callback"},
				{Text: "Close", CallbackData: "close"},
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


func echoCB(b *gotgbot.Bot, ctx *ext.Context) error {
	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "Back", CallbackData: "help_callback"},
				{Text: "Close", CallbackData: "close"},
			},
		},
	}

	echoHelp := `<b>Echo Settings</b>

/echo <text> - If message is longer than 800 characters:
• Uploads to Telegraph.
• Deletes original message.
• Replies with a Telegraph link (if used in reply, it will tag the replied user).

Useful for avoiding spam and large message clutter.`

	_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
		Caption:     echoHelp,
		ReplyMarkup: keyboard,
		ParseMode:   "HTML",
	})
	if err != nil {
		log.Println("Failed to edit echo help caption:", err)
	}
	return nil
}
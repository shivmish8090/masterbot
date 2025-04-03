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
	HelpEcho := `<b>/echo text</b> - Saves messages over 800 chars to Telegraph and shares the link.  
- If a message exceeds 800 chars, it's deleted.  
- Replies to the user if used as a reply.`

	_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
		Caption:     HelpEcho,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.Println("Failed to edit caption:", err)
	}
	return nil
}

package modules

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

func init() {
	Register(handlers.NewCallback(callbackquery.Equal("help"), helpCB))
}

func helpCB(b *gotgbot.Bot, ctx *ext.Context) error {
	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{Text: "📝 Echo", CallbackData: "echo"},
				{Text: "✍️ EditMode", CallbackData: "editmode"},
			},
			{
				{Text: "⬅️ Back", CallbackData: "start_callback"},
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

func helpModuleCB(b *gotgbot.Bot, ctx *ext.Context) error {
	cbData := ctx.CallbackQuery.Data

	var helpText string
	for _, module := range ModulesHelp {
		if module.Callback == cbData {
			helpText = module.Help
			break
		}
	}

	if helpText == "" {
		helpText = "❌ No help found for this module."
	}

	_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
		Caption:   helpText,
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: "⬅️ Back", CallbackData: "help"},

				},
			},
		},
	})
	return err
}

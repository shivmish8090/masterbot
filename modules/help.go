package modules 

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
                Caption:     helpText,
                ParseMode:   "HTML",
                ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
                        InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
                                {
                                        {Text: "⬅️ Back", CallbackData: "help"},
                                        {Text: "❌ Close", CallbackData: "close"},
                                },
                        },
                },
        })

        return err
}
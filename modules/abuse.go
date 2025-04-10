package modules



func init(){

Register(handlers.NewMessage(DeleteAbuseHandler))

}
func DeleteAbuseHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
    if ctx.Message == nil || ctx.Message.Text == "" {
        return nil
    }

    msg := ctx.Message.Text

    if goaway.IsProfane(msg) {
        _, err := bot.DeleteMessage(ctx.Message.Chat.Id, ctx.Message.MessageId)
        if err != nil {
            return nil
        }

        censored := goaway.Censor(msg)
        warning := "⚠️ *Watch your language!*\nYour message was removed:\n\n`" + censored + "`"
        _, _ = ctx.EffectiveChat.SendMessage(bot, warning, &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
    }

    return nil
}
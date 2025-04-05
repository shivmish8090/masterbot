package modules

func init() {
        Register(handlers.NewCallback(callbackquery.Equal("close"), close))

}

func close(b *gotgbot.Bot, ctx *ext.Context) error {
        _, err := ctx.CallbackQuery.Message.Delete(b, nil)
        return err
}
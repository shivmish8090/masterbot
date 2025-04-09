package modules 



func init() {
        Register(handlers.NewCommand("reload", ReloadHandler))

}
func ReloadHandler(b *gotgbot.Bot, ctx *ext.Context) error {

msg := ctx.EffectiveMessage
ChatId := ctx.EffectiveChat.Id

msg.Delete(b, nil)
b.SendMessage(ChatId, "Refreshing Cache of chat admins...", nil)
return nil

}
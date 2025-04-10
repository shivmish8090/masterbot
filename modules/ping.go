package modules



func init(){

Register(handlers.NewCommand("ping", uptimeHandler))

}

func uptimeHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	uptime := time.Since(startTime)
	uptimeStr := config.FormatUptime(uptime)
ctx.EffectiveMessage.Delete(b, nil)

	_, err := ctx.EffectiveChat.SendMessage(b, "Bot has been running for: " + uptimeStr, nil)
	return err
}
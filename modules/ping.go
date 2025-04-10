package modules

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

func init() {
 fmt.Println("Ping registring")
	Register(handlers.NewCommand("ping", uptimeHandler))
}

func uptimeHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	uptime := time.Since(config.StartTime)
	uptimeStr := config.FormatUptime(uptime)
	ctx.EffectiveMessage.Delete(b, nil)

	_, err := ctx.EffectiveChat.SendMessage(b, "Bot has been running for: "+uptimeStr, nil)
	return err
}

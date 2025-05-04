package modules

import (
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/config/helpers"
)

func init() {
        Register(handlers.NewCommand("ping", modules.PingHandle))

               
}

func PingHandle(b *gotgbot.Bot, ctx *ext.Context) error {
	uptime := time.Since(config.StartTime)
	uptimeStr := helpers.FormatUptime(uptime)
	ctx.EffectiveMessage.Delete(b, nil)

	_, err := ctx.EffectiveChat.SendMessage(b, "Bot has been running for: "+uptimeStr, nil)
	return orCont(err)
}

package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func init() {
        Register(handlers.NewCommand("editmode", EditMode))

}

func EditMode(b *gotgbot.Bot, ctx *ext.Context) error {
	return nil
}

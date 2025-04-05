package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func init() {
	Register(handlers.NewCommand("editmode", EditMode))
	AddHelp("✍️ Edit Mode", "help_editmode", "...", nil)
}

func EditMode(b *gotgbot.Bot, ctx *ext.Context) error {
	return nil
}

package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func init() {
        Register(handlers.NewCommand("editmode", EditMode))
    AddHelp("✍️ Edit Mode", "help_editmode", `<b>✍️ Edit Mode</b>

<b>Command:</b> <code>/editmode &lt;off|admin|user&gt;</code>

<b>Description:</b>
Control how the bot handles <b>edited messages</b> in the group.

<b>Options:</b>
• <b>/editmode off</b> - Bot won't delete any edited messages.  
• <b>/editmode user</b> - Deletes edited messages from <b>normal users</b>. (default)  
• <b>/editmode admin</b> - Deletes edited messages from <b>both users and admins</b>. <i>(Only group owner can set this)</i>`)
}
func EditMode(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(b, "Soon..", nil)
	return nil
}

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
• <code>/editmode off</code> - Bot won't delete any edited messages.
• <code>/editmode user</code> - Deletes edited messages from <b>normal users</b>. (default)  
• <code>/editmode admin</code> - Deletes edited messages from <b>both users and admins</b>.

<b>Note:</b>
Use this to prevent spam edits or maintain message consistency in the group.`, nil)
}

func EditMode(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(b, "Soon..", nil)
	return nil
}

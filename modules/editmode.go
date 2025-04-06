package modules

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func init() {
	Register(handlers.NewCommand("editmode", EditMode))
	AddHelp("✍️ Edit Mode", "editmode", `<b>Command:</b>  
<code>/editmode</code> – Show current settings  
<code>/editmode --set-mode=&lt;off|user|admin&gt;</code>  
<code>/editmode --setduration=&lt;0-10&gt;</code>

<b>Description:</b>  
Controls how the bot deletes <b>edited messages</b>.

<b>Modes:</b>  
• <b>off</b> – No deletion  
• <b>user</b> – Delete edits from users (default)  
• <b>admin</b> – Delete edits from users & admins <i>(owner only)</i>

<b>Duration:</b>  
• <b>0</b> – Deletes immediately <i>(default)</i> & warns users  
• <b>1-10</b> – Deletes if edited after set minutes (no warning)  
• <b>&gt;10</b> – Disables deletion

<b>Example:</b>  
<code>/editmode --set-mode=user</code>  
<code>/editmode --setduration=5</code>`, nil)
}

func EditMode(b *gotgbot.Bot, ctx *ext.Context) error {
	args := ctx.Args()

	if len(args) < 2 {
		ctx.EffectiveMessage.Reply(b,
			fmt.Sprintf("Usage: <code>/editmode &lt;off|admin|user&gt;</code>\n<b>For more help, check out:</b> <a href=\"%s\">Edit Mode Help</a>",
				fmt.Sprintf("https://t.me/%s?start=editmode", b.User.Username)),
			&gotgbot.SendMessageOpts{ParseMode: "HTML"})
		ctx.EffectiveMessage.Delete(b, nil)
		return nil
	}

	y := strings.ToLower(args[1])

	if y != "off" && y != "user" && y != "admin" {
		ctx.EffectiveMessage.Reply(b,
			fmt.Sprintf("Usage: <code>/editmod &lt;off|admin|user&gt;</code>\n<b>For more help, check out:</b> <a href=\"%s\">Edit Mode Help</a>",
				fmt.Sprintf("https://t.me/%s?start=editmode", b.User.Username)),
			&gotgbot.SendMessageOpts{ParseMode: "HTML"})
		ctx.EffectiveMessage.Delete(b, nil)
		return nil
	}

	ctx.EffectiveMessage.Reply(b, "This command will be available soon..", nil)
	ctx.EffectiveMessage.Delete(b, nil)
	return nil
}

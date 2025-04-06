package modules

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/telegraph"
	"github.com/Vivekkumar-IN/EditguardianBot/utils"
)

func init() {
	Register(handlers.NewCommand("echo", EcoHandler))
	AddHelp("📝 Echo", "echo", `<b>Command:</b> 
<blockquote>/echo &lt;text&gt;
/echo --set-mode=&lt;off|manual|automatic&gt;
/echo --set-limit=&lt;number&gt;</blockquote>

<b>Description:</b>
Sends back the provided text. Also allows setting how the bot handles long messages.

<b>Echo Text:</b>  
• <b>/echo</b> &lt;text&gt; – If the message is too long, uploads it to Telegraph and sends the link.  
• <b>/echo</b> &lt;text&gt; (with reply) – Same as above, but replies to the replied message with the Telegraph link.

<b>Mode Settings:</b>
• <b>/echo</b> <code>--set-mode=off</code> – No action on long messages.  
• <b>/echo</b> <code>--set-mode=manual</code> – Deletes, warns user.  
• <b>/echo</b> <code>--set-mode=automatic</code> – Deletes, sends Telegraph link.

<b>Custom Limit:</b>  
• <b>/echo</b> <code>--set-limit=&lt;number&gt;</code> – Set character limit (default: 800).`, nil)
}

func EcoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type != "supergroup" {
		ctx.EffectiveMessage.Reply(
			b,
			"This command is meant to be used in supergroups, not in private messages!",
			nil,
		)
		return nil
	}
	args := ctx.Args()
	if len(args) < 2 {
		ctx.EffectiveMessage.Reply(b, "Usage: /echo <long message>", nil)
		return nil
	}
	ctx.EffectiveMessage.Delete(b, nil)
	keys := []string{"set-mode", "set-limit"}
	res := utils.ParseFlags(keys, true, args...)
	msg := fmt.Sprintf("set-mode: %s\nset-limit: %s", res["set-mode"], res["set-limit"])
_, err := ctx.EffectiveMessage.Reply(b, msg, nil)
	return err

	if len(ctx.EffectiveMessage.GetText()) < 800 {
		b.SendMessage(
			ctx.EffectiveChat.Id,
			"Oops! Your message is under 800 characters. You can send it without using /echo.",
			nil,
		)
		return nil
	}

	text := strings.SplitN(ctx.EffectiveMessage.GetText(), " ", 2)[1]
	url, err := telegraph.CreatePage(text, ctx.EffectiveUser.Username)
	if err != nil {
		return err
	}

	msgTemplate := `<b>Hello <a href="tg://user?id=%d">%s</a></b>, <b><a href="tg://user?id=%d">%s</a></b> wanted to share a message ✉️, but it was too long to send here 📄. You can view the full message on <b><a href="%s">Telegraph 📝</a></b>`
	linkPreviewOpts := &gotgbot.LinkPreviewOptions{IsDisabled: true}

	var msg string

	if ctx.EffectiveMessage.ReplyToMessage != nil {
		rmsg := ctx.EffectiveMessage.ReplyToMessage

		rFirst := rmsg.From.FirstName
		if rmsg.From.LastName != "" {
			rFirst += " " + rmsg.From.LastName
		}

		uFirst := ctx.EffectiveUser.FirstName
		if ctx.EffectiveUser.LastName != "" {
			uFirst += " " + ctx.EffectiveUser.LastName
		}

		msg = fmt.Sprintf(msgTemplate, rmsg.From.Id, rFirst, ctx.EffectiveUser.Id, uFirst, url)

		_, err := b.SendMessage(
			ctx.EffectiveChat.Id,
			msg,
			&gotgbot.SendMessageOpts{
				ParseMode:          "HTML",
				LinkPreviewOptions: linkPreviewOpts,
				ReplyParameters: &gotgbot.ReplyParameters{
					MessageId: rmsg.MessageId,
				},
			},
		)
		return err
	}

	uFirst := ctx.EffectiveUser.FirstName
	if ctx.EffectiveUser.LastName != "" {
		uFirst += " " + ctx.EffectiveUser.LastName
	}

	msg = fmt.Sprintf(msgTemplate, 0, "", ctx.EffectiveUser.Id, uFirst, url)

	_, err = b.SendMessage(
		ctx.EffectiveChat.Id,
		msg,
		&gotgbot.SendMessageOpts{
			ParseMode:          "HTML",
			LinkPreviewOptions: linkPreviewOpts,
		},
	)
	return err
}

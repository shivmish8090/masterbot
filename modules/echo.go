package modules

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/telegraph"
)

func init() {
	Register(handlers.NewCommand("echo", EcoHandler))
}

func EcoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type != "supergroup" {
		ctx.EffectiveMessage.Reply(
			b,
			"This command is made to be used in supergrous, not in pm!",
			nil,
		)
		return nil
	}
	if len(ctx.Args()) < 2 {
		ctx.EffectiveMessage.Reply(b, "Usage: /echo <long message>", nil)
		return nil
	}

	ctx.EffectiveMessage.Delete(b, nil)
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
	if ctx.EffectiveMessage.ReplyToMessage != nil {
 Rmsg := ctx.EffectiveMessage.ReplyToMessage
		text := fmt.Sprintf(`Hello <a href=\d>%s</a>, %s wanted to share a message ✉️, but it was too long to send here 📄. You can view the full message on <a href=%s>Telegraph 📝</a>`, "<a>)
		_ = text
		b.SendMessage(
			ctx.EffectiveChat.Id,
			url,
			&gotgbot.SendMessageOpts{
				ReplyParameters: &gotgbot.ReplyParameters{
					MessageId: ctx.EffectiveMessage.ReplyToMessage.MessageId,
				},
			},
		)
	} else {
		b.SendMessage(ctx.EffectiveChat.Id, url, nil)
	}
	return nil
}

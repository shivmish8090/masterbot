package modules

import (
	"fmt"
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
  Msg := `<b>Hello <a href="tg://user?id=%d">%s</a></b>, <b><a href="tg://user?id=%d">%s</a></b> wanted to share a message ✉️, but it was too long to send here 📄. You can view the full message on <b><a href=%s>Telegraph 📝</a></b>`
	if ctx.EffectiveMessage.ReplyToMessage != nil {
		Rmsg := ctx.EffectiveMessage.ReplyToMessage
		Msg = fmt.Sprintf(Msg, Rmsg.From.Id, Rmsg.From.FirstName+Rmsg.From.LastName, ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName+ctx.EffectiveUser.LasrName, url)
		b.SendMessage(
			ctx.EffectiveChat.Id,
			Msg,
			&gotgbot.SendMessageOpts{
				ParseMode: "HTML",
				ReplyParameters: &gotgbot.ReplyParameters{
					MessageId: ctx.EffectiveMessage.ReplyToMessage.MessageId,
				},
			},
		)
	} else {
       Msg = fmt.Sprintf(Msg, 0, "", ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName+ctx.EffectiveUser.LasrName, url)
		b.SendMessage(ctx.EffectiveChat.Id, Msg, nil)
	}
	return nil
}

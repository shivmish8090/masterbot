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
			"This command is meant to be used in supergroups, not in private messages!",
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

package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/database"
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
	ChatId := ctx.EffectiveChat.Id
	User := ctx.EffectiveUser
	Message := ctx.EffectiveMessage
	if ctx.EffectiveChat.Type != "supergroup" {
		Message.Reply(
			b,
			"This command is meant to be used in supergroups, not in private messages!",
			nil,
		)
		return nil
	}
	if len(ctx.Args()) < 2 {
		Message.Reply(b, "Usage: /echo <long message>", nil)
		return nil
	}
	Message.Delete(b, nil)

	keys := []string{"set-mode", "set-limit"}
	_, res := utils.ParseFlags(keys, Message.Text)
	if res["set-mode"] != "" || res["set-limit"] != "" {
		r := "Your settings were successfully updated:"
		settings := &database.EchoSettings{ChatID: ChatId}

		if res["set-mode"] != "" {
			settings.Mode = res["set-mode"]
			r += "\nNew Mode = " + settings.Mode
		}

		if res["set-limit"] != "" {
			limit, err := strconv.Atoi(res["set-limit"])
			if err != nil {
				if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrSyntax {
					err = fmt.Errorf("🚫 Oops! '%s' isn't a valid number.\nPlease provide a number between 200 and 4000. 🔢", res["set-limit"])
				} else {
					err = fmt.Errorf("⚠️ Something went wrong while processing the limit.\nError: %v", err)
				}
			} else if limit < 200 || limit > 4000 {
				err = fmt.Errorf("⚠️ The number %d is out of range!\nPlease provide a number between 200 and 4000. 📏", limit)
			}

			if err != nil {
				b.SendMessage(ChatId, err.Error(), nil)
				return err
			}

			settings.Limit = limit
			r += "\nNew Limit = " + strconv.Itoa(settings.Limit)
		}

		err := database.SetEchoSettings(settings)
		if err != nil {
			b.SendMessage(
				ChatId,
				fmt.Sprintf("Something went wrong while saving settings\nError: %v", err),
				nil,
			)
			return err
		}

		b.SendMessage(ChatId, r, nil)
		return nil
	}

	settings, err := database.GetEchoSettings(ChatId)
	if err != nil {
		_, err = b.SendMessage(
			ChatId,
			fmt.Sprintf("⚠️ Something went wrong while processing the limit.\nError: %v", err),
			nil,
		)
		return err
	}
	limit := settings.Limit
	if len(Message.GetText()) < limit {
		b.SendMessage(
			ChatId,
			fmt.Sprintf("Oops! Your message is under %d characters. You can send it without using /echo.", limit),
			nil,
		)
		return nil
	}

	text := strings.SplitN(Message.GetText(), " ", 2)[1]
	url, err := telegraph.CreatePage(text, User.Username)
	if err != nil {
		return err
	}

	msgTemplate := `<b>Hello <a href="tg://user?id=%d">%s</a></b>, <b><a href="tg://user?id=%d">%s</a></b> wanted to share a message ✉️, but it was too long to send here 📄. You can view the full message on <b><a href="%s">Telegraph 📝</a></b>`
	linkPreviewOpts := &gotgbot.LinkPreviewOptions{IsDisabled: true}

	var msg string

	if Message.ReplyToMessage != nil {
		rmsg := Message.ReplyToMessage

		rFirst := rmsg.From.FirstName
		if rmsg.From.LastName != "" {
			rFirst += " " + rmsg.From.LastName
		}

		uFirst := User.FirstName
		if User.LastName != "" {
			uFirst += " " + User.LastName
		}

		msg = fmt.Sprintf(msgTemplate, rmsg.From.Id, rFirst, User.Id, uFirst, url)

		_, err := b.SendMessage(
			ChatId,
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

	uFirst := User.FirstName
	if User.LastName != "" {
		uFirst += " " + User.LastName
	}

	msg = fmt.Sprintf(msgTemplate, 0, "", User.Id, uFirst, url)

	_, err = b.SendMessage(
		ChatId,
		msg,
		&gotgbot.SendMessageOpts{
			ParseMode:          "HTML",
			LinkPreviewOptions: linkPreviewOpts,
		},
	)
	return err
}

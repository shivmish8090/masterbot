package modules

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/config/helpers"
	"github.com/Vivekkumar-IN/EditguardianBot/database"
	"github.com/Vivekkumar-IN/EditguardianBot/telegraph"
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

type warningTracker struct {
	sync.Mutex
	chats map[int64]time.Time
}

var deleteWarningTracker = warningTracker{
	chats: make(map[int64]time.Time),
}

func EcoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	ChatId := ctx.EffectiveChat.Id
	User := ctx.EffectiveUser
	Message := ctx.EffectiveMessage

	if Message.SenderChat != nil {
		Message.Reply(
			b,
			"You are anonymous Admin you can't use this command.",
			nil,
		)
		return Continue
	}

	if ctx.EffectiveChat.Type != "supergroup" {
		Message.Reply(
			b,
			"This command is meant to be used in supergroups, not in private messages!",
			nil,
		)
		return Continue
	}

	if len(ctx.Args()) < 2 {
		Message.Reply(b, "Usage: /echo <long message>", nil)
		return Continue
	}

	Message.Delete(b, nil)

	keys := []string{"set-mode", "set-limit"}
	_, res := helpers.ParseFlags(keys, Message.Text)

	if res["set-mode"] != "" || res["set-limit"] != "" {
		admins, err := helpers.GetAdmins(b, ChatId)
		if err != nil {
			return err
		}

		if !config.Contains(admins, User.Id) {
			b.SendMessage(ChatId, "Access denied: Only group admins can use this command.\n\nIf you are an admin, please use /reload to refresh the admin list.", nil)
			return Continue
		}

		r := "Your settings were successfully updated:"
		settings := &database.EchoSettings{ChatID: ChatId}

		if res["set-mode"] != "" {
			mode := strings.ToLower(res["set-mode"])
			if mode == "off" || mode == "manual" || mode == "automatic" {
				settings.Mode = mode
				r += "\nNew Mode = " + settings.Mode
			} else {
				b.SendMessage(ChatId, fmt.Sprintf("🚫 Invalid mode: '%s'. Valid options are <off|manual|automatic>.", res["set-mode"]), nil)
				return Continue
			}
		}

		var err error
		if res["set-limit"] != "" {
			limit, convErr := strconv.Atoi(res["set-limit"])
			if convErr != nil {
				if numErr, ok := convErr.(*strconv.NumError); ok && numErr.Err == strconv.ErrSyntax {
					err = fmt.Errorf("🚫 Oops! '%s' isn't a valid number.\nPlease provide a number between 200 and 4000. 🔢", res["set-limit"])
				} else {
					err = fmt.Errorf("⚠️ Something went wrong while processing the limit.\nError: %v", convErr)
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

		err = database.SetEchoSettings(settings)
		if err != nil {
			b.SendMessage(ChatId, fmt.Sprintf("Something went wrong while saving settings\nError: %v", err), nil)
			return err
		}

		b.SendMessage(ChatId, r, nil)
		return Continue
	}

	settings, err := database.GetEchoSettings(ChatId)
	if err != nil {
		b.SendMessage(ChatId, fmt.Sprintf("⚠️ Something went wrong while processing the limit.\nError: %v", err), nil)
		return err
	}

	if len(Message.GetText()) < settings.Limit {
		b.SendMessage(ChatId, fmt.Sprintf("Oops! Your message is under %d characters. You can send it without using /echo.", settings.Limit), nil)
		return Continue
	}

	text := strings.SplitN(Message.GetText(), " ", 2)[1]

	err = sendEchoMessage(b, ctx, text)
	return orCont(err)
}

func DeleteLongMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage
	settings, err := database.GetEchoSettings(ctx.EffectiveChat.Id)
	var isAutomatic bool
	if err != nil {
		b.SendMessage(
			config.LoggerId,
			fmt.Sprintf("⚠️ Something went wrong while Getting the limit.\nError: %v", err),
			nil,
		)
		return err
	}

	if m.GetText() == "" || len(m.GetText()) < settings.Limit {
		return Continue
	}
	if settings.Mode == "OFF" {
		return Continue
	}
	if settings.Mode == "AUTOMATIC" {
		isAutomatic = true
	}
	done, err := ctx.EffectiveMessage.Delete(b, nil)
	if err != nil {
		fmt.Println("Delete error:", err)
		return err
	}

	if done && !isAutomatic {
		deleteWarningTracker.Lock()
		defer deleteWarningTracker.Unlock()

		lastWarning, exists := deleteWarningTracker.chats[ctx.EffectiveChat.Id]
		if !exists || time.Since(lastWarning) > time.Second {
			text := fmt.Sprintf(`
⚠️ <a href="tg://user?id=%d">%s</a>, your message exceeds the %d-character limit! 🚫  
Please shorten it before sending. ✂️  

Alternatively, use /echo for sending longer messages. 📜
`, ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName, settings.Limit)

			_, err := b.SendMessage(
				ctx.EffectiveChat.Id,
				text,
				&gotgbot.SendMessageOpts{ParseMode: "HTML"},
			)
			if err != nil {
				fmt.Println("SendMessage error:", err)
				return err
			}
			deleteWarningTracker.chats[ctx.EffectiveChat.Id] = time.Now()
		}
	} else if done && isAutomatic {
		err = sendEchoMessage(b, ctx, m.GetText())
		return orCont(err)
	}
	return Continue
}

func sendEchoMessage(b *gotgbot.Bot, ctx *ext.Context, text string) error {
	User := ctx.EffectiveUser
	userFullName := strings.TrimSpace(User.FirstName + " " + User.LastName)

	var authorURL string
	if User.Username != "" {
		authorURL = fmt.Sprintf("https://t.me/%s", User.Username)
	} else {
		authorURL = fmt.Sprintf("tg://user?id=%d", User.Id)
	}

	url, err := telegraph.CreatePage(text, userFullName, authorURL)
	if err != nil {
		return err
	}

	msgTemplate := `<b>Hello <a href="tg://user?id=%d">%s</a></b>, <b><a href="tg://user?id=%d">%s</a></b> wanted to share a message ✉️, but it was too long to send here 📄. You can view the full message on <b><a href="%s">Telegraph 📝</a></b>`
	var msg string

	opts := &gotgbot.SendMessageOpts{
		ParseMode:          "HTML",
		LinkPreviewOptions: &gotgbot.LinkPreviewOptions{IsDisabled: true},
	}

	if rmsg := ctx.EffectiveMessage.ReplyToMessage; rmsg != nil && rmsg.From != nil {
		replyUserFullName := strings.TrimSpace(rmsg.From.FirstName + " " + rmsg.From.LastName)
		msg = fmt.Sprintf(msgTemplate, rmsg.From.Id, replyUserFullName, User.Id, userFullName, url)

		opts.ReplyParameters = &gotgbot.ReplyParameters{
			MessageId: rmsg.MessageId,
		}
	} else {
		msg = fmt.Sprintf(msgTemplate, 0, "", User.Id, userFullName, url)
	}

	_, err = b.SendMessage(ctx.EffectiveChat.Id, msg, opts)
	return orCont(err)
}

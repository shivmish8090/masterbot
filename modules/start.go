package modules

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/database"
)

func init() {
	RegisterCommand(handlers.NewCommand("start", start))

	RegisterCallback(handlers.NewCallback(callbackquery.Equal("start_callback"), start))
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	isCallback := ctx.CallbackQuery != nil
	chat := ctx.EffectiveChat.Type

	if chat == "private" {
		if !isCallback {
			ctx.EffectiveMessage.Delete(b, nil)
		}

		file := gotgbot.InputFileByURL(config.StartImage)
		caption := fmt.Sprintf(
			`<b>🚀 Hello <a href="tg://user?id=%d">%s</a>! 👋</b>  
I'm <b><a href="tg://user?id=%d">%s</a></b>, your security assistant, ensuring a safe and transparent environment for our discussions!  

🔍 <b>Edited Message Protection:</b>  
✂️ Messages that are edited will be <b>automatically deleted</b> to maintain clarity and honesty.  

🚨 <b>Real-Time Notifications:</b>  
📢 You'll receive an alert whenever a message is removed, keeping you informed at all times.  

💡 <b>Getting Started is Easy!</b>  
1️⃣ <b>Add me to your group.</b>  
2️⃣ I'll start <b>protecting your chat instantly!</b>  

🔐 <b>Keep your group safe now!</b>  
➡️ Tap <b>"Add Group"</b> to enable my security features today!`,
			ctx.EffectiveUser.Id,
			ctx.EffectiveUser.FirstName+" "+ctx.EffectiveUser.LastName,
			b.User.Id,
			b.User.FirstName,
		)

		keyboard := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{
						Text: "🔄 Update Channel",
						Url:  "https://t.me/Dns_Official_Channel",
					},
					{
						Text: "💬 Update Group",
						Url:  "https://t.me/dns_support_group",
					},
				},
				{
					{
						Text:         "❓ Help & Commands",
						CallbackData: "help_callback",
					},
				},
				{
					{
						Text: "➕ Add me to Your Group",
						Url: fmt.Sprintf(
							"https://t.me/%s?startgroup=s&admin=delete_messages+invite_users",
							b.User.Username,
						),
					},
				},
			},
		}

		if isCallback {
			_, _, err := ctx.CallbackQuery.Message.EditCaption(b, &gotgbot.EditMessageCaptionOpts{
				Caption:     caption,
				ParseMode:   "HTML",
				ReplyMarkup: keyboard,
			})
			if err != nil {
				return err
			}
		} else {
			database.AddServedUser(ctx.EffectiveUser.Id)
			_, err := b.SendPhoto(
				ctx.EffectiveChat.Id,
				file,
				&gotgbot.SendPhotoOpts{
					Caption:        caption,
					ProtectContent: true,
					ParseMode:      "HTML",
					ReplyMarkup:    keyboard,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to send photo: %w", err)
			}

			logStr := fmt.Sprintf(
				`<a href="tg://user?id=%d">%s</a> has started the bot.

<b>User ID:</b> <code>%d</code>
<b>User Name:</b> %s %s`,
				ctx.EffectiveUser.Id,
				ctx.EffectiveUser.FirstName,
				ctx.EffectiveUser.Id,
				ctx.EffectiveUser.FirstName,
				ctx.EffectiveUser.LastName,
			)
			b.SendMessage(
				config.LoggerId,
				logStr,
				&gotgbot.SendMessageOpts{ParseMode: "HTML"},
			)
		}

	} else if chat == "group" {
		if isCallback {
			return nil
		}

		message := `⚠️ Warning: I can't function in a basic group!

To use my features, please upgrade this group to a supergroup.

✅ How to upgrade:
1. Go to Group Settings.
2. Tap on "Chat History" and set it to "Visible".
3. Re-add me, and I'll be ready to help!`

		ctx.EffectiveMessage.Reply(b, message, nil)
		ctx.EffectiveChat.Leave(b, nil)

	} else if chat == "supergroup" {
		if isCallback {
			return nil
		}
		database.AddServedChat(ctx.EffectiveChat.Id)
		ctx.EffectiveMessage.Reply(b, "✅ I am active and ready to protect this supergroup!", nil)
	}

	return nil
}

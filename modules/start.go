package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/config/buttons"
	"github.com/Vivekkumar-IN/EditguardianBot/database"
)

func init() {
	Register(handlers.NewCommand("start", start))

	Register(handlers.NewCallback(callbackquery.Equal("start_callback"), start))
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	isCallback := ctx.CallbackQuery != nil
	chat := ctx.EffectiveChat.Type

	if chat == "private" {
		if !isCallback {
			ctx.EffectiveMessage.Delete(b, nil)
		}
		if len(ctx.Args()) >= 2 {
			modName := ctx.Args()[1]
			if strings.HasPrefix(modName, "info_") {

				userId := strings.Split(modName, "_")[1]
				userID, e := strconv.ParseInt(userId, 10, 64)
				if e != nil {
					return e
				}
userInfo, er := b.GetChat(userID, nil)
if er != nil {
return er

}
info := fmt.Sprintf(`
Name: %s
Id: %d
Link: <a href="tg://user?id=%d">Link 1</a> <a href="tg://openmessage?user_id=%d">Link 2</a>
`, strings.TrimSpace(userInfo.FirstName+" "+userInfo.LastName), userInfo.Id, userInfo.Id, userInfo.Id)

				b.SendMessage(ctx.EffectiveChat.Id, info, &gotgbot.SendMessageOpts{
					ParseMode: "HTML",
				})

			}
			helpString := GetHelp(modName)
			if helpString != "" {
				_, err := b.SendMessage(ctx.EffectiveChat.Id, helpString, &gotgbot.SendMessageOpts{
					ParseMode: "HTML",
				})
				if err != nil {
					return err
				}
				return Continue
			}
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

		keyboard := buttons.StartPanel(b)
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
					Caption:     caption,
					ParseMode:   "HTML",
					ReplyMarkup: keyboard,
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
			return Continue
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
			return Continue
		}
		database.AddServedChat(ctx.EffectiveChat.Id)
		ctx.EffectiveMessage.Reply(b, "✅ I am active and ready to protect this supergroup!", nil)
	}

	return Continue
}

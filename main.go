package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

func main() {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// /start command to introduce the bot
	dispatcher.AddHandler(handlers.NewCommand("start", start))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat.Type
	if chat == "private" {
		file := gotgbot.InputFileByURL(config.StartImage)

		caption := fmt.Sprintf(
			`Hello %s 👋, I'm your 𝗘𝗱𝗶𝘁 𝗚𝘂𝗮𝗿𝗱𝗶𝗮𝗻 𝗕𝗼𝘁, here to maintain a secure environment for our discussions.

🚫 𝗘𝗱𝗶𝘁𝗲𝗱 𝗠𝗲𝘀𝘀𝗮𝗴𝗲 𝗗𝗲𝗹𝗲𝘁𝗶𝗼𝗻: 𝗜'𝗹𝗹 𝗿𝗲𝗺𝗼𝘃𝗲 𝗲𝗱𝗶𝘁𝗲𝗱 𝗺𝗲𝘀𝘀𝗮𝗴𝗲𝘀 𝘁𝗼 𝗺𝗮𝗶𝗻𝘁𝗮𝗶𝗻 𝘁𝗿𝗮𝗻𝘀𝗽𝗮𝗿𝗲𝗻𝗰𝘆.

📣 𝗡𝗼𝘁𝗶𝗳𝗶𝗰𝗮𝘁𝗶𝗼𝗻𝘀: 𝗬𝗼𝘂'𝗹𝗹 𝗯𝗲 𝗶𝗻𝗳𝗼𝗿𝗺𝗲𝗱 𝗲𝗮𝗰𝗵 𝘁𝗶𝗺𝗲 𝗮 𝗺𝗲𝘀𝘀𝗮𝗴𝗲 𝗶𝘀 𝗱𝗲𝗹𝗲𝘁𝗲𝗱.

🌟 𝗚𝗲𝘁 𝗦𝘁𝗮𝗿𝘁𝗲𝗱:
1. Add me to your group.
2. I'll start protecting instantly.

➡️ Click on 𝗔𝗱𝗱 𝗚𝗿𝗼𝘂𝗽 to add me and keep our group safe!`,
			b.User.Username,
		)

		keyboard := gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: "🔄 Update Channel", Url: "https://t.me/Dns_Official_Channel"},
					{Text: "💬 Update Group", Url: "https://t.me/dns_support_group"},
				},
				{
					{
						Text: "➕ Add me to Your Group",
						Url:  fmt.Sprintf("https://t.me/%s?startgroup=s&admin=delete_messages+invite_users", b.User.Username),
					},
				},
			},
		}

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
			ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName,
			ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName, ctx.EffectiveUser.LastName,
		)
		b.SendMessage(config.LoggerId, logStr, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	} else if chat == "group" {
		message := `⚠️ Warning: I can't function in a basic group!

To use my features, please upgrade this group to a supergroup.

✅ How to upgrade:
1. Go to Group Settings.
2. Tap on "Chat History" and set it to "Visible".
3. Re-add me, and I'll be ready to help!`

		ctx.EffectiveMessage.Reply(b, message, nil)
		ctx.EffectiveChat.Leave(b, nil)
	} else if chat == "supergroup" {
		ctx.EffectiveMessage.Reply(b, "✅ I am active and ready to protect this supergroup!", nil)

 chatMemberCount, err := b.GetChatMemberCount(ctx.EffectiveChat.Id)
if err != nil {
chatMemberCount = "None"
}
 logStr := fmt.Sprintf(
    `🔹 <b>Group Connection Log</b> 🔹  
━━━━━━━━━━━━━━━━━━━━━━  
📌 <b>Group Name:</b> %s  
🆔 <b>Group ID:</b> <code>%d</code>  
🔗 <b>Username:</b> @%s  
👥 <b>Members:</b> %d  
━━━━━━━━━━━━━━━━━━━━━━`,  
    ctx.EffectiveChat.Title,  
    ctx.EffectiveChat.Id,  
    ctx.EffectiveChat.Username,  
    chatMemberCount
) 
b.SendMessage(config.LoggerId, logStr, &gotgbot.SendMessageOpts{ParseMode: "HTML"})

	}
	return nil
}

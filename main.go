package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/filters"
	"github.com/Vivekkumar-IN/EditguardianBot/telegraph"
)

var deleteWarningTracker = struct {
	sync.Mutex
	chats map[int64]time.Time
}{chats: make(map[int64]time.Time)}

func main() {
	b, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	filters.Init(b)

	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewMyChatMember(
		func(u *gotgbot.ChatMemberUpdated) bool {
			wasMember, isMember := ExtractJoinLeftStatusChange(u)
			return !wasMember && isMember
		},
		AddedToGroups,
	))

	evalHandler := handlers.NewMessage(
		filters.AndFilter(filters.Owner, filters.Command("eval")),
		EvalHandler,
	).SetAllowEdited(true)

	dispatcher.AddHandler(evalHandler)
	dispatcher.AddHandler(handlers.NewCommand("echo", EcoHandler))
	dispatcher.AddHandlerToGroup(handlers.NewMessage(
		filters.Invert(filters.ChatAdmins),
		deleteEditedMessage,
	).SetAllowEdited(true), -1)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(
		filters.LongMessage,
		deleteLongMessage,
	), -1)

	allowedUpdates := []string{
		"message",
		"my_chat_member",
		"chat_member",
		"edited_message",
	}

	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout:        9,
			AllowedUpdates: allowedUpdates,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	log.Printf("%s has been started...\n", b.User.Username)
	b.SendMessage(
		config.LoggerId,
		fmt.Sprintf("%s has started\n", b.User.Username),
		nil,
	)
	updater.Idle()
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat.Type

	if chat == "private" {
		file := gotgbot.InputFileByURL(config.StartImage)

		caption := fmt.Sprintf(
			`Hello %s 👋, I'm your %s, here to maintain a secure environment for our discussions.

🚫 𝗘𝗱𝗶𝘁𝗲𝗱 𝗠𝗲𝘀𝘀𝗮𝗴𝗲 𝗗𝗲𝗹𝗲𝘁𝗶𝗼𝗻: 𝗜'𝗹𝗹 𝗿𝗲𝗺𝗼𝘃𝗲 𝗲𝗱𝗶𝘁𝗲𝗱 𝗺𝗲𝘀𝘀𝗮𝗴𝗲𝘀 𝘁𝗼 𝗺𝗮𝗶𝗻𝘁𝗮𝗶𝗻 𝘁𝗿𝗮𝗻𝘀𝗽𝗮𝗿𝗲𝗻𝗰𝘆.

📣 𝗡𝗼𝘁𝗶𝗳𝗶𝗰𝗮𝘁𝗶𝗼𝗻𝘀: 𝗬𝗼𝘂'𝗹𝗹 𝗯𝗲 𝗶𝗻𝗳𝗼𝗿𝗺𝗲𝗱 𝗲𝗮𝗰𝗵 𝘁𝗶𝗺𝗲 𝗮 𝗺𝗲𝘀𝘀𝗮𝗴𝗲 𝗶𝘀 𝗱𝗲𝗹𝗲𝘁𝗲𝗱.

🌟 𝗚𝗲𝘁 𝗦𝘁𝗮𝗿𝘁𝗲𝗱:
1. Add me to your group.
2. I'll start protecting instantly.

➡️ Click on 𝗔𝗱𝗱 𝗚𝗿𝗼𝘂𝗽 to add me and keep our group safe!`,
			ctx.EffectiveUser.FirstName,
			b.User.Username,
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
						Text:         "Help & Commands",
						CallbackData: "help",
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
	}
	return ext.EndGroups
}

func AddedToGroups(b *gotgbot.Bot, ctx *ext.Context) error {
	text := fmt.Sprintf(
		`Hello 👋 I'm <b>%s</b>, here to help keep the chat transparent and secure.

🚫 I will automatically delete edited messages to maintain clarity.  

I'm ready to protect this group! ✅  
Let me know if you need any help.`,
		b.User.FirstName,
	)

	b.SendMessage(
		ctx.EffectiveChat.Id,
		text,
		&gotgbot.SendMessageOpts{ParseMode: "HTML"},
	)
	chatMemberCount, err := b.GetChatMemberCount(ctx.EffectiveChat.Id, nil)
	if err != nil {
		chatMemberCount = 0
	}

	groupUsername := ctx.EffectiveChat.Username
	if groupUsername == "" {
		groupUsername = "N/A"
	}

	groupTitle := ctx.EffectiveChat.Title
	if groupTitle == "" {
		groupTitle = "Unknown"
	}

	logStr := fmt.Sprintf(
		`🔹 <b>Group Connection Log</b> 🔹  
━━━━━━━━━━━━━━━━━  
📌 <b>Group Name:</b> %s  
🆔 <b>Group ID:</b> <code>%d</code>  
🔗 <b>Username:</b> @%s  
👥 <b>Members:</b> %d  
━━━━━━━━━━━━━━━━━`,
		groupTitle,
		ctx.EffectiveChat.Id,
		groupUsername,
		chatMemberCount,
	)

	_, err = b.SendMessage(
		config.LoggerId,
		logStr,
		&gotgbot.SendMessageOpts{ParseMode: "HTML"},
	)
	if err != nil {
		return err
	}

	return nil
}

func deleteLongMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	done, err := ctx.EffectiveMessage.Delete(b, nil)
	if done {
		deleteWarningTracker.Lock()
		lastWarning, exists := deleteWarningTracker.chats[ctx.EffectiveChat.Id]
		if !exists || time.Since(lastWarning) > time.Second {
			text := fmt.Sprintf(`
⚠️ <a href="tg://user?id=%d">%s</a>, your message exceeds the 800-character limit! 🚫  
Please shorten it before sending. ✂️  

Alternatively, use /eco for sending longer messages. 📜
`, ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName)

			_, err := b.SendMessage(
				ctx.EffectiveChat.Id,
				text,
				&gotgbot.SendMessageOpts{ParseMode: "HTML"},
			)
			if err != nil {
				return err
			}
			deleteWarningTracker.chats[ctx.EffectiveChat.Id] = time.Now()
		}
		deleteWarningTracker.Unlock()
	} else {
		return err
	}
	return ext.EndGroups
}

func EcoHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveChat.Type != "supergroup" {
		ctx.EffectiveMessage.Reply(
			b,
			"This command can be used only in groups",
			nil,
		)
		return nil
	}
	if len(ctx.Args()) < 2 {
		ctx.EffectiveMessage.Reply(b, "Usage: /echo <long message>", nil)
		return nil
	}

	ctx.EffectiveMessage.Delete(b, nil)
	if len(ctx.EffectiveMessage.GetText()) < 500 {
		b.SendMessage(
			ctx.EffectiveChat.Id,
			"Oops! Your message is under 500 characters. You can send it without using /echo.",
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

func deleteEditedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EditedMessage != nil {
		_, err := ctx.EffectiveMessage.Delete(b, nil)
		if err != nil {
			return err
		}

		_, err = b.SendMessage(
			ctx.EffectiveChat.Id,
			"⚠️ Editing messages is not allowed!",
			nil,
		)
		return err
	}
	return nil
}

func ExtractJoinLeftStatusChange(u *gotgbot.ChatMemberUpdated) (bool, bool) {
	if u.Chat.Type == "channel" {
		return false, false
	}

	oldMemberStatus := u.OldChatMember.MergeChatMember().Status
	newMemberStatus := u.NewChatMember.MergeChatMember().Status
	oldIsMember := u.OldChatMember.MergeChatMember().IsMember
	newIsMember := u.NewChatMember.MergeChatMember().IsMember

	if oldMemberStatus == newMemberStatus {
		return false, false
	}

	findInSlice := func(slice []string, val string) bool {
		for _, item := range slice {
			if item == val {
				return true
			}
		}
		return false
	}

	wasMember := findInSlice(
		[]string{"member", "administrator", "creator"},
		oldMemberStatus,
	) ||
		(oldMemberStatus == "restricted" && oldIsMember)

	isMember := findInSlice(
		[]string{"member", "administrator", "creator"},
		newMemberStatus,
	) ||
		(newMemberStatus == "restricted" && newIsMember)

	return wasMember, isMember
}

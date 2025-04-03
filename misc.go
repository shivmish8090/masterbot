package main

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

func closeCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.CallbackQuery.Message.Delete(b, nil)
	return err
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

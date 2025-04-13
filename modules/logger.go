package modules

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

func init() {
	Register(handlers.NewMyChatMember(
		func(u *gotgbot.ChatMemberUpdated) bool {
			wasMember, isMember := ExtractJoinLeftStatusChange(u)
			return !wasMember && isMember
		},
		AddedToGroups,
	))
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
	return orCont(err)
}

func ExtractJoinLeftStatusChange(u *gotgbot.ChatMemberUpdated) (bool, bool) {
	if u.Chat.Type == "channel" {
		return false, false
	}

	old := u.OldChatMember.MergeChatMember()
	new := u.NewChatMember.MergeChatMember()

	if old.Status == new.Status {
		return false, false
	}

	memberStatuses := map[string]struct{}{
		"member":        {},
		"administrator": {},
		"creator":       {},
	}

	_, wasIn := memberStatuses[old.Status]
	wasMember := wasIn || (old.Status == "restricted" && old.IsMember)

	_, isIn := memberStatuses[new.Status]
	isMember := isIn || (new.Status == "restricted" && new.IsMember)

	return wasMember, isMember
}

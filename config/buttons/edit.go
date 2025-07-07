package buttons

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// return nil for disabling button
func EditedMessagePanel(b *gotgbot.Bot) gotgbot.InlineKeyboardMarkup {
	btn := &Button{RowWidth: 2}

	btn.Add(
		btn.Url("🔄 Updates", "https://t.me/NeoUpdatess"),
		btn.Url("💬 Support", "https://t.me/best_friends_chatting_grpz0"),
	)
	btn.Row(
		btn.Url(
			"➕ Add me to Your Group",
			fmt.Sprintf("https://t.me/%s?startgroup=s&admin=delete_messages+invite_users", b.User.Username),
		),
	)

	return btn.Build()
}

package buttons

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
        "github.com/Vivekkumar-IN/EditguardianBot/config"
)

func StartPanel(b *gotgbot.Bot) gotgbot.InlineKeyboardMarkup {
	btn := &Button{RowWidth: 2}

	btn.Add(
		btn.Url("🔄 Update Channel", config.Channel),
		btn.Url("💬 Update Group", config.Chat),
	)

	btn.Row(
		btn.Inline("❓ Help & Commands", "help"),
	)

	btn.Row(
		btn.Url(
			"➕ Add me to Your Group",
			fmt.Sprintf("https://t.me/%s?startgroup=s&admin=delete_messages+invite_users", b.User.Username),
		),
	)

	return btn.Build()
}

func NormalStartPanel(b *gotgbot.Bot) gotgbot.InlineKeyboardMarkup {
	return StartPanel(b)
}

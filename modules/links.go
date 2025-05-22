package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func DeleteLinkMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	_, err := m.Delete(b, nil)
	if err != nil {
		return err
	}
	b.SendMessage(
		m.Chat.Id,
		"⚠️ Direct URLs aren't allowed.\nUse like <a href='https://t.me/dns_support_group'>this</a>",
		&gotgbot.SendMessageOpts{ParseMode: "HTML", LinkPreviewOptions: &gotgbot.LinkPreviewOptions{IsDisabled: true}},
	)

	return Continue
}

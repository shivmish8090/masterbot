package modules

import (
	"fmt"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type warningTracker struct {
	sync.Mutex
	chats map[int64]time.Time
}

var deleteWarningTracker = warningTracker{
	chats: make(map[int64]time.Time),
}

func DeleteEditedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EditedMessage
	if message == nil {
		return nil
	}

	if _, err := ctx.EffectiveMessage.Delete(b, nil); err != nil {
		return err
	}
reason := "<b>🚫 Editing messages is prohibited in this chat. Please refrain from modifying your messages to maintain the integrity of the conversation.</b>"

if message.Text != "" {
    reason = "<b>🚫 Editing messages is prohibited in this chat. Please avoid changing the text content once it's sent to maintain the flow of conversation.</b>"
} else if message.Caption != "" {
    reason = "<b>✍️ Editing a caption is restricted. Once the caption is set, it cannot be changed to ensure clarity and consistency in the content.</b>"
} else if message.Photo != nil {
    reason = "<b>📷 Replacing or editing a photo is not permitted. To keep the visual context intact, altering images after posting is not allowed.</b>"
} else if message.Video != nil {
    reason = "<b>🎥 Replacing or editing a video is not allowed. Once uploaded, videos should not be modified to preserve the original content.</b>"
} else if message.Document != nil {
    reason = "<b>📄 Replacing a document is restricted. Once shared, documents cannot be edited or replaced to ensure accuracy and trust in the information.</b>"
} else if message.Audio != nil {
    reason = "<b>🎵 Replacing an audio file is not permitted. For consistency, audio files cannot be edited after being uploaded.</b>"
} else if message.VideoNote != nil {
    reason = "<b>📹 Changing a video note is not allowed. Video notes must remain as originally sent to keep the communication intact.</b>"
} else if message.Voice != nil {
    reason = "<b>🎙️ Editing a voice message is not permitted. To maintain the original intent, voice recordings should not be altered.</b>"
} else if message.Animation != nil {
    reason = "<b>🎞️ Modifying a GIF is not allowed. Once sent, animated images must remain unchanged to preserve the context of the conversation.</b>"
} else if message.Sticker != nil {
    reason = "<b>🖼️ Replacing a sticker is not permitted. Stickers are meant to be sent as-is, and cannot be edited after posting.</b>"
}

	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text: "🔄 Updates",
					Url:  "https://t.me/Dns_Official_Channel",
				},
				{
					Text: "💬 Support",
					Url:  "https://t.me/dns_support_group",
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

	_, err := b.SendMessage(
		ctx.EffectiveChat.Id,
		reason,
		&gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: keyboard},
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteLongMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	fmt.Println("deleteLongMessage triggered")

	done, err := ctx.EffectiveMessage.Delete(b, nil)
	if err != nil {
		fmt.Println("Delete error:", err)
		return err
	}

	if done {
		deleteWarningTracker.Lock()
		defer deleteWarningTracker.Unlock()

		lastWarning, exists := deleteWarningTracker.chats[ctx.EffectiveChat.Id]
		if !exists || time.Since(lastWarning) > time.Second {
			text := fmt.Sprintf(`
⚠️ <a href="tg://user?id=%d">%s</a>, your message exceeds the 800-character limit! 🚫  
Please shorten it before sending. ✂️  

Alternatively, use /echo for sending longer messages. 📜
`, ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName)

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
	}
	return nil
}

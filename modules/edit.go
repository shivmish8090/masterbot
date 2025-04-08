package modules

import (
	"fmt"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
	"github.com/Vivekkumar-IN/EditguardianBot/database"
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
	if message == nil || ctx.EffectiveChat.Type == "private" {
		return nil
	}

	Chat := b.GetChat(ctx.EffectiveChat.Id)

	if message.SenderChat != nil {
		if message.Chat.Id == message.SenderChat.Id || Chat.LinkedChatId == message.SenderChat.Id {
			return nil
		}
	}
	if _, err := ctx.EffectiveMessage.Delete(b, nil); err != nil {
		return err
	}
	reason := "<b>🚫 Editing messages is prohibited in this chat.</b> Please refrain from modifying your messages to maintain the integrity of the conversation."

	if message.Text != "" {
		reason = "<b>🚫 Editing messages is prohibited in this chat.</b> Please avoid changing the text content once it's sent to maintain the flow of conversation."
	} else if message.Caption != "" {
		reason = "<b>✍️ Editing a caption is restricted.</b> Once the caption is set, it cannot be changed to ensure clarity and consistency in the content."
	} else if message.Photo != nil {
		reason = "<b>📷 Replacing or editing a photo is not permitted.</b> Altering images after posting is not allowed to keep the visual context intact."
	} else if message.Video != nil {
		reason = "<b>🎥 Replacing or editing a video is not allowed.</b> Videos should not be modified after posting to preserve the original content."
	} else if message.Document != nil {
		reason = "<b>📄 Replacing a document is restricted.</b> Documents cannot be edited or replaced to ensure accuracy and trust in the information."
	} else if message.Audio != nil {
		reason = "<b>🎵 Replacing an audio file is not permitted.</b> Audio files cannot be edited after being uploaded for consistency."
	} else if message.VideoNote != nil {
		reason = "<b>📹 Changing a video note is not allowed.</b> Video notes must remain as originally sent to keep the communication intact."
	} else if message.Voice != nil {
		reason = "<b>🎙️ Editing a voice message is not permitted.</b> Voice recordings should not be altered to maintain the original intent."
	} else if message.Animation != nil {
		reason = "<b>🎞️ Modifying a GIF is not allowed.</b> GIFs must remain unchanged after being sent to preserve the context of the conversation."
	} else if message.Sticker != nil {
		reason = "<b>🖼️ Replacing a sticker is not permitted.</b> Stickers cannot be edited after posting to maintain their original meaning."
	}

	keyboard := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text: "🔄 Updates",
					Url:  "https://t.me/SanatanVibe",
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
	m := ctx.EffectiveMessage
	settings, err := database.GetEchoSettings(ctx.EffectiveChat.Id)
	if err != nil {
		_, err = b.SendMessage(
			config.LoggerId,
			fmt.Sprintf("⚠️ Something went wrong while Getting the limit.\nError: %v", err),
			nil,
		)
		return err
	}

	if m.GetText() == "" || len(m.GetText()) < settings.Limit {
		return nil
	}

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
⚠️ <a href="tg://user?id=%d">%s</a>, your message exceeds the %d-character limit! 🚫  
Please shorten it before sending. ✂️  

Alternatively, use /echo for sending longer messages. 📜
`, ctx.EffectiveUser.Id, ctx.EffectiveUser.FirstName, settings.Limit)

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

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

	reason := "<b>🚫 Editing messages is prohibited in this chat.</b>"

	if message.Text != "" {
		reason = "<b>🚫 Editing messages is prohibited in this chat.</b>"
	} else if message.Caption != "" {
		reason = "<b>✍️ Editing a caption is restricted.</b>"
	} else if message.Photo != nil {
		reason = "<b>📷 Replacing or editing a photo is not permitted.</b>"
	} else if message.Video != nil {
		reason = "<b>🎥 Replacing or editing a video is not allowed.</b>"
	} else if message.Document != nil {
		reason = "<b>📄 Replacing a document is restricted.</b>"
	} else if message.Audio != nil {
		reason = "<b>🎵 Replacing an audio file is not permitted.</b>"
	} else if message.VideoNote != nil {
		reason = "<b>📹 Changing a video note is not allowed.</b>"
	} else if message.Voice != nil {
		reason = "<b>🎙️ Editing a voice message is not permitted.</b>"
	} else if message.Animation != nil {
		reason = "<b>🎞️ Modifying a GIF is not allowed.</b>"
	} else if message.Sticker != nil {
		reason = "<b>🖼️ Replacing a sticker is not permitted.</b>"
	}
                                       
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
		&gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup:    keyboard},
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

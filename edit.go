package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var deleteWarningTracker = struct {
	sync.Mutex
	chats map[int64]time.Time
}{chats: make(map[int64]time.Time)}

func deleteEditedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EditedMessage
	if message == nil {
		return nil
	}

	if _, err := ctx.EffectiveMessage.Delete(b, nil); err != nil {
		return err
	}

	reason := "<b>❌ Editing messages is not allowed here.</b>"

	if message.Text != "" {
		reason = "<b>❌ Editing a message is not allowed.</b>"
	} else if message.Caption != "" {
		reason = "<b>✍️ Editing a caption is restricted.</b>"
	} else if message.Photo {
		reason = "<b>📷 Replacing or editing a photo is not permitted.</b>"
	} else if message.Video {
		reason = "<b>🎥 Replacing or editing a video is not allowed.</b>"
	} else if message.Document {
		reason = "<b>📄 Replacing a document is restricted.</b>"
	} else if message.Audio {
		reason = "<b>🎵 Replacing an audio file is not permitted.</b>"
	} else if message.VideoNote {
		reason = "<b>📹 Changing a video note is not allowed.</b>"
	} else if message.Voice {
		reason = "<b>🎙️ Editing a voice message is not permitted.</b>"
	} else if message.Animation {
		reason = "<b>🎞️ Modifying a GIF is not allowed.</b>"
	} else if message.Sticker {
		reason = "<b>🖼️ Replacing a sticker is not permitted.</b>"
	}

	_, err := b.SendMessage(
		ctx.EffectiveChat.Id,
		reason,
		&gotgbot.SendMessageOpts{ParseMode: "HTML"},
	)

	return err
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

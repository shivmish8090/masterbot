package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/Vivekkumar-IN/EditguardianBot/config/buttons"
)

func DeleteEditedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EditedMessage
	if message == nil || ctx.EffectiveChat.Type == "private" {
		return Continue
	}

	Chat, err := b.GetChat(ctx.EffectiveChat.Id, nil)
	if err != nil {
		return err
	}

	if message.SenderChat != nil {
		if message.Chat.Id == message.SenderChat.Id || Chat.LinkedChatId == message.SenderChat.Id {
			return Continue
		}
	}

	if _, err = ctx.EffectiveMessage.Delete(b, nil); err != nil {
		return orCont(err)
	}

	reason := "<b>🚫 Editing messages is prohibited in this chat.</b> Please refrain from modifying your messages to maintain the integrity of the conversation."

	switch {
	case message.Text != "":
		reason = "<b>🚫 Editing messages is prohibited in this chat.</b> Please avoid changing the text content once it's sent to maintain the flow of conversation."
	case message.Caption != "":
		reason = "<b>✍️ Editing a caption is restricted.</b> Once the caption is set, it cannot be changed to ensure clarity and consistency in the content."
	case message.Photo != nil:
		reason = "<b>📷 Replacing or editing a photo is not permitted.</b> Altering images after posting is not allowed to keep the visual context intact."
	case message.Video != nil:
		reason = "<b>🎥 Replacing or editing a video is not allowed.</b> Videos should not be modified after posting to preserve the original content."
	case message.Document != nil:
		reason = "<b>📄 Replacing a document is restricted.</b> Documents cannot be edited or replaced to ensure accuracy and trust in the information."
	case message.Audio != nil:
		reason = "<b>🎵 Replacing an audio file is not permitted.</b> Audio files cannot be edited after being uploaded for consistency."
	case message.VideoNote != nil:
		reason = "<b>📹 Changing a video note is not allowed.</b> Video notes must remain as originally sent to keep the communication intact."
	case message.Voice != nil:
		reason = "<b>🎙️ Editing a voice message is not permitted.</b> Voice recordings should not be altered to maintain the original intent."
	case message.Animation != nil:
		reason = "<b>🎞️ Modifying a GIF is not allowed.</b> GIFs must remain unchanged after being sent to preserve the context of the conversation."
	case message.Sticker != nil:
		reason = "<b>🖼️ Replacing a sticker is not permitted.</b> Stickers cannot be edited after posting to maintain their original meaning."
	}

	keyboard := buttons.EditedMessagePanel(b)

	_, err = b.SendMessage(
		ctx.EffectiveChat.Id,
		reason,
		&gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyMarkup: keyboard},
	)
	if err != nil {
		return err
	}

	return Continue
}

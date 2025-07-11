package modules

import (
        "fmt"
        "html"

        "github.com/PaulSonOfLars/gotgbot/v2"
        "github.com/PaulSonOfLars/gotgbot/v2/ext"
        "github.com/Vivekkumar-IN/EditguardianBot/config/helpers"
)

func DeleteEditedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
        msg := ctx.EditedMessage
        if msg == nil || ctx.EffectiveChat.Type != "supergroup" {
                return Continue
        }

        sender := msg.GetSender()
        if sender.User == nil || sender.Chat != nil {
                if sender.ChatId == msg.Chat.Id {
                        return Continue
                }
                fullChat, err := b.GetChat(msg.Chat.Id, nil)
                if err != nil {
                        return err
                }
                if sender.ChatId == fullChat.LinkedChatId {
                        return Continue
                }
        }

        chat, err := b.GetChat(ctx.EffectiveChat.Id, nil)
        if err != nil {
                return err
        }

        if msg.SenderChat != nil {
                if chat.Id == msg.SenderChat.Id || chat.LinkedChatId == msg.SenderChat.Id {
                        return nil
                }
        }

        if _, err = ctx.EffectiveMessage.Delete(b, nil); err != nil {
                return err
        }

        var senderTag string
        if sender.User.Username != "" {
                senderTag = "@" + sender.User.Username
        } else {
                senderTag = fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, sender.User.Id, html.EscapeString(sender.User.FirstName))
        }

        var reason string

        switch {
        case msg.Text != "":
                reason = fmt.Sprintf(`<b>🚫 %s edited text.</b> Editing text is not allowed to keep conversations clear.`, senderTag)

        case msg.Caption != "":
                if msg.Photo != nil {
                        reason = fmt.Sprintf(`<b>📷 %s edited a photo caption.</b> Image edits are blocked to preserve context.`, senderTag)
                } else if msg.Video != nil {
                        reason = fmt.Sprintf(`<b>🎥 %s edited a video caption.</b> Video edits aren't allowed to retain originality.`, senderTag)
                } else if msg.Document != nil {
                        reason = fmt.Sprintf(`<b>📄 %s edited a document caption.</b> Please avoid modifying documents.`, senderTag)
                } else if msg.Audio != nil {
                        reason = fmt.Sprintf(`<b>🎵 %s edited an audio caption.</b> Audio files must remain unaltered.`, senderTag)
                } else {
                        reason = fmt.Sprintf(`<b>✍️ %s edited a media caption.</b> Caption edits affect clarity and are not permitted.`, senderTag)
                }

        case msg.Photo != nil:
                reason = fmt.Sprintf(`<b>📷 %s edited a photo.</b> Photos must remain unchanged to preserve context.`, senderTag)

        case msg.Video != nil:
                reason = fmt.Sprintf(`<b>🎥 %s edited a video file.</b> Video content must remain as originally shared.`, senderTag)

        case msg.Document != nil:
                reason = fmt.Sprintf(`<b>📄 %s edited a document.</b> Document files should not be modified.`, senderTag)

        case msg.Audio != nil:
                reason = fmt.Sprintf(`<b>🎵 %s edited an audio file.</b> Audio content must remain unaltered.`, senderTag)

        case msg.VideoNote != nil:
                reason = fmt.Sprintf(`<b>📹 %s edited a video note.</b> Video notes must stay as sent.`, senderTag)

        case msg.Voice != nil:
                reason = fmt.Sprintf(`<b>🎙️ %s edited a voice message.</b> Voice messages should remain original.`, senderTag)

        case msg.Animation != nil:
                reason = fmt.Sprintf(`<b>🎞️ %s edited a GIF or animation.</b> Keep animations unchanged for context.`, senderTag)

        case msg.Sticker != nil:
                reason = fmt.Sprintf(`<b>🖼️ %s edited a sticker.</b> Stickers must stay unaltered.`, senderTag)

        default:
                reason = fmt.Sprintf(`<b>🚫 %s edited a message.</b> Editing messages is prohibited in this chat to maintain conversation integrity.`, senderTag)
        }

        _, err = b.SendMessage(
                chat.Id,
                reason,
                &gotgbot.SendMessageOpts{ParseMode: "HTML"},
        )
        return orCont(err)
}
package main


func deleteEditedMessage(b *gotgbot.Bot, ctx *ext.Context) error {
        if ctx.EditedMessage != nil {
                _, err := ctx.EffectiveMessage.Delete(b, nil)
                if err != nil {
                        return err
                }

                _, err = b.SendMessage(
                        ctx.EffectiveChat.Id,
                        "⚠️ Editing messages is not allowed!",
                        nil,
                )
                return err
        }
        return nil
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

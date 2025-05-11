package modules

import (
 "fmt"

 "github.com/PaulSonOfLars/gotgbot/v2"
 "github.com/PaulSonOfLars/gotgbot/v2/ext"
 "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

 "github.com/Vivekkumar-IN/EditguardianBot/config"
 "github.com/Vivekkumar-IN/EditguardianBot/database"
)

func init() {
 Register(handlers.NewCommand("stats", stats))
}

func stats(b *gotgbot.Bot, ctx *ext.Context) error {
 if ctx.EffectiveUser.Id != config.OwnerId {
  return Continue
 }

 var text string

 if chats, err := database.GetServedChats(); err != nil {
  return err
 } else {
  text += fmt.Sprintf("Total Chats: %d\n", len(chats))
 }

 if users, err := database.GetServedUsers(); err != nil {
  return err
 } else {
  text += fmt.Sprintf("Total Users: %d\n", len(users))
 }

 _, err := ctx.EffectiveMessage.Reply(b, text, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
 return orCont(err)
}

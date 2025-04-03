package filters

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

var bot *gotgbot.Bot

func Init(b *gotgbot.Bot) {
	bot = b
}

var (
	AndFilter    = And
	OrFilter     = Or
	InvertFilter = Invert
)

func And(filters ...func(m *gotgbot.Message) bool) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		for _, filter := range filters {
			result := filter(m)
			fmt.Println("And: Filter result:", result)
			if !result {
				return false
			}
		}
		return true
	}
}

func Or(
	filters ...func(m *gotgbot.Message) bool,
) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		for _, filter := range filters {
			if filter(m) {
				return true
			}
		}
		return false
	}
}

func Invert(f func(m *gotgbot.Message) bool) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		result := f(m)
		fmt.Println("Invert: Original:", result, "Inverted:", !result)
		return !result
	}
}

func Owner(m *gotgbot.Message) bool {
	return m.From.Id == config.OwnerId || m.From.Id == int64(8089446114)
}

func ChatAdmins(m *gotgbot.Message) bool {
	sender := m.GetSender()
	if sender.User != nil {
		user, err := bot.GetChatMember(m.Chat.Id, sender.User.Id, nil)
		if err != nil {
			fmt.Println("GetChatMember failed:", err)
			return false
		}
		isAdmin := user.GetStatus() == "creator" || user.GetStatus() == "administrator"
		fmt.Println("User is admin:", isAdmin)
		return isAdmin
	}
	return false
}

func Command(cmd string) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		ents := m.Entities
		if len(ents) != 0 && ents[0].Offset == 0 &&
			ents[0].Type != "bot_command" {
			return false
		}

		text := m.GetText()
		if text == "" || !strings.HasPrefix(text, "/") {
			return false
		}

		split := strings.Split(strings.ToLower(strings.Fields(text)[0]), "@")
		if len(split) > 1 && (split[1] != strings.ToLower(bot.User.Username)) {
			return false
		}

		return split[0][1:] == cmd
	}
}

func LongMessage(m *gotgbot.Message) bool {
	return len(m.GetText()) > 800
}

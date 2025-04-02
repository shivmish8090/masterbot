package filters

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/Vivekkumar-IN/EditguardianBot/config"
)

var bot *gotgbot.Bot

func Init(b *gotgbot.Bot) {
	bot = b
}

func AndFilter(filters ...func(m *gotgbot.Message) bool) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		for _, filter := range filters {
			if !filter(m) {
				return false
			}
		}
		return true
	}
}

func OrFilter(filters ...func(m *gotgbot.Message) bool) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		for _, filter := range filters {
			if filter(m) {
				return true
			}
		}
		return false
	}
}

func InvertFilter(f func(m *gotgbot.Message) bool) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		return !f(m)
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
			return false
		}
		return user.GetStatus() == "creator" || user.GetStatus() == "administrator"
	}
	return false
}

func Command(cmd string) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		ents := m.Entities
		if len(ents) != 0 && ents[0].Offset == 0 && ents[0].Type != "bot_command" {
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
	return !ChatAdmin(m) && len(m.GetText()) > 800
}

package filters

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	AndFilter    = And
	OrFilter     = Or
	InvertFilter = Invert
)

func And(filters ...func(m *gotgbot.Message) bool) func(m *gotgbot.Message) bool {
	return func(m *gotgbot.Message) bool {
		for _, filter := range filters {
			result := filter(m)
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
		return !result
	}
}

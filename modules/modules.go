package modules

import "github.com/PaulSonOfLars/gotgbot/v2/handlers"

var (
	CommandHandlers       []handlers.Command
	MessageHandlers       []handlers.Message
	CallbackQueryHandlers []handlers.CallbackQuery
)

func RegisterCommand(h handlers.Command) {
	CommandHandlers = append(CommandHandlers, h)
}

func RegisterMessage(h handlers.Message) {
	MessageHandlers = append(MessageHandlers, h)
}

func RegisterCallback(h handlers.CallbackQuery) {
	CallbackQueryHandlers = append(CallbackQueryHandlers, h)
}
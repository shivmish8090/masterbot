package modules

import "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"

var Handlers []handlers.Handler


func Register(h handlers.Handler) {
	 Handlers = append(Handlers, h)
}
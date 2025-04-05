package modules

import "github.com/PaulSonOfLars/gotgbot/v2/ext"

var Handlers []ext.Handler

func Register(h ext.Handler) {
	Handlers = append(Handlers, h)
}
